package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	log "github.com/sirupsen/logrus"

	"github.com/farnasirim/rex"
	rex_grpc "github.com/farnasirim/rex/grpc"
)

var pathToCACert string
var pathToCert string
var pathToKey string

func readFileOrFatal(filepath string) []byte {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v\n", filepath, err)
	}
	return content
}

func main() {
	log.SetLevel(log.DebugLevel)

	flag.StringVar(&pathToCACert, "ca", "", "path to ca certificate in pem format")
	flag.StringVar(&pathToCert, "cert", "", "path to server certificate in pem format")
	flag.StringVar(&pathToKey, "key", "", "path to server private key in pem format")

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
	}

	tlsCredentials := credentials.NewTLS(config)
	conn, err := grpc.Dial("localhost:9090",
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithUnaryInterceptor(rex_grpc.ErrorUnmarshallerInterceptor))
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalln(err.Error())
		}
	}()

	var client rex.Service = rex_grpc.NewClient(conn)

	log.Debugln("Created new client")
	if flag.NArg() <= 1 {
		log.Fatalln("missing action")
	}
	action := flag.Arg(0)
	rest := flag.Args()[1:]

	switch action {
	case "exec":
		if len(rest) == 0 {
			log.Fatalln("missing path to executable file")
		}
		processUUID, err := client.Exec(rest[0], rest[1:]...)
		if err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				log.Debugln("Got exec.ErrNotFound")
			}
			log.Fatalln(err.Error())
		}
		fmt.Println(processUUID)
	case "kill":
		log.Fatalln("Not implemented")
	case "ps":
		log.Fatalln("Not implemented")
	case "get":
		log.Fatalln("Not implemented")
	case "read":
		log.Fatalln("Not implemented")
	default:
		log.Fatalf("Invalid action: %q", action)
	}
}
