package localexec_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/farnasirim/rex"
	"github.com/farnasirim/rex/localexec"
)

func TestRead_Stdout(t *testing.T) {
	dataDir := os.TempDir()
	s := localexec.NewServer(dataDir)

	ctx := context.Background()
	ctx = rex.WithUserID(ctx, uuid.New().String())

	procID, err := s.Exec(ctx, "echo", []string{"hello"}...)

	if err != nil {
		t.Errorf("While calling Exec: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	info, err := s.GetProcessInfo(ctx, procID)
	if err != nil {
		t.Errorf("While calling GetProcessInfo: %v", err)
	}

	if info.Running == true {
		t.Errorf("Expected the process to finish running after 10 milliseconds")
	}

	content, err := s.Read(ctx, procID, rex.StdoutStream)
	if err != nil {
		t.Errorf("While calling Read: %v", err)
	}

	exp := fmt.Sprintf("hello\n")
	if string(content) != exp {
		t.Errorf("Expected: %s, actual: %s", exp, content)
	}
}

func TestGetProcessInfo_Exit(t *testing.T) {
	dataDir := os.TempDir()
	s := localexec.NewServer(dataDir)

	ctx := context.Background()
	ctx = rex.WithUserID(ctx, uuid.New().String())

	procID, err := s.Exec(ctx, "sleep", []string{"1"}...)

	if err != nil {
		t.Errorf("While calling Exec: %v", err)
	}

	info, err := s.GetProcessInfo(ctx, procID)
	if err != nil {
		t.Errorf("While calling GetProcessInfo: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	if info.Running != true {
		t.Errorf("Expected sleep 1 to be running after 100 milliseconds")
	}

	for info.Running == true {
		info, err = s.GetProcessInfo(ctx, procID)
		time.Sleep(100 * time.Millisecond)
	}

	runningTime := info.Exit.Sub(info.Create)
	if runningTime.Round(time.Millisecond) < 1*time.Second || 2*time.Second < runningTime.Round(time.Millisecond) {
		t.Errorf("Expected running time of sleep 1 to be approximately 1 second. Actual: %d ms",
			runningTime.Round(time.Millisecond))
	}
}

func TestKill(t *testing.T) {
	dataDir := os.TempDir()
	s := localexec.NewServer(dataDir)
	ctx := context.Background()
	ctx = rex.WithUserID(ctx, uuid.New().String())

	procID, err := s.Exec(ctx, "sleep", []string{"1"}...)

	if err != nil {
		t.Errorf("While calling Exec: %v", err)
	}

	info, err := s.GetProcessInfo(ctx, procID)
	if err != nil {
		t.Errorf("While calling GetProcessInfo: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	if info.Running != true {
		t.Errorf("Expected sleep 1 to be running after 100 milliseconds")
	}

	err = s.Kill(ctx, procID, int(syscall.SIGINT))
	if err != nil {
		t.Errorf("Expected Kill to succeed, being called on a process that should be running")
	}

	for info.Running == true && time.Now().UTC().Sub(info.Create) < time.Second {
		info, err = s.GetProcessInfo(ctx, procID)
		time.Sleep(100 * time.Millisecond)
	}

	err = s.Kill(ctx, procID, int(syscall.SIGINT))
	if err == nil {
		t.Errorf("Expected Kill to fail, being called on a process that should has exited")
	}

	runningTime := info.Exit.Sub(info.Create)
	if runningTime.Round(time.Millisecond) > 1*time.Second {
		t.Errorf("Expected running time of sleep 1 to be less than 1 second. When it's killed. Actual: %d ms",
			runningTime.Round(time.Millisecond))
	}
}

func TestExec_Unauthenticated(t *testing.T) {
	dataDir := os.TempDir()
	s := localexec.NewServer(dataDir)

	ctx := context.Background()

	_, err := s.Exec(ctx, "sleep", []string{"1"}...)
	if err != rex.ErrUnauthenticated {
		t.Errorf("Expected error %v, actual: %v", rex.ErrUnauthenticated, err)
	}
}

func TestExec_AuthBadDataDir(t *testing.T) {
	dir := os.TempDir()
	file, err := ioutil.TempFile(dir, "regularfile")
	if err != nil {
		t.Errorf("Failed to create temp file in the temp dir: %v", err)
	}
	badDirName := file.Name()
	s := localexec.NewServer(badDirName)

	ctx := context.Background()
	ctx = rex.WithUserID(ctx, uuid.New().String())

	_, err = s.Exec(ctx, "sleep", []string{"1"}...)
	if err == nil || err == rex.ErrUnauthenticated {
		t.Errorf("Expected non-nil and non %v error, got: %v", rex.ErrUnauthenticated, err)
	}
}

func TestExec_NoAuthBadDataDir(t *testing.T) {
	dir := os.TempDir()
	file, err := ioutil.TempFile(dir, "regularfile")
	if err != nil {
		t.Errorf("Failed to create temp file in the temp dir: %v", err)
	}
	badDirName := file.Name()
	s := localexec.NewServer(badDirName)

	ctx := context.Background()

	_, err = s.Exec(ctx, "sleep", []string{"1"}...)
	if err != rex.ErrUnauthenticated {
		t.Errorf("Expected error %v, actual: %v", rex.ErrUnauthenticated, err)
	}
}

func TestListProcessInfo_Sort(t *testing.T) {
	dataDir := os.TempDir()
	s := localexec.NewServer(dataDir)
	ctx := context.Background()
	ctx = rex.WithUserID(ctx, uuid.New().String())

	args := []int{2, 3, 1, 5, 4}
	for _, arg := range args {
		strArg := strconv.FormatInt(int64(arg), 10)
		_, _ = s.Exec(ctx, "echo", []string{strArg}...)
	}

	ls, err := s.ListProcessInfo(ctx)
	if err != nil {
		t.Errorf("While calling ListProcessInfo: %v", err)
	}
	if len(ls) != 5 {
		t.Errorf("Expected 5 elements in the process list, got: %d", len(ls))
	}
	for i, _ := range args {
		if ls[i].Args[0] != strconv.FormatInt(int64(args[len(args)-1-i]), 10) {
			t.Errorf("Expected %s at index %d, got %s",
				strconv.FormatInt(int64(args[len(args)-1-i]), 10), i, ls[i].Args[0])
		}
	}
}
