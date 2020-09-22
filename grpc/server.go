package grpc

import (
	"context"

	"github.com/farnasirim/rex"
	"github.com/farnasirim/rex/log"
	"github.com/farnasirim/rex/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements rex.Service by translating the GRPC api and passing
// the request to a concrete implementation of rex.Service
type Server struct {
	proto.UnimplementedRexServer
	ps rex.Service
}

// Exec implements the Exec function from the Rex GRPC api.
func (s *Server) Exec(ctx context.Context, req *proto.ExecRequest) (*empty.Empty, error) {
	log.Debugln("In grpc.Exec")
	if err := s.ps.Exec(req.Path, req.Args...); err != nil {
		err := status.Error(codes.Internal, err.Error())
		return nil, err
	}
	return &empty.Empty{}, nil
}

// NewServer creates a new Server capable of serving its API
// over GRPC.
func NewServer(ps rex.Service) *Server {
	return &Server{
		ps: ps,
	}
}
