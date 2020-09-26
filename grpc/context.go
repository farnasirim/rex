package grpc

import (
	"context"
)

const (
	userIDContextKey     = "Rex-Context-UserID"
	methodNameContextKey = "Rex-Context-MethodName"
)

func stringFromContext(ctx context.Context, key string) (string, bool) {
	val, ok := ctx.Value(key).(string)
	if !ok {
		return "", false
	}
	return val, true
}
