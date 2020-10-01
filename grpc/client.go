package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/farnasirim/rex/proto"
)

// Client implements rex.Service by translating the API to GRPC.
type Client struct {
	grpcClient proto.RexClient
}

// Exec implementes rex.Service.Exec by sending it over GRPC to a remote
// implementation of rex.Service
func (c *Client) Exec(path string, args ...string) (uuid.UUID, error) {
	req := &proto.ExecRequest{
		Path: path,
		Args: args,
	}

	execResponse, err := c.grpcClient.Exec(context.Background(), req)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return uuid.Nil, errors.New(st.Message())
		}
		return uuid.Nil, err
	}

	return uuid.Parse(execResponse.ProcessUUID)
}

// NewClient creates a new GRPC Client
func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		grpcClient: proto.NewRexClient(conn),
	}
}
