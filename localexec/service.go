package localexec

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"

	"github.com/farnasirim/rex"
)

// ProcessServer implements rex.Service in the Linux environment
type ProcessServer struct {
	processes sync.Map
	dataDir   string
}

// Exec creates a process from the supplied path and args
func (ps *ProcessServer) Exec(ctx context.Context,
	path string, args ...string) (uuid.UUID, error) {
	cmd := exec.Command(path, args...)

	ownerID, ok := rex.UserIDFromContext(ctx)
	if !ok {
		return uuid.Nil, rex.ErrUnauthenticated
	}

	processID := uuid.New().String()
	stdout, stderr, err := ps.createOutputFiles(processID)
	if err != nil {
		return uuid.Nil, err
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// TODO: would be better to get the exact start time from /proc/$pid/stat
	// I still don't see an easy way to find the exact exit time however.
	create := time.Now().UTC()
	if err := cmd.Start(); err != nil {
		log.Infof("failed starting a process: %v", err)
		return uuid.Nil, err
	}
	ps.registerProcess(processID, ownerID, cmd, create)

	return uuid.MustParse(processID), nil
}

// ListProcessInfo returs a list of all processes that have ever been
// successfully Exec'd into the system, sorted by their creation time (newest
// first).
func (ps *ProcessServer) ListProcessInfo(ctx context.Context) ([]rex.ProcessInfo, error) {
	var infoList []rex.ProcessInfo
	ps.processes.Range(func(key, value interface{}) bool {
		infoList = append(infoList, value.(*processHandle).getProcessInfo())
		return true
	})
	sort.Slice(infoList, func(i, j int) bool {
		return infoList[i].Create.After(infoList[j].Create)
	})
	return infoList, nil
}

// GetProcessInfo returns the process info corresponding to the givne processID
func (ps *ProcessServer) GetProcessInfo(ctx context.Context, processID uuid.UUID) (rex.ProcessInfo, error) {
	mustBeProcessHandle, ok := ps.processes.Load(processID.String())
	// Must be a little careful here. We are leaking information about the
	// process UUID's in the system. Might make more sense to make the not found
	// case indistinguishable from the unauthorized/unauthenticated case.
	//
	// A more thorough and natural (not hardcoded) RBAC implementation will
	// take care of this in a natural way, since for example in this case, it
	// *cannot find* a process with the given ID, satisfying the condition of
	// proc.OwnerID == currentUser
	if !ok {
		return rex.ProcessInfo{}, rex.ErrNotFound
	}
	userID, ok := rex.UserIDFromContext(ctx)
	if !ok {
		return rex.ProcessInfo{}, rex.ErrUnauthenticated
	}
	if mustBeProcessHandle.(*processHandle).ownerID != userID {
		return rex.ProcessInfo{}, rex.ErrAccessDenied
	}
	return mustBeProcessHandle.(*processHandle).getProcessInfo(), nil
}

// Kill sends a signal to the given process
func (ps *ProcessServer) Kill(ctx context.Context, processID uuid.UUID, signal int) error {
	mustBeProcessHandle, ok := ps.processes.Load(processID.String())
	if !ok {
		return rex.ErrNotFound
	}
	userID, ok := rex.UserIDFromContext(ctx)
	if !ok {
		return rex.ErrUnauthenticated
	}
	handle := mustBeProcessHandle.(*processHandle)
	if handle.ownerID != userID {
		return rex.ErrAccessDenied
	}

	handle.m.Lock()
	defer handle.m.Unlock()
	return handle.cmd.Process.Signal(syscall.Signal(signal))
}

// Read reads either the stdout or the stderr of the given process
func (ps *ProcessServer) Read(ctx context.Context, processID uuid.UUID, target rex.OutputStream) ([]byte, error) {
	mustBeProcessHandle, ok := ps.processes.Load(processID.String())
	if !ok {
		return nil, rex.ErrNotFound
	}
	userID, ok := rex.UserIDFromContext(ctx)
	if !ok {
		return nil, rex.ErrUnauthenticated
	}
	if mustBeProcessHandle.(*processHandle).ownerID != userID {
		return nil, rex.ErrAccessDenied
	}

	var targetFile string
	if target == rex.StderrStream {
		targetFile = ps.getStderrFilename(processID.String())
	} else if target == rex.StdoutStream {
		targetFile = ps.getStdoutFilename(processID.String())
	}

	if targetFile == "" {
		return nil, rex.ErrInvalidArgument
	}

	return ioutil.ReadFile(targetFile)
}

func (ps *ProcessServer) registerProcess(processID, ownerID string,
	cmd *exec.Cmd, create time.Time) {

	handle := &processHandle{
		id:      processID,
		ownerID: ownerID,
		cmd:     cmd,
		running: true,
		create:  create,
		pid:     cmd.Process.Pid,
	}
	ps.processes.Store(processID, handle)
	go func() {
		err := handle.cmd.Wait()

		handle.m.Lock()
		defer handle.m.Unlock()

		handle.running = false
		handle.exit = time.Now().UTC()
		handle.exitcode = handle.cmd.ProcessState.ExitCode()
		handle.waitError = err
	}()
}

func (ps *ProcessServer) getStdoutFilename(processID string) string {
	return path.Join(ps.dataDir, "proc", processID, "stdout")
}

func (ps *ProcessServer) getStderrFilename(processID string) string {
	return path.Join(ps.dataDir, "proc", processID, "stderr")
}

// createOutputFiles leaves the responsibility of closing the returned files
// to the caller if error != nil
func (ps *ProcessServer) createOutputFiles(processID string) (*os.File, *os.File, error) {
	err := os.MkdirAll(path.Dir(ps.getStderrFilename(processID)), 0755)
	if err != nil {
		return nil, nil, err
	}
	stdout, err := os.Create(ps.getStdoutFilename(processID))
	if err != nil {
		return nil, nil, err
	}

	stderr, err := os.Create(ps.getStderrFilename(processID))
	if err != nil {
		{
			if err := stdout.Close(); err != nil {
				log.Errorf("Failed to close stdout file: %v", err)
			}
		}
		return nil, nil, err
	}

	return stdout, stderr, nil
}

type processHandle struct {
	id        string
	ownerID   string
	pid       int
	cmd       *exec.Cmd
	create    time.Time
	exit      time.Time
	exitcode  int
	running   bool
	waitError error
	m         sync.RWMutex
}

func (ph *processHandle) getProcessInfo() rex.ProcessInfo {
	ph.m.RLock()
	defer ph.m.RUnlock()
	info := rex.ProcessInfo{
		ID:      uuid.MustParse(ph.id),
		PID:     ph.pid,
		Running: ph.running,
		Path:    ph.cmd.Path,
		Args:    ph.cmd.Args[1:],
		Create:  ph.create,
		OwnerID: uuid.MustParse(ph.ownerID),
	}
	if !ph.running {
		info.Exit = ph.exit
		info.ExitCode = ph.exitcode
	}
	return info
}

// NewServer creates a ProcessServer which is a concrete implementation of
// rex.Server.
func NewServer(dataDir string) *ProcessServer {
	return &ProcessServer{
		dataDir: dataDir,
	}
}
