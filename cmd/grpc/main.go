package main

import (
	"log"
	"net"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zfusion_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zpi_client"
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

	log.Println("Dialing ZIPC Crypt Core...")
	zipcClient, err := ipc.DialIPC()
	if err != nil {
		log.Println("Cannot initialize IPC dialer: %v", err)
	}
	defer zipcClient.Close()

	log.Println("Dialing ZPiScanner...")
	zpiClient, err := zpi_client.RunZPiClient("192.168.1.229:50051")
	if err != nil {
		log.Println("Cannot Connect To ZPiScanner: %v", err)
	}
	_, err = zpi_client.InitializeZPiClient(zpiClient, 320)
	if err != nil {
		log.Println("Failed To Configure ZPiClient..")
	}

	log.Println("Dialing ZFusionCore...")
	zfusionClient, err := zfusion_client.RunZFusionClient("")
	if err != nil {
		log.Println("Failed To Connect To ZFusion Core !")
	}

	ZCoreService := &zcore.ZCoreService{
		ZIPCClient: zipcClient,
		ZPiClient:  zpiClient,
		ZFusion:    zfusionClient,
	}

	ZCoreService.ZCoreEngine()

}
