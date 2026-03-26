package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/grpc/zpi_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zscanproto"
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
	client := zpi_client.RunZPiClient("192.168.1.229:50051")

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	//defer cancel()

	statusResp, err := client.SetLEDStatus(context.Background(), &zscanproto.LEDStatusRequest{
		Status: zscanproto.LEDStatus_PENDING,
	})

	if err != nil {
		log.Fatalf("RPC failed: %v", err)
	}

	fmt.Println("SetLEDStatus:", statusResp.Message)

}
