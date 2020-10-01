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

Next, run the server:
```bash
$ ./rexd
```

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


## Design
[Design document](https://docs.google.com/document/d/1ICGf0mDO4sh1-PH73gvYQXFNxD0CGETNpy9wnxx1UWM/edit?usp=sharing)

## License
MIT
