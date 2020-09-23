package grpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"

	"github.com/farnasirim/rex/proto"
)

// Client implements rex.Service by translating the API to GRPC.
type Client struct {
	grpcClient proto.RexClient
}

// Exec implementes rex.Service.Exec by sending it over GRPC to a remote
// implementation of rex.Service
func (c *Client) Exec(path string, args ...string) error {
	req := &proto.ExecRequest{
		Path: path,
		Args: args,
	}
	if _, err := c.grpcClient.Exec(context.Background(), req); err != nil {
		if st, ok := status.FromError(err); ok {
			log.Errorln(err.Error())
			return errFromCode(st.Code())
		}
		return err
	}

	return nil
}

// NewClient creates a new GRPC Client
func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		grpcClient: proto.NewRexClient(conn),
	}
}
