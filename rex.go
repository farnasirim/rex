package rex

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// OutputStream specifies either stdout or stderr of a process
type OutputStream int

const (
	// StdoutStream specifies stdout of a process
	StdoutStream OutputStream = 0x1
	// StderrStream specifies stderr of a process
	StderrStream OutputStream = 0x2
)

type rexContextKey string

const (
	userIDContextKey rexContextKey = "Rex-Context-UserID"
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

// ProcessInfo contains various informations about a process.
type ProcessInfo struct {
	// ID is the unique identifier of a process.
	ID uuid.UUID
	// PID is the pid of the process in the OS. Can be recycled by the OS
	// after the process exits.
	PID int
	// ExitCode will hold the exit code of the process after it exits. It is
	// undefined if Running=true.
	ExitCode int
	// Running specifies whether or not the process is currently running.
	Running bool
	// Path is the address to the executable corresponding to the process.
	Path string
	// Args is the list of the command line arguments that were passed to the
	// process upon creation.
	Args []string
	// OwnerID is the unique identifier of the owner of the process.
	OwnerID uuid.UUID
	// Create is the point in time (UTC) at which the process was created.
	Create time.Time
	// Exit is the point in time (UTC) at which the process exited. It is
	// undefined if Running=true.
	Exit time.Time
}

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

// UserIDFromContext gets the unique identifier of the API user. Returns
// false as the second argument if no such key is found in the context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(userIDContextKey).(string)
	return val, ok
}

// WithUserID adds the supplied user ID to the given context and returns
// the resulting context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}
