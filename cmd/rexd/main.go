package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"path"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	log "github.com/sirupsen/logrus"

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

var policyFlags variadicFlag
var pathToCACert string
var pathToCert string
var pathToKey string
var dataDirFlag string

func readFileOrFatal(filepath string) []byte {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", filepath, err)
	}
	return content
}

// TODO: lots of duplication in rex/main.go and rexd/main.go
func main() {
	log.SetLevel(log.DebugLevel)

	flag.Var(&policyFlags, "policy",
		"JSON formatted policy with keys Principal, Action, and Effect. Can be passed multiple times.")

	flag.StringVar(&pathToCACert, "ca", "", "path to ca certificate in pem format")
	flag.StringVar(&pathToCert, "cert", "", "path to server certificate in pem format")
	flag.StringVar(&pathToKey, "key", "", "path to server private key in pem format")

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

	var policies []rex_grpc.Policy
	for _, fl := range policyFlags {
		x, err := rex_grpc.SimpleAccessRuleFromJSON([]byte(fl))
		if err != nil {
			log.Fatalf("Policy argument malformed: %v", err)
		}
		policies = append(policies, x)
	}

	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(readFileOrFatal(pathToCACert)); !ok {
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

	tlsCredentials := credentials.NewTLS(config)

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
