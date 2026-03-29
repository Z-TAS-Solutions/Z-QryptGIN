package main

import (
	"log"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/znode_remote"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"google.golang.org/grpc"
)

func RunRemoteGRPC(compute bool, remoteAddr string) {
	if !compute && remoteAddr == "" {
		log.Fatal("invalid config: no compute available")
	}

	grpcServer := grpc.NewServer()

	if compute {
		ipc.RunZIPCHub(grpcServer)
	}

	//localServer, err := grpc.Dial(remoteAddr, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("failed to dial remote: %v", err)
	//}

}

func main() {
	//nodeID := "ZTAS@0001"
	//nodeAddr := "localhost:50052"
	//hubAddr := "localhost:50051"

	znode_remote.RunZCoreRemote()
	//node_engine.RunZCoreWHub(nodeID, nodeAddr, hubAddr)
}
