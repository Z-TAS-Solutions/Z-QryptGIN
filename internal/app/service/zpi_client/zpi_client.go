package zpi_client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zscanproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func RunZPiClient(ip string) (zscanproto.ZPiControllerClient, *grpc.ClientConn, error) {
	log.Printf("[ZPi] Attempting to connect to ZPi Controller at %s...", ip)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		kpc := keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}

		zPiConn, err := grpc.DialContext(ctx, ip,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(kpc),
		)

		if err != nil {
			log.Printf("[ZPi] Failed to connect to %s: %v. Retrying in 7 seconds...", ip, err)
			cancel()
			time.Sleep(2 * time.Second)
			continue
		}

		cancel()
		log.Printf("[ZPi] Successfully connected to ZPi GRPC Host: %s", ip)

		zPiClient := zscanproto.NewZPiControllerClient(zPiConn)
		return zPiClient, zPiConn, nil
	}
}

func InitializeZPiClient(zPiClient zscanproto.ZPiControllerClient, threshold uint32) (*zscanproto.Status, error) {

	configCtx, configCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer configCancel()

	tofConfigResp, err := zPiClient.ConfigureToF(configCtx, &zscanproto.ToFConfig{
		Threshold: uint32(threshold),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to configure ToF: %w", err)
	}

	log.Printf("ToF Configured Successfully: %s", tofConfigResp.GetMessage())
	return tofConfigResp, nil

}

func StartToFStream(zPiClient zscanproto.ZPiControllerClient) (zscanproto.ZPiController_ToFEventStreamClient, error) {
	ToFEventStream, err := zPiClient.ToFEventStream(context.Background())
	if err != nil {
		log.Fatalf("Failed To Start ToF Event Stream: %v", err)
		return nil, err
	}

	return ToFEventStream, nil
}
