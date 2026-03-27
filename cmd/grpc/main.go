package main

import (
	"log"
	"net"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/grpc/zpi_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"google.golang.org/grpc"
)

func RunLocalGRPC(compute bool) {

	grpcServer := grpc.NewServer()

	if compute {
		ipc.RunZIPCHub(grpcServer)
	}

	localServer, error := net.Listen("tcp", ":50051")
	if error != nil {
		log.Fatalf("failed to listen: %v", error)
	}

	grpcServer.Serve(localServer)

}

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
	//RunLocalGRPC(true)
	
	log.Println("Dialing Rust IPC Service...")
	zpipcClient, err := ipc.DialIPC()
	if err != nil {
		log.Fatalf("Cannot initialize IPC dialer: %v", err)
	}
	defer zpipcClient.Close()

	zpi_client.RunZPiClient("192.168.1.229:50051", zpipcClient)
}
