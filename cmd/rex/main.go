package main

import (
	"google.golang.org/grpc"

	rex_grpc "github.com/farnasirim/rex/grpc"
	"github.com/farnasirim/rex/log"
)

func main() {
	log.SetLogLevel(log.LevelDebug)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:9090", opts...)
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

	if err := client.Exec("ls", "~"); err != nil {
		log.Fatalln(err.Error())
	}
}
