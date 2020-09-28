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

	for _, info := range protoInfos.Processes {
		processes = append(processes,
			rex.ProcessInfo{
				ID:       uuid.MustParse(info.ProcessUUID),
				PID:      int(info.Pid),
				ExitCode: int(info.ExitCode),
				Path:     info.Path,
				Args:     info.Args,
				OwnerID:  uuid.MustParse(info.OwnerUUID),
				Create:   time.Unix(info.Create.GetSeconds(), int64(info.Create.GetNanos())).UTC(),
				Exit:     time.Unix(info.Exit.GetSeconds(), int64(info.Exit.GetNanos())).UTC(),
			},
		)
	}

	return processes, nil
}

// NewClient creates a new GRPC Client
func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		grpcClient: proto.NewRexClient(conn),
	}
}
