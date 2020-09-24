package grpc

import (
	"context"

	"github.com/farnasirim/rex"
	"github.com/farnasirim/rex/proto"
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
func (s *Server) Exec(ctx context.Context, req *proto.ExecRequest) (*proto.ExecResponse, error) {
	processUUID, err := s.ps.Exec(req.Path, req.Args...)
	if err != nil {
		err := status.Error(codes.Internal, err.Error())
		return nil, err
	}
	return &proto.ExecResponse{ProcessUUID: processUUID.String()}, nil
}

// NewServer creates a new Server capable of serving its API
// over GRPC.
func NewServer(ps rex.Service) *Server {
	return &Server{
		ps: ps,
	}
}
