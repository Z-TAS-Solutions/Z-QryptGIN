package main

import (
	"log"
	"sync"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/grpc/zcore/zcore_node"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zfusion_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zpi_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"google.golang.org/grpc"
)

func RunZClientHandler(ZCoreService *zcore.ZCoreService, nodeID, nodeAddr, hubAddr string) {
	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()

		log.Println("Dialing ZPiScanner...")
		zpiClient, err := zpi_client.RunZPiClient("192.168.1.229:50051")
		if err != nil {
			log.Println("Cannot Connect To ZPiScanner: %v", err)
		}
		_, err = zpi_client.InitializeZPiClient(zpiClient, 320)
		if err != nil {
			log.Println("Failed To Configure ZPiClient..")
		}

		ZCoreService.ZPiClient = zpiClient

	}()

	go func() {
		defer wg.Done()

		log.Println("Dialing ZFusionCore...")
		zfusionClient, err := zfusion_client.RunZFusionClient("")
		if err != nil {
			log.Println("Failed To Connect To ZFusion Core !")
		}
		ZCoreService.ZFusion = zfusionClient

	}()

	go func() {
		defer wg.Done()

		log.Println("Dialing ZIPC Crypt Core...")
		zIPCClient, err := ipc.RunZIPCClient()
		if err != nil {
			log.Println("Cannot Connect To ZIPC Crypt Core : %v", err)
		}

		ZCoreService.ZIPCClient = zIPCClient

	}()

	go func() {

		log.Println("Dialing Remote ZCoreHub...")
		zCoreHubClient, err := zcore_node.ConnectZCoreHub(nodeID, nodeAddr, hubAddr, true)
		if err != nil {
			log.Println("Cannot Connect to Remote ZCoreHub: %v", err)
		}

		ZCoreService.ZCoreHub = zCoreHubClient

	}()

}

func RunZCoreWHub() {

	nodeID := "ZTAS@0001"
	nodeAddr := "localhost:50052"
	hubAddr := "localhost:50051"

	eventQueue := make(chan zcore.ZEvent, 200)

	go zcore_node.RunZCoreNode(nodeAddr, eventQueue)

	ZCoreService := &zcore.ZCoreService{}

	RunZClientHandler(ZCoreService, nodeID, nodeAddr, hubAddr)

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
