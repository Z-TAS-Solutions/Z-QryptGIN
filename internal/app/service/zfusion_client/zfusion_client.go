package zfusion_client

import (
	"context"
	"log"
	"runtime"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zfusionproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func RunZFusionClient(ip string) (zfusionproto.FusionCaptureServiceClient, *grpc.ClientConn, error) {
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

	kpc := keepalive.ClientParameters{
		Time:                30 * time.Second,
		Timeout:             5 * time.Second,
		PermitWithoutStream: true,
	}

	log.Printf("[ZTAS] Attempting to connect to ZFusionCore via %s (%s detected)...", target, runtime.GOOS)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		zFusionConn, err := grpc.DialContext(ctx, target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(kpc),
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
		return zFusionClient, zFusionConn, nil
	}
}
