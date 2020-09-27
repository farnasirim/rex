package localexec

import (
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// ProcessServer implements rex.Service in the Linux environment
type ProcessServer struct {
	processes sync.Map
	dataDir   string
}

// Exec creates a process from the supplied path and args
func (ps *ProcessServer) Exec(path string, args ...string) (uuid.UUID, error) {
	cmd := exec.Command(path, args...)

	processUUID := uuid.New()
	stdout, stderr, err := ps.createOutputFiles(processUUID)

	if err != nil {
		return uuid.Nil, err
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// TODO: would be better to get the exact start time from /proc/$pid/stat
	// I still don't see an easy way to find the exact exit time however.
	startTime := time.Now().UTC()
	if err := cmd.Start(); err != nil {
		log.Infof("failed starting a process: %v", err)
		return uuid.Nil, err
	}
	ps.registerProcess(processUUID, cmd, startTime)

	return processUUID, nil
}

func (ps *ProcessServer) registerProcess(processUUID uuid.UUID,
	cmd *exec.Cmd, start time.Time) {

	handle := &processHandle{
		id:        processUUID,
		cmd:       cmd,
		startTime: start,
	}
	ps.processes.Store(processUUID, handle)
	go func() {
		err := handle.cmd.Wait()
		func() {
			handle.m.Lock()
			defer handle.m.Unlock()
			handle.finishTime = time.Now().UTC()
			handle.waitError = err
		}()
	}()
}

func (ps *ProcessServer) getStdoutFilename(processUUID uuid.UUID) string {
	return path.Join(ps.dataDir, "proc", processUUID.String(), "stdout")
}

func (ps *ProcessServer) getStderrFilename(processUUID uuid.UUID) string {
	return path.Join(ps.dataDir, "proc", processUUID.String(), "stderr")
}

func (ps *ProcessServer) createOutputFiles(processUUID uuid.UUID) (*os.File, *os.File, error) {
	err := os.MkdirAll(path.Dir(ps.getStderrFilename(processUUID)), 0755)
	if err != nil {
		return nil, nil, err
	}
	stdout, err := os.Create(ps.getStdoutFilename(processUUID))
	if err != nil {
		return nil, nil, err
	}

	stderr, err := os.Create(ps.getStderrFilename(processUUID))
	if err != nil {
		defer stdout.Close()
		return nil, nil, err
	}

	return stdout, stderr, nil
}

type processHandle struct {
	id         uuid.UUID
	cmd        *exec.Cmd
	startTime  time.Time
	finishTime time.Time
	waitError  error
	m          sync.RWMutex
}

func NewServer(dataDir string) *ProcessServer {
	return &ProcessServer{
		dataDir: dataDir,
	}
}
