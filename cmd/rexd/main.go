package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"net"
	"os"
	"path"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	log "github.com/sirupsen/logrus"

	"github.com/farnasirim/rex/cmd/internal/io"
	rex_grpc "github.com/farnasirim/rex/grpc"
	"github.com/farnasirim/rex/localexec"
	"github.com/farnasirim/rex/proto"
)

type variadicFlag []string

func (f *variadicFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *variadicFlag) String() string {
	return ""
}

var (
	policyFlags  variadicFlag
	pathToCACert string
	pathToCert   string
	pathToKey    string
	dataDirFlag  string
	serveAddr    string
)

func main() {
	log.SetLevel(log.DebugLevel)
	parseAndValidate()

	var policies []rex_grpc.Policy
	for _, fl := range policyFlags {
		x, err := rex_grpc.SimpleAccessRuleFromJSON([]byte(fl))
		if err != nil {
			log.Fatalf("Policy argument malformed: %v", err)
		}
		policies = append(policies, x)
	}

	lis, err := net.Listen("tcp", serveAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tlsCredentials := getTLSCredentials()
	policyEnforcer := rex_grpc.NewPolicyEnforcer(policies...)

	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(
			rex_grpc.AuthInfoInterceptor,
			rex_grpc.PolicyEnforcementInterceptor(policyEnforcer),
			rex_grpc.ErrorMarshallerInterceptor,
		),
	)
	linuxProcessServer := localexec.NewServer(dataDirFlag)
	rexGRPCServer := rex_grpc.NewServer(linuxProcessServer)

	proto.RegisterRexServer(grpcServer, rexGRPCServer)
	log.Debugln("Serving...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln(err.Error())
	}
}

func parseAndValidate() {
	flag.Var(&policyFlags, "policy",
		"JSON formatted policy with keys Principal, Action, and Effect. Can be passed multiple times.")

	flag.StringVar(&pathToCACert, "ca", "", "path to ca certificate in pem format")
	flag.StringVar(&pathToCert, "cert", "", "path to server certificate in pem format")
	flag.StringVar(&pathToKey, "key", "", "path to server private key in pem format")
	flag.StringVar(&serveAddr, "addr", "localhost:9090", "serve address of format [ip]:port")

	dataDirDefault := os.Getenv("TMPDIR")
	if len(dataDirDefault) == 0 {
		dataDirDefault = "/tmp"
	}
	dataDirDefault = path.Join(dataDirDefault, "rex")
	flag.StringVar(&dataDirFlag, "datadir", dataDirDefault,
		"Directory to store process stdout/stderr files")

	flag.Parse()

	if pathToCACert == "" {
		log.Fatalln("Missing -ca arg")
	}

	if pathToKey == "" {
		log.Fatalln("Missing -key arg")
	}

	if pathToCert == "" {
		log.Fatalln("Missing -cert arg")
	}
}

func getTLSCredentials() credentials.TransportCredentials {
	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(io.ReadFileOrFatal(pathToCACert)); !ok {
		log.Fatalln("CA cert malformed")
	}

	cert, err := tls.LoadX509KeyPair(pathToCert, pathToKey)
	if err != nil {
		log.Fatalf("Failed to load key pair: %v\n", err)
	}

	config := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: false,
		RootCAs:            caPool,
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          caPool,
	}
	return credentials.NewTLS(config)
}
