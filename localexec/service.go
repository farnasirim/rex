package localexec

import (
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// ProcessServer implements rex.Service in the Linux environment
type ProcessServer struct {
}

// Exec creates a process from the supplied path and args
func (ps *ProcessServer) Exec(path string, args ...string) (uuid.UUID, error) {
	cmd := exec.Command(path, args...)

	if err := cmd.Start(); err != nil {
		log.Infof("failed starting a process: %v", err)
		return uuid.Nil, err
	}

	if err := cmd.Wait(); err != nil {
		return uuid.Nil, err
	}

	processUUID := uuid.New()

	return processUUID, nil
}

// NewServer creates a new Server capable of serving its API
// over GRPC.
func NewServer() *ProcessServer {
	return &ProcessServer{}
}
