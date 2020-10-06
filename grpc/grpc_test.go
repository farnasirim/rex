package grpc_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"

	"github.com/farnasirim/rex"
	rex_grpc "github.com/farnasirim/rex/grpc"
	"github.com/farnasirim/rex/proto"
)

func TestService_TLSHandshake(t *testing.T) {
	lis := bufconn.Listen(1e4)
	serverTLSCredentials := getServerTLSCredentials(t)
	processID := uuid.New()
	var linuxProcessServer rex.Service = &processServerMock{
		t,
		func(ctx context.Context, processID uuid.UUID) (rex.ProcessInfo, error) {
			return rex.ProcessInfo{PID: 1234}, nil
		},
	}
	rexGRPCServer := rex_grpc.NewServer(linuxProcessServer)
	grpcServer := grpc.NewServer(grpc.Creds(serverTLSCredentials))

	proto.RegisterRexServer(grpcServer, rexGRPCServer)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Errorf("Error while serving: %v", err)
		}
	}()

	clientConnection, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(getClietnTLSCredentials(t)))
	if err != nil {
		t.Errorf("Failed to create client connection: %v", err)
	}
	defer clientConnection.Close()

	client := proto.NewRexClient(clientConnection)
	info, err := client.GetProcessInfo(context.Background(),
		&proto.GetProcessInfoRequest{ProcessUUID: processID.String()})
	if err != nil {
		t.Errorf("Error in calling GetProcessInfo: %v", err)
	}
	if info.Pid != 1234 {
		t.Errorf("Expected info.Pid to be set to %q, got %q", 1234, info.Pid)
	}
}

func TestService_GetProcessInfo_APITranslation(t *testing.T) {
	lis := bufconn.Listen(1e4)
	serverTLSCredentials := getServerTLSCredentials(t)
	originalProcessID := uuid.New()
	originalOwnerID := uuid.New()
	createTime := time.Now().UTC().AddDate(-1, -1, -1)
	exitTime := createTime.Add(1 * time.Second)
	var linuxProcessServer rex.Service = &processServerMock{
		t,
		func(ctx context.Context, processID uuid.UUID) (rex.ProcessInfo, error) {
			if processID != originalProcessID {
				t.Errorf("Expected original processID to be the same as the process ID that is passed to GetProcessInfo")
				return rex.ProcessInfo{}, rex.ErrInvalidArgument
			}
			return rex.ProcessInfo{
				ID:       originalProcessID,
				PID:      123,
				ExitCode: -1,
				Running:  false,
				Path:     "/usr/bin/sleep",
				Args:     []string{"1"},
				OwnerID:  originalOwnerID,
				Create:   createTime,
				Exit:     exitTime,
			}, nil
		},
	}
	rexGRPCServer := rex_grpc.NewServer(linuxProcessServer)
	grpcServer := grpc.NewServer(grpc.Creds(serverTLSCredentials))

	proto.RegisterRexServer(grpcServer, rexGRPCServer)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Errorf("Error while serving: %v", err)
		}
	}()

	clientConnection, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(getClietnTLSCredentials(t)))
	if err != nil {
		t.Errorf("Failed to create client connection: %v", err)
	}
	defer clientConnection.Close()

	client := proto.NewRexClient(clientConnection)
	info, err := client.GetProcessInfo(context.Background(),
		&proto.GetProcessInfoRequest{ProcessUUID: originalProcessID.String()})
	if err != nil {
		t.Errorf("Error in calling GetProcessInfo: %v", err)
	}
	if info.ProcessUUID != originalProcessID.String() {
		t.Errorf("Expected the returned info to contain the same information as the translated received info: ProcessUUID")
	}
	if info.Pid != 123 {
		t.Errorf("Expected the returned info to contain the same information as the translated received info: Pid")
	}
	if info.ExitCode != -1 {
		t.Errorf("Expected the returned info to contain the same information as the translated received info: ExitCode")
	}
	if info.Path != "/usr/bin/sleep" {
		t.Errorf("Expected the returned info to contain the same information as the translated received info: Running")
	}
	if len(info.Args) != 1 {
		t.Errorf("Expected the returned info to contain the same information as the translated received info: len(info.Args)")
	}
	if info.Args[0] != "1" {
		t.Errorf("Expected the returned info to contain the same information as the translated received info: Args[0]")
	}
	if info.OwnerUUID != originalOwnerID.String() {
		t.Errorf("Expected the returned info to contain the same information as the translated received info: OwnerUUID")
	}
	if info.Create.AsTime().UTC() != createTime {
		t.Errorf("Expected the returned info to contain the same information as the translated received info: Create")
	}
	if info.Exit.AsTime().UTC() != exitTime {
		t.Errorf("Expected the returned info to contain the same information as the translated received info: Exit")
	}

}

func readFileOrFatal(filepath string, t *testing.T) []byte {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Errorf("Failed to read %s: %v\n", filepath, err)
	}
	return content
}

func getServerTLSCredentials(t *testing.T) credentials.TransportCredentials {
	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(readFileOrFatal("./../fixtures/tls/ca/ca.crt", t)); !ok {
		t.Errorf("CA cert malformed")
	}

	cert, err := tls.LoadX509KeyPair("./../fixtures/tls/server/1.pem", "./../fixtures/tls/server/1.key")
	if err != nil {
		t.Errorf("Failed to load key pair: %v\n", err)
	}

	config := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: false,
		RootCAs:            caPool,
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          caPool,
	}
	return credentials.NewTLS(config)
}

func getClietnTLSCredentials(t *testing.T) credentials.TransportCredentials {
	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(readFileOrFatal("./../fixtures/tls/ca/ca.crt", t)); !ok {
		t.Errorf("CA cert malformed")
	}

	cert, err := tls.LoadX509KeyPair("./../fixtures/tls/client/1.pem", "./../fixtures/tls/client/1.key")
	if err != nil {
		t.Errorf("Failed to load key pair: %v\n", err)
	}

	config := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: false,
		ServerName:         "localhost",
		RootCAs:            caPool,
		Certificates:       []tls.Certificate{cert},
	}

	tlsCredentials := credentials.NewTLS(config)
	return tlsCredentials
}

type processServerMock struct {
	t                  *testing.T
	GetProcessInfoFunc func(ctx context.Context, processID uuid.UUID) (rex.ProcessInfo, error)
}

func (m *processServerMock) Exec(ctx context.Context, path string, args ...string) (uuid.UUID, error) {
	m.t.Errorf("Not implemented")
	return uuid.Nil, rex.ErrNotImplemented
}
func (m *processServerMock) ListProcessInfo(ctx context.Context) ([]rex.ProcessInfo, error) {
	m.t.Errorf("Not implemented")
	return nil, rex.ErrNotImplemented
}
func (m *processServerMock) GetProcessInfo(ctx context.Context, processID uuid.UUID) (rex.ProcessInfo, error) {
	return m.GetProcessInfoFunc(ctx, processID)
}
func (m *processServerMock) Kill(ctx context.Context, processID uuid.UUID, signal int) error {
	m.t.Errorf("Not implemented")
	return rex.ErrNotImplemented
}
func (m *processServerMock) Read(ctx context.Context, processID uuid.UUID, target rex.OutputStream) ([]byte, error) {
	m.t.Errorf("Not implemented")
	return nil, rex.ErrNotImplemented
}
