package main

import (
	"net"

	"google.golang.org/grpc"

	rex_grpc "github.com/farnasirim/rex/grpc"
	"github.com/farnasirim/rex/log"
	"github.com/farnasirim/rex/proto"
)

func main() {
	log.SetLogLevel(log.LevelDebug)

	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	rexGRPCServer := rex_grpc.NewServer(nil)

	proto.RegisterRexServer(grpcServer, rexGRPCServer)
	log.Debugln("Serving...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln(err.Error())
	}
}
