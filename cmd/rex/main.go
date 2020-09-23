package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	rex_grpc "github.com/farnasirim/rex/grpc"
	log "github.com/sirupsen/logrus"
)

func readFileOrFatal(filepath string) []byte {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v\n", filepath, err)
	}
	return content
}

func main() {
	log.SetLevel(log.DebugLevel)

	flag.Parse()

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(readFileOrFatal("scripts/ca.crt")); !ok {
		log.Fatalln("CA cert malformed")
	}

	cert, err := tls.LoadX509KeyPair("scripts/client.pem", "scripts/client.key")
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
		grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalln(err.Error())
		}
	}()

	client := rex_grpc.NewClient(conn)

	log.Debugln("Created new client")
	if flag.NArg() > 1 {
		action := flag.Arg(0)
		rest := flag.Args()[1:]

		switch action {
		case "exec":
			if err := client.Exec(rest[0], rest[1:]...); err != nil {
				log.Fatalln(err.Error())
			}
			break
		case "kill":
			log.Fatalln("Not implemented")
			break
		case "ps":
			log.Fatalln("Not implemented")
			break
		case "get":
			log.Fatalln("Not implemented")
			break
		case "read":
			log.Fatalln("Not implemented")
			break
		}
	} else {
		log.Fatalln("invalid/missing action")
	}

}
