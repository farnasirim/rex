# Rex: Remote process execution service

## TODO
 - Command line parsing and error handling is tightly knitted with the code.
   which might prevent good error reporting. Also makes it non-trivial to
   keep cli documentation in sync with code. Given some more time, I would
   have cleaned this up to separate the CLI function implementations away from
   CLI syntax (and possibly semantics) descriptions.
   (the [cobra](https://github.com/spf13/cobra) and
   [viper](https://github.com/spf13/viper) combo is my favorite solution).
 - If serializing/deserializing error chains is useful and keeping it would
   be beneficial, it has to be clear what it does with grpc status codes and
   errors that are generated in the grpc layer itself.
 - Authz: Using grpc's authorization [engine](https://pkg.go.dev/google.golang.org/grpc/security/authorization@v0.0.0-20201001231224-bebda80b05da/engine) or any well known authorization scheme
   would have made things less arbitrary and spontaneous.
 - Use of UUIDs (vs strings) for ID objects is almost arbitrary across the
   code. It would be better to completely get rid of the UUID-ness in the
   API descriptions as an implementation detail and encapsulate all of its
   usage in `localexec`. Everyone else will treat unique IDs as unique strings
   without any certain conditions.
 - A bit low level, but in retrospect `rex.ProcessInfo` being treated as a
   value type (as opposed to pointer type) across API boundaries was a mistake.
   Not only there is not enough reason to treat it as a value type, but also
   it's nullability is quite desired across the APIs.
 - Also low level: creating global structs (`rex.ProcessInfo`) must be done
   with more care. Currently no API is safe against sudden change to addition
   of a field to `rex.ProcessInfo`. The behavior is not unnatural as the new
   field will simply be unsupported across the system, but still it would have
   been nice if we could have a compile time trap for this. Maybe `NewProcessInfo`?

## Build
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

An optional timeout (milliseconds) argument can be passed to the cli:
```
$ TASK_ID=$(./rex $CL2_ARGS exec find / | grep \\-)
$ ./rex $CL2_ARGS -timeout 1000 read $TASK_ID stdout
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
