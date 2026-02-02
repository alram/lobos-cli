# Lobos cli

This is a quick gRPC client implementation for Lobos in Go, the protobufs are stored in the repo for ease of use.

After cloning:
```bash
$ go build 
$ ./lobosctl -h
NAME:
   lobosctl - Interact with lobos control plane

USAGE:
   lobosctl [global options] [command [command options]]

COMMANDS:
   users    Add/Remove/List S3 users
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --endpoint string  gRPC server endpoint (default: "127.0.0.1:50051")
   --help, -h         show help

```

Creating a user:
```bash
$ ./lobosctl users add --name lobos
User lobos created. Access: LB9BQU07HFHG7WF1GBDA, secret: +nyiuICCGh07rlbbLc+qvtRSF8QxmSFjnk4gLVrP
```


You can generate the Go bindings from the main lobos project via:

```bash
protoc -I ../lobos/src/controlplane \       
       --go_out=./protos --go_opt=paths=source_relative \
       --go-grpc_out=./protos --go-grpc_opt=paths=source_relative \
       ../lobos/src/controlplane/loboscontrol.proto
```
