# Rex: Remote process execution service

# build
Install [gRPC toolkit](https://grpc.io/docs/languages/go/quickstart/) for go.
Afterwards you can build the client and server binaries:
```bash
$ make rex rexd
```

Next, run the server:
```bash
$ ./rexd
```

Then on another terminal:
```bash
$ ./rex
```

## Design
[Design document](https://docs.google.com/document/d/1ICGf0mDO4sh1-PH73gvYQXFNxD0CGETNpy9wnxx1UWM/edit?usp=sharing)

## License
MIT
