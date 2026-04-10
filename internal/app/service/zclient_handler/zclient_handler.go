package zclient_handler

import (
	"context"
	"log"
	"sync"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zfusion_client"
	znodecontroller "github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/znode_controller"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/znode_monitor"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zpi_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
)

// this premium version includes the remote hub as a service while the basic version doesn't
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

	// setting up context for service monitor
	ctx, EndMonitor := context.WithCancel(context.Background())

	ZCoreService.ServiceMonitor.MonitorContext = ctx
	ZCoreService.ServiceMonitor.EndMonitor = EndMonitor

}

// no remote hub meaning no remote admin panel/ 2FA
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

	// setting up context for service monitor
	ctx, EndMonitor := context.WithCancel(context.Background())

	ZCoreService.ServiceMonitor.MonitorContext = ctx
	ZCoreService.ServiceMonitor.EndMonitor = EndMonitor

}

func connectZPi(ZCoreService *zcore.ZCoreService) {
	log.Println("[ZClientHandler] Dialing ZPiScanner...")
	zpiClient, zpiConn, err := zpi_client.RunZPiClient("192.168.1.229:50051")
	if err != nil {
		log.Printf("[ZClientHandler] Initial ZPi connection failed: %v", err)
		return
	}
	_, err = zpi_client.InitializeZPiClient(zpiClient, 320)
	if err != nil {
		log.Println("[ZClientHandler] Initial ZPi configuration failed.")
	}
	ZCoreService.ZPiClient = zpiClient

	// Starting the ZNode monitor
	znodemonitor.RunNodeMonitor(ZCoreService.ServiceMonitor.MonitorContext, "[ZPi]", zpiConn)

}

func connectZFusion(ZCoreService *zcore.ZCoreService) {
	log.Println("[ZClientHandler] Dialing ZFusionCore...")
	zfusionClient, zfusionConn, err := zfusion_client.RunZFusionClient("")
	if err != nil {
		log.Printf("[ZClientHandler] Initial ZFusion connection failed: %v", err)
		return
	}
	ZCoreService.ZFusion = zfusionClient

	// Starting the ZNode monitor
	znodemonitor.RunNodeMonitor(ZCoreService.ServiceMonitor.MonitorContext, "[ZFusion]", zfusionConn)

}

func connectZIPC(ZCoreService *zcore.ZCoreService) {
	log.Println("[ZClientHandler] Dialing ZIPC Crypt Core...")
	zIPCClient, err := ipc.RunZIPCClient()
	if err != nil {
		log.Printf("[ZClientHandler] Initial ZIPC connection failed: %v", err)
		return
	}
	ZCoreService.ZIPCClient = zIPCClient

	// Starting the ZNode monitor
	znodemonitor.RunNodeMonitor(ZCoreService.ServiceMonitor.MonitorContext, "[ZIPC]", zIPCClient.ZIPCConn)

}

func connectZCoreHub(ZCoreService *zcore.ZCoreService, nodeID, nodeAddr, hubAddr string) {
	log.Println("[ZClientHandler] Dialing Remote ZCoreHub...")
	zCoreHubClient, zCoreHubConn, err := znodecontroller.ConnectZCoreHub(nodeID, nodeAddr, hubAddr, true)
	if err != nil {
		log.Printf("[ZClientHandler] Initial ZCoreHub connection/registration failed: %v", err)
		return
	}
	ZCoreService.ZCoreHub = zCoreHubClient

	// Starting the ZNode monitor
	znodemonitor.RunNodeMonitor(ZCoreService.ServiceMonitor.MonitorContext, "[ZIPC]", zCoreHubConn)

}
