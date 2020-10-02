package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/farnasirim/rex"
	"github.com/farnasirim/rex/proto"
)

// Client implements rex.Service by translating the API to GRPC.
type Client struct {
	grpcClient proto.RexClient
}

// Exec implementes rex.Service.Exec by sending it over GRPC to a remote
// implementation of rex.Service
func (c *Client) Exec(ctx context.Context, path string, args ...string) (uuid.UUID, error) {
	req := &proto.ExecRequest{
		Path: path,
		Args: args,
	}

	execResponse, err := c.grpcClient.Exec(ctx, req)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return uuid.Nil, errors.New(st.Message())
		}
		return uuid.Nil, err
	}

	return uuid.Parse(execResponse.ProcessUUID)
}

// ListProcessInfo forwards a ListProcessInfo request to a remote GRPC
// implementation of rex.Service
func (c *Client) ListProcessInfo(ctx context.Context) ([]rex.ProcessInfo, error) {
	protoInfos, err := c.grpcClient.ListProcessInfo(ctx, &proto.ListProcessInfoRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, errors.New(st.Message())
		}
		return nil, err
	}
	var processes []rex.ProcessInfo

	for _, pInfo := range protoInfos.Processes {
		processes = append(processes,
			processInfoNativeFromProto(pInfo),
		)
	}

	return processes, nil
}

// GetProcessInfo translates GetProcessInfo from the native API to the GRPC
// api to get the process info of a given process
func (c *Client) GetProcessInfo(ctx context.Context, processID uuid.UUID) (rex.ProcessInfo, error) {
	pInfo, err := c.grpcClient.GetProcessInfo(ctx,
		&proto.GetProcessInfoRequest{ProcessUUID: processID.String()},
	)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return rex.ProcessInfo{}, errors.New(st.Message())
		}
		return rex.ProcessInfo{}, err
	}
	return processInfoNativeFromProto(pInfo), nil
}

// Kill translates Kill from the native API to the GRPC api to send a signal
// to a specific process
func (c *Client) Kill(ctx context.Context, processID uuid.UUID, signal int) error {
	_, err := c.grpcClient.Kill(ctx,
		&proto.KillRequest{ProcessUUID: processID.String(), Signal: int32(signal)},
	)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return errors.New(st.Message())
		}
		return err
	}
	return nil
}

// Read translates Read from the native API to the GRPC api to read the stdout
// or the stderr of a specific process
func (c *Client) Read(ctx context.Context, processID uuid.UUID, target rex.OutputStream) ([]byte, error) {
	readRequest := &proto.ReadRequest{ProcessUUID: processID.String()}
	if target == rex.StdoutStream {
		readRequest.Target = proto.ReadRequest_STDOUT
	} else if target == rex.StderrStream {
		readRequest.Target = proto.ReadRequest_STDERR
	}
	readResponse, err := c.grpcClient.Read(ctx, readRequest)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, errors.New(st.Message())
		}
		return nil, err
	}
	return readResponse.Content, nil
}

func processInfoNativeFromProto(pInfo *proto.ProcessInfo) rex.ProcessInfo {
	return rex.ProcessInfo{
		ID:       uuid.MustParse(pInfo.ProcessUUID),
		PID:      int(pInfo.Pid),
		ExitCode: int(pInfo.ExitCode),
		Path:     pInfo.Path,
		Args:     pInfo.Args,
		Running:  pInfo.Running,
		OwnerID:  uuid.MustParse(pInfo.OwnerUUID),
		Create:   time.Unix(pInfo.Create.GetSeconds(), int64(pInfo.Create.GetNanos())).UTC(),
		Exit:     time.Unix(pInfo.Exit.GetSeconds(), int64(pInfo.Exit.GetNanos())).UTC(),
	}
}

// NewClient creates a new GRPC Client
func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		grpcClient: proto.NewRexClient(conn),
	}
}
