package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/yaml.v2"

	"github.com/google/uuid"
	"github.com/kataras/tablewriter"

	log "github.com/sirupsen/logrus"

	"github.com/farnasirim/rex"
	"github.com/farnasirim/rex/cmd/internal/io"
	rex_grpc "github.com/farnasirim/rex/grpc"
)

const (
	maxMsgSize int = 1e12
)

var (
	pathToCACert string
	pathToCert   string
	pathToKey    string
	serverAddr   string
	timeout      int
)

func main() {
	log.SetLevel(log.DebugLevel)
	parseAndValidate()
	conn := getGRPCConnection()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalln(err.Error())
		}
	}()

	var client rex.Service = rex_grpc.NewClient(conn)

	if flag.NArg() < 1 {
		log.Fatalln("missing action")
	}

	action := flag.Arg(0)
	rest := flag.Args()[1:]

	ctx := context.Background()
	if timeout != -1 {
		var cancelFunc func()
		ctx, cancelFunc = context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		defer cancelFunc()
	}

	switch action {
	case "exec":
		if len(rest) < 1 {
			log.Fatalln("Missing executable path")
		}
		processUUID, err := client.Exec(ctx, rest[0], rest[1:]...)
		if err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				log.Debugln("Got exec.ErrNotFound")
			}
			log.Fatalln(err.Error())
		}
		fmt.Println(processUUID)
	case "kill":
		if len(rest) < 1 {
			log.Fatalln("Missing process id")
		} else if len(rest) > 1 {
			log.Fatalf("Too many arguments to kill: got: %d, expected: %d", len(rest), 1)
		}

		processID, err := uuid.Parse(rest[0])
		if err != nil {
			log.Fatalf("Error while parsing processUUID: %v", err)
		}

		// only supports sigint for now
		err = client.Kill(ctx, processID, int(syscall.SIGINT))
		if err != nil {
			log.Fatalln(err.Error())
		}
	case "ps":
		if len(rest) > 0 {
			log.Warnf("Ignoring %d extra arguments to %q", len(rest), "ps")
		}
		processes, err := client.ListProcessInfo(ctx)
		if err != nil {
			log.Fatalln(err.Error())
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Owner ID", "Created", "State"})
		now := time.Now().UTC()
		for _, p := range processes {
			var row []string
			row = append(row, p.ID.String())
			row = append(row, p.OwnerID.String())
			row = append(row, now.Sub(p.Create).Round(time.Second).String())
			state := "running"
			if !p.Exit.IsZero() {
				state = fmt.Sprintf("Exited with code %d (%s ago)",
					p.ExitCode, now.Sub(p.Exit).Round(time.Second).String())
			}
			row = append(row, state)
			table.Append(row)
		}
		table.Render()
	case "get":
		if len(rest) < 1 {
			log.Fatalln("Missing processID argument")
		} else if len(rest) > 1 {
			log.Fatalf("Too many arguments to get: got: %d, expected: %d", len(rest), 1)
		}
		processUUID, err := uuid.Parse(rest[0])
		if err != nil {
			log.Fatalf("Bad argument %q: %v", rest[0], err)
		}
		procInfo, err := client.GetProcessInfo(ctx, processUUID)
		if err != nil {
			log.Fatalln(err.Error())
		}
		output, err := yaml.Marshal(procInfo)
		if err != nil {
			log.Fatalf("Error while presenting results: %v", err)
		}
		fmt.Print(string(output))

	case "read":
		if len(rest) < 1 {
			log.Fatalln("Missing process id")
		} else if len(rest) == 1 {
			log.Fatalln("Missing target stream (stdout/stderr)")
		} else if len(rest) > 2 {
			log.Fatalf("Too many arguments: got: %d, expected: %d", len(rest), 2)
		}
		processID, err := uuid.Parse(rest[0])
		if err != nil {
			log.Fatalf("Error while parsing processUUID: %v", err)
		}

		var targetStream rex.OutputStream
		if rest[1] != "stdout" && rest[1] != "stderr" {
			log.Fatalf("Target stream must be either %q or %q", "stdout", "stderr")
		}
		if rest[1] == "stdout" {
			targetStream = rex.StdoutStream
		} else if rest[1] == "stderr" {
			targetStream = rex.StderrStream
		}

		content, err := client.Read(ctx, processID, targetStream)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Print(string(content))

	default:
		log.Fatalf("Invalid action: %q", action)
	}
}

func parseAndValidate() {
	flag.StringVar(&pathToCACert, "ca", "", "path to ca certificate in pem format")
	flag.StringVar(&pathToCert, "cert", "", "path to server certificate in pem format")
	flag.StringVar(&pathToKey, "key", "", "path to server private key in pem format")
	flag.StringVar(&serverAddr, "addr", "localhost:9090", "server address of form [ip]:port")
	flag.IntVar(&timeout, "timeout", -1, "time limit of the exection of the command (milliseconds)")

	flag.Parse()

	if pathToCACert == "" {
		log.Fatalln("Missing -ca arg")
	}

	if pathToKey == "" {
		log.Fatalln("Missing -key arg")
	}

	if pathToCert == "" {
		log.Fatalln("Missing -cert arg")
	}

}

func getGRPCConnection() *grpc.ClientConn {
	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(io.ReadFileOrFatal(pathToCACert)); !ok {
		log.Fatalln("CA cert malformed")
	}

	cert, err := tls.LoadX509KeyPair(pathToCert, pathToKey)
	if err != nil {
		log.Fatalf("Failed to load key pair: %v\n", err)
	}

	config := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: false,
		RootCAs:            caPool,
		Certificates:       []tls.Certificate{cert},
	}

	tlsCredentials := credentials.NewTLS(config)
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithUnaryInterceptor(rex_grpc.ErrorUnmarshallerInterceptor),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return conn
}
