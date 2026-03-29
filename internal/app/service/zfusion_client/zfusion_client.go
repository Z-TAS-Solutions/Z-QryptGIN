package zfusion_client

import (
	"context"
	"log"
	"runtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zfusionproto"
)

func RunZFusionClient(ip string) (zfusionproto.FusionCaptureServiceClient, error) {
	var target string

	if ip != "" {
		target = ip
	} else {
		if runtime.GOOS == "windows" {
			target = "127.0.0.1:50051"
		} else {
			target = "unix:///tmp/zfusion.sock"
		}
	}

	log.Printf("[ZTAS] Attempting to connect to ZFusionCore via %s (%s detected)...", target, runtime.GOOS)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		zFusionConn, err := grpc.DialContext(ctx, target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)

		if err != nil {
			log.Printf("[FusionClient] Connection failed: %v. Retrying in 7 seconds...", err)
			cancel()
			time.Sleep(2 * time.Second)
			continue
		}

		cancel()
		log.Printf("[FusionClient] Successfully connected to ZCrypt-FusionEngine at %s", target)

		zFusionClient := zfusionproto.NewFusionCaptureServiceClient(zFusionConn)
		return zFusionClient, nil
	}
}
