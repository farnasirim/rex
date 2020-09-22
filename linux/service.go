package linux

import (
	"os/exec"
)

// ProcessServer implements rex.Service in the Linux environment
type ProcessServer struct {
}

// Exec creates a process from the supplied path and args
func (ps *ProcessServer) Exec(path string, args ...string) error {
	cmd := exec.Command(path, args...)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// NewServer creates a new Server capable of serving its API
// over GRPC.
func NewServer() *ProcessServer {
	return &ProcessServer{}
}
