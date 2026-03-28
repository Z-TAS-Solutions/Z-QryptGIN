package main

import (
	"log"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/grpc/zcore/zcore_node"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zfusion_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zpi_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"google.golang.org/grpc"
)

func RunZCoreWHub() {

	nodeID := "ZTAS@0001"
	nodeAddr := "localhost:50052"
	hubAddr := "localhost:50051"

	eventQueue := make(chan zcore.ZEvent, 200)

	log.Println("Dialing Remote ZCoreHub...")
	zCoreHubClient, err := zcore_node.ConnectZCoreHub(nodeID, nodeAddr, hubAddr, true)
	if err != nil {
		log.Println("Cannot Connect to Remote ZCoreHub: %v", err)
	}

	zcore_node.RunZCoreNode(nodeAddr, eventQueue)

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
		ZCoreHub:   zCoreHubClient,
	}

	ZCoreService.ZCoreEngine()

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

}
