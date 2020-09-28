package rex

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Service defines the Rex interface within Go.
type Service interface {
	// Exec executes a given executable with the supplied args.
	Exec(ctx context.Context, path string, args ...string) (uuid.UUID, error)

	// ListProcessInfo returns a list containing ProcessInfo objects, one
	// for each process previously Exec'd on this server.
	ListProcessInfo(ctx context.Context) ([]ProcessInfo, error)
}

type ProcessInfo struct {
	ID       uuid.UUID
	PID      int
	ExitCode int
	Running  bool
	Path     string
	Args     []string
	OwnerID  uuid.UUID
	Create   time.Time
	Exit     time.Time
}

type rexContextKey string

const (
	UserIDContextKey     rexContextKey = "Rex-Context-UserID"
	MethodNameContextKey rexContextKey = "Rex-Context-MethodName"
)

var (
	// ErrNotImplemented is returned by any of the API implementations
	// that are not yet implemented.
	ErrNotImplemented = errors.New("not implemented")

	// ErrTLSCredentials indicates that the TLS credentials of the peer
	// could not be extracted.
	ErrTLSCredentials = errors.New("unable to read tls credentials")

	// ErrAccessDenied is returned when a user attempts to call an API but is
	// unauthorized to do so
	ErrAccessDenied = errors.New("access denied")

	// ErrUnauthenticated is returned whenever the received request does not
	// contain information to prove authentication of the caller.
	ErrUnauthenticated = errors.New("unauthenticated")
)

func UserIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(UserIDContextKey).(string)
	return val, ok
}

func MethodNameFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(MethodNameContextKey).(string)
	return val, ok
}
