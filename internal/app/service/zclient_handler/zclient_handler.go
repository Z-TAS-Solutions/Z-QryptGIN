package zclient_handler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zfusion_client"
	znodecontroller "github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/znode_controller"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zpi_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zscanproto"
)

func RunZClientHandlerEx(ZCoreService *zcore.ZCoreService, nodeID, nodeAddr, hubAddr string) {
	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()
		connectZPi(ZCoreService)
	}()

	go func() {
		defer wg.Done()
		connectZFusion(ZCoreService)
	}()

	go func() {
		defer wg.Done()
		connectZIPC(ZCoreService)
	}()

	go func() {
		defer wg.Done()
		connectZCoreHub(ZCoreService, nodeID, nodeAddr, hubAddr)
	}()

	log.Println("[ZClientHandler] Awaiting initial connection of all services...")
	wg.Wait()
	log.Println("[ZClientHandler] All services connected! Starting persistent monitor...")

	// Start the batched health monitor
	go monitorConnections(ZCoreService, nodeID, nodeAddr, hubAddr, true)
}

func RunZClientHandler(ZCoreService *zcore.ZCoreService) {
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		connectZPi(ZCoreService)
	}()

	go func() {
		defer wg.Done()
		connectZFusion(ZCoreService)
	}()

	go func() {
		defer wg.Done()
		connectZIPC(ZCoreService)
	}()

	log.Println("[ZClientHandler] Awaiting initial connection of core services...")
	wg.Wait()
	log.Println("[ZClientHandler] Core services connected! Starting persistent monitor...")

	go monitorConnections(ZCoreService, "", "", "", false)
}

func connectZPi(ZCoreService *zcore.ZCoreService) {
	log.Println("[ZClientHandler] Dialing ZPiScanner...")
	zpiClient, err := zpi_client.RunZPiClient("192.168.1.229:50051")
	if err != nil {
		log.Printf("[ZClientHandler] Initial ZPi connection failed: %v", err)
		return
	}
	_, err = zpi_client.InitializeZPiClient(zpiClient, 320)
	if err != nil {
		log.Println("[ZClientHandler] Initial ZPi configuration failed.")
	}
	ZCoreService.ZPiClient = zpiClient
}

func connectZFusion(ZCoreService *zcore.ZCoreService) {
	log.Println("[ZClientHandler] Dialing ZFusionCore...")
	zfusionClient, err := zfusion_client.RunZFusionClient("")
	if err != nil {
		log.Printf("[ZClientHandler] Initial ZFusion connection failed: %v", err)
		return
	}
	ZCoreService.ZFusion = zfusionClient
}

func connectZIPC(ZCoreService *zcore.ZCoreService) {
	log.Println("[ZClientHandler] Dialing ZIPC Crypt Core...")
	zIPCClient, err := ipc.RunZIPCClient()
	if err != nil {
		log.Printf("[ZClientHandler] Initial ZIPC connection failed: %v", err)
		return
	}
	ZCoreService.ZIPCClient = zIPCClient
}

func connectZCoreHub(ZCoreService *zcore.ZCoreService, nodeID, nodeAddr, hubAddr string) {
	log.Println("[ZClientHandler] Dialing Remote ZCoreHub...")
	zCoreHubClient, err := znodecontroller.ConnectZCoreHub(nodeID, nodeAddr, hubAddr, true)
	if err != nil {
		log.Printf("[ZClientHandler] Initial ZCoreHub connection/registration failed: %v", err)
		return
	}
	ZCoreService.ZCoreHub = zCoreHubClient
}

func monitorConnections(ZCoreService *zcore.ZCoreService, nodeID, nodeAddr, hubAddr string, hasHub bool) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if ZCoreService.ZPiClient == nil {
			connectZPi(ZCoreService)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, err := ZCoreService.ZPiClient.ConfigureToF(ctx, &zscanproto.ToFConfig{Threshold: 320})
			cancel()
			if err != nil {
				log.Printf("[ConnectionMonitor] ZPiScanner appears offline: %v. Reconnecting...", err)
				connectZPi(ZCoreService)
			}
		}

		if ZCoreService.ZIPCClient == nil {
			connectZIPC(ZCoreService)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, err := ZCoreService.ZIPCClient.Ping(ctx, "keep-alive")
			cancel()
			if err != nil {
				log.Printf("[ConnectionMonitor] ZIPC Core appears offline: %v. Reconnecting...", err)
				connectZIPC(ZCoreService)
			}
		}

		if hasHub {
			if ZCoreService.ZCoreHub == nil {
				connectZCoreHub(ZCoreService, nodeID, nodeAddr, hubAddr)
			} else {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				_, err := ZCoreService.ZCoreHub.Ping(ctx, &zcoreproto.PingRequest{Message: "keep-alive"})
				cancel()
				if err != nil {
					log.Printf("[ConnectionMonitor] ZCoreHub appears offline: %v. Reconnecting and Re-registering...", err)
					connectZCoreHub(ZCoreService, nodeID, nodeAddr, hubAddr)
				}
			}
		}

		if ZCoreService.ZFusion == nil {
			connectZFusion(ZCoreService)
		}
	}
}
