# Rex: Remote process execution service

## build
Install [gRPC toolkit](https://grpc.io/docs/languages/go/quickstart/) for go.
Afterwards you can build the client and server binaries:
```bash
$ make rex rexd
```

Visit `scripts/README.md` to generate certificates. You can also use
pre-generated certificates in `./fixtures/tls/{client,server,ca}`.

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
    -policy '{"Principal": "*", "Action": "*", "Effect": "Allow"}' \
    -policy '{"Principal": "'$CL2_ID'", "Action": "/Rex/Exec", "Effect": "Deny"}'
```
The `-datadir` flag specifies the directory that rex will use to store
process stdout/stderr.
The `-policy` flag can be passed multiple times.
The former policy allows all
API calls by all users, otherwise no user is authorized to access any API.
The latter disallows a user with UUID equal to `$CL2_ID` from calling `/Rex/Exec`.

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

Since CL2 is not allowed to call Exec, the following command would fail:
```bash
$ ./rex $CL2_ARGS exec touch another_file
```

## Design
[Design document](https://docs.google.com/document/d/1ICGf0mDO4sh1-PH73gvYQXFNxD0CGETNpy9wnxx1UWM/edit?usp=sharing)

## License
MIT
