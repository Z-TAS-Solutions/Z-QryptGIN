package main

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	go ipc.RunZIPCHub(grpcServer)

}
