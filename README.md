# Rex: Remote process execution service

## build
Install [gRPC toolkit](https://grpc.io/docs/languages/go/quickstart/) for go.
Afterwards you can build the client and server binaries:
```bash
$ make rex rexd
```

Visit `scripts/README.md` to generate certificates. You can also use
pre-generated certificates in `./fixtures/tls/{client,server,ca}`.
## Run

You can save the clients' UUIDs for future reference:
```
CL1_ID="$(openssl x509 -subject -in fixtures/tls/client/1.pem -noout -dates | grep CN | grep -o '[^ ]*$')"
CL2_ID="$(openssl x509 -subject -in fixtures/tls/client/2.pem -noout -dates | grep CN | grep -o '[^ ]*$')"
```

Next, run the server, passing the ca certificate, server key pair, plus the
required policies:
```bash
$ ./rexd \
    -ca fixtures/tls/ca/ca.crt -cert fixtures/tls/server/1.pem -key fixtures/tls/server/1.key \
    -datadir data
    -policy '{"Principal": "*", "Action": "*", "Effect": "Allow"}' \
    -policy '{"Principal": "'$CL2_ID'", "Action": "/Rex/ListProcessInfo", "Effect": "Deny"}'
```
The `-datadir` flag specifies the directory that rex will use to store
process stdout/stderr.
The `-policy` flag can be passed multiple times.
The former policy allows all
API calls by all users, otherwise no user is authorized to access any API.
The latter disallows a user with UUID equal to `$CL2_ID` from calling `/Rex/ListProcessInfo`.

Then on another terminal, first set `CL1_ARGS` and `CL2_ARGS` to contain the
TLS-related arguments for the client:
```
CL1_ARGS="-ca fixtures/tls/ca/ca.crt -cert fixtures/tls/client/1.pem -key fixtures/tls/client/1.key"
CL2_ARGS="-ca fixtures/tls/ca/ca.crt -cert fixtures/tls/client/2.pem -key fixtures/tls/client/2.key"
```

Then run
```bash
$ ./rex $CL1_ARGS exec $YOUR_COMMAND_HERE
```

For example:
```bash
$ ./rex $CL1_ARGS exec touch some_file
```

Verify the results (on the same directory that you ran `./rexd` from):
```bash
$ ls some_file
```

Executing a nonexistent file, which allows the client to catch `ErrNotFound`:
```bash
$ ./rex $CL1_ARGS exec nonexistent-binary
```

Executing a file without the execute permission:
```bash
$ ./rex $CL1_ARGS exec ./rex.go
```

Executing a command by client 2 which will likely write to both stdout and stderr:
```bash
$ TASK_ID=$(./rex $CL2_ARGS exec find / -maxdepth 3 2>/dev/null | grep \\-)
```
And to peek at the results:
```bash
$ ./rex $CL2_ARGS $TASK_ID stdout
$ ./rex $CL2_ARGS $TASK_ID stderr
```
And to make sure clients cannot access each others' resources
```bash
$ ./rex $CL2_ARGS $TASK_ID stderr
```

To get a list of processes (access is not limited to owned processes of the current user
for demonstration purposes):
```bash
$ ./rex $CL1_ARGS ps
```

To verify that client with UUID `$CL2_ID` is not allowed to call
`/Rex/ListProcessInfo/`:
```bash
$ ./rex $CL2_ARGS ps
```

To read the output of a process while it is (probably) still writing to it:
```bash
$ ./rex $CL2_ARGS read $(./rex $CL2_ARGS exec find / -maxdepth 6 2>/dev/null | grep \\-) stdout
```

To send check the status of a process and send `SIGINT` to it while it's running:
```bash
$ TASK_ID=$(./rex $CL2_ARGS exec sleep 100 2>/dev/null | grep \\-)
$ ./rex $CL2_ARGS get $TASK_ID
$ ./rex $CL2_ARGS kill $TASK_ID
```

To verify that an error is received if a signal is sent to a process that is
not running:
```bash
$ ./rex $CL2_ARGS get $TASK_ID
$ ./rex $CL2_ARGS kill $TASK_ID
```

To run tests:
```
$ make test
```
To show coverage:
```
$ make coverage
```

## Design
[Design document](https://docs.google.com/document/d/1ICGf0mDO4sh1-PH73gvYQXFNxD0CGETNpy9wnxx1UWM/edit?usp=sharing)

## License
MIT
