package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	log "github.com/sirupsen/logrus"

	rex_grpc "github.com/farnasirim/rex/grpc"
	"github.com/farnasirim/rex/os"
	"github.com/farnasirim/rex/proto"
)

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

	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(readFileOrFatal("scripts/ca.crt")); !ok {
		log.Fatalln("CA cert malformed")
	}

	cert, err := tls.LoadX509KeyPair("scripts/server.pem", "scripts/server.key")
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

	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	linuxProcessServer := os.NewServer()
	rexGRPCServer := rex_grpc.NewServer(linuxProcessServer)

	proto.RegisterRexServer(grpcServer, rexGRPCServer)
	log.Debugln("Serving...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln(err.Error())
	}
}
