package grpc

import (
	"context"
)

type grpcContextKey string

const (
	methodNameContextKey = "Rex-GRPC-Context-MethodName"
)

func methodNameFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(methodNameContextKey).(string)
	return val, ok
}

func withMethodName(ctx context.Context, methodName string) context.Context {
	return context.WithValue(ctx, methodNameContextKey, methodName)
}
