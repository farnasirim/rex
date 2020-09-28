package grpc

import (
	"context"

	"github.com/farnasirim/rex"
	"github.com/farnasirim/rex/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// Server implements rex.Service by translating the GRPC api and passing
// the request to a concrete implementation of rex.Service
type Server struct {
	proto.UnimplementedRexServer
	ps rex.Service
}

// Exec implements the Exec function from the Rex GRPC api.
func (s *Server) Exec(ctx context.Context, req *proto.ExecRequest) (*proto.ExecResponse, error) {
	processUUID, err := s.ps.Exec(ctx, req.Path, req.Args...)
	if err != nil {
		return nil, err
	}
	return &proto.ExecResponse{ProcessUUID: processUUID.String()}, nil
}

// ListProcessInfo passes an incoming ListProcessInfo GRPC request to a
// concrete implementation of rex.Service
func (s *Server) ListProcessInfo(ctx context.Context, req *proto.ListProcessInfoRequest) (*proto.ProcessInfoList, error) {
	processes, err := s.ps.ListProcessInfo(ctx)
	if err != nil {
		return nil, err
	}

	var protoInfos []*proto.ProcessInfo
	for _, proc := range processes {
		protoInfos = append(protoInfos, &proto.ProcessInfo{
			ProcessUUID: proc.ID.String(),
			Pid:         int32(proc.PID),
			ExitCode:    int32(proc.ExitCode),
			Path:        proc.Path,
			Args:        proc.Args,
			OwnerUUID:   proc.OwnerID.String(),
			Create: &timestamp.Timestamp{
				Seconds: proc.Create.Unix(),
				Nanos:   int32(proc.Create.Nanosecond())},
			Exit: &timestamp.Timestamp{
				Seconds: proc.Exit.Unix(),
				Nanos:   int32(proc.Exit.Nanosecond())},
		})
	}

	return &proto.ProcessInfoList{Processes: protoInfos}, nil
}

// NewServer creates a new Server capable of serving its API
// over GRPC.
func NewServer(ps rex.Service) *Server {
	return &Server{
		ps: ps,
	}
}
