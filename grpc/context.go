package grpc

type contextKey string

const (
	userIDContextKey     contextKey = "Rex-Context-UserID"
	methodNameContextKey contextKey = "Rex-Context-MethodName"
)
