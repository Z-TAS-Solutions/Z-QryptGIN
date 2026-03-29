package zclient_handler

import (
	"log"
	"sync"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zfusion_client"
	znodecontroller "github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/znode_controller"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zpi_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
)

func RunZClientHandlerEx(ZCoreService *zcore.ZCoreService, nodeID, nodeAddr, hubAddr string) {
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
		defer wg.Done()

		log.Println("Dialing Remote ZCoreHub...")
		zCoreHubClient, err := znodecontroller.ConnectZCoreHub(nodeID, nodeAddr, hubAddr, true)
		if err != nil {
			log.Println("Cannot Connect to Remote ZCoreHub: %v", err)
		}

		ZCoreService.ZCoreHub = zCoreHubClient

	}()

	log.Println("Awaiting All Autobots...")
	wg.Wait()
	log.Println("All services connected!")

}

func RunZClientHandler(ZCoreService *zcore.ZCoreService) {
	var wg sync.WaitGroup

	wg.Add(3)

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

	log.Println("Awaiting All Autobots...")
	wg.Wait()
	log.Println("All services connected!")

}
