# Rex: Remote process execution service

# build
Install [gRPC toolkit](https://grpc.io/docs/languages/go/quickstart/) for go.
Afterwards you can build the client and server binaries:
```bash
$ make rex rexd
```

To generate the certificates (opnessl required):
```bash
$ (cd scripts && ./gen_certs.sh)
```

The certificates will be created in the scripts directory:
```bash
$ ls scripts
```

Next, run the server, passing the required policies:
```bash
$ ./rexd -policy '{"Principal": "*", "Action": "*", "Effect": "Allow"}' \
    -policy '{"Principal": "USER_UUID_HERE", "Action": "/Rex/Exec", "Effect": "Deny"}'
```
The `-policy` flag can be passed multiple times.
The former policy allows all
API calls by all users, otherwise no user is authorized to access any API.
The latter disallows a particular user from calling `/Rex/Exec`.

Then on another terminal:
```bash
$ ./rex exec $YOUR_COMMAND_HERE
```

For example:
```bash
$ ./rex exec touch some_file
```

Verify the results (on the same directory that you ran `./rexd` from):
```bash
$ ls some_file
```

Executing a nonexistent file, which allows the client to catch `ErrNotFound`:
```bash
$ ./rex exec nonexistent-binary
```

Executing a file without the execute permission:
```bash
$ ./rex exec ./rex.go
```


## Design
[Design document](https://docs.google.com/document/d/1ICGf0mDO4sh1-PH73gvYQXFNxD0CGETNpy9wnxx1UWM/edit?usp=sharing)

## License
MIT
