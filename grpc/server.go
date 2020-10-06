package grpc

import (
	"context"

	"github.com/farnasirim/rex"
	"github.com/farnasirim/rex/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
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
		protoInfos = append(protoInfos, processInfoProtoFromNative(proc))
	}

	return &proto.ProcessInfoList{Processes: protoInfos}, nil
}

// GetProcessInfo translates the request to get the info of a specific process
// from the gRPC API to the native API.
func (s *Server) GetProcessInfo(ctx context.Context, req *proto.GetProcessInfoRequest) (*proto.ProcessInfo, error) {
	processUUID, err := uuid.Parse(req.ProcessUUID)
	if err != nil {
		return nil, err
	}
	processInfo, err := s.ps.GetProcessInfo(ctx, processUUID)
	if err != nil {
		return nil, err
	}

	return processInfoProtoFromNative(processInfo), nil
}

// Kill will send the specified signal to the specified proces. Here it will be
// translated from the gRPC API to the native API and passed to the concrete
// implementation.
func (s *Server) Kill(ctx context.Context, req *proto.KillRequest) (*proto.KillResponse, error) {
	processUUID, err := uuid.Parse(req.GetProcessUUID())
	if err != nil {
		return nil, err
	}
	return &proto.KillResponse{}, s.ps.Kill(ctx, processUUID, int(req.GetSignal()))
}

// Read forwards a requet to read the stdout/stderr of a process to the
// underlying (concrete) rex.Service.
func (s *Server) Read(ctx context.Context, req *proto.ReadRequest) (*proto.ReadResponse, error) {
	processUUID, err := uuid.Parse(req.GetProcessUUID())
	if err != nil {
		return nil, err
	}

	var outputStream rex.OutputStream
	if req.GetTarget() == proto.ReadRequest_STDOUT {
		outputStream = rex.StdoutStream
	} else if req.GetTarget() == proto.ReadRequest_STDERR {
		outputStream = rex.StderrStream
	} else {
		return nil, rex.ErrInvalidArgument
	}

	output, err := s.ps.Read(ctx, processUUID, outputStream)
	if err != nil {
		return nil, err
	}
	return &proto.ReadResponse{Content: output}, nil
}

func processInfoProtoFromNative(proc rex.ProcessInfo) *proto.ProcessInfo {
	return &proto.ProcessInfo{
		ProcessUUID: proc.ID.String(),
		Pid:         int32(proc.PID),
		ExitCode:    int32(proc.ExitCode),
		Path:        proc.Path,
		Args:        proc.Args,
		Running:     proc.Running,
		OwnerUUID:   proc.OwnerID.String(),
		Create: &timestamp.Timestamp{
			Seconds: proc.Create.Unix(),
			Nanos:   int32(proc.Create.Nanosecond())},
		Exit: &timestamp.Timestamp{
			Seconds: proc.Exit.Unix(),
			Nanos:   int32(proc.Exit.Nanosecond())},
	}
}

// NewServer creates a new Server capable of serving its API
// over GRPC.
func NewServer(ps rex.Service) *Server {
	return &Server{
		ps: ps,
	}
}
