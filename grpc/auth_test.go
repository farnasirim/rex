package grpc_test

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rex_grpc "github.com/farnasirim/rex/grpc"
)

func TestAuthInfoInterceptor_Unauthenticated(t *testing.T) {

	_, err := rex_grpc.AuthInfoInterceptor(context.Background(), nil, nil, nil)
	if err == nil {
		t.Errorf("Expected non nill err")
	}
	status, ok := status.FromError(err)
	if !ok {
		t.Errorf("Expected grpc error")
	}
	if status.Code() != codes.Unauthenticated {
		t.Errorf("Expected code %d (unauthenticated), received: %d", codes.Unauthenticated, status.Code())
	}
}
