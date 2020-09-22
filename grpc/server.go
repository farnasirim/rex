package grpc

import (
	"context"

	"github.com/farnasirim/rex"
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

	err := status.Error(codes.Unimplemented,
		errFromCode(codes.Unimplemented).Error())
	return &empty.Empty{}, err
}

// NewServer creates a new Server capable of serving its API
// over GRPC.
func NewServer(ps rex.Service) *Server {
	return &Server{
		ps: ps,
	}
}
