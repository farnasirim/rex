package rex

import (
	"errors"

	"github.com/google/uuid"
)

// Service defines the Rex interface within Go.
type Service interface {
	// Exec executes a given executable with the supplied args.
	Exec(path string, args ...string) (uuid.UUID, error)
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
)
