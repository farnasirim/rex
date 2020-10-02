package rex

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type OutputStream int

const (
	StdoutStream OutputStream = 0x1
	StderrStream OutputStream = 0x2
)

// Service defines the Rex interface within Go.
type Service interface {
	// Exec executes a given executable with the supplied args.
	Exec(ctx context.Context, path string, args ...string) (uuid.UUID, error)

	// ListProcessInfo returns a list containing ProcessInfo objects, one
	// for each process previously Exec'd on this server.
	ListProcessInfo(ctx context.Context) ([]ProcessInfo, error)

	// GetProcessInfo returns a ProcessInfo object corresponding to a
	// the processID that is provided
	GetProcessInfo(ctx context.Context, processID uuid.UUID) (ProcessInfo, error)

	// Kill sends the specified signal to the specified process
	Kill(ctx context.Context, processID uuid.UUID, signal int) error

	// Read returns the content of the stdout or the stderr of a process
	Read(ctx context.Context, processID uuid.UUID, target OutputStream) ([]byte, error)
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

	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("not found")

	// ErrInvalidArgument is when an invalid arugment is given to a function
	ErrInvalidArgument = errors.New("not found")
)

func UserIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(UserIDContextKey).(string)
	return val, ok
}

func MethodNameFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(MethodNameContextKey).(string)
	return val, ok
}
