package znode_controller

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/znode"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc/keepalive"
)

func RunZCoreNode(nodeAddr string, eventChannel chan zcore.ZEvent) {
	for {
		log.Println("[ZCoreNodeManager] Starting ZCore gRPC Server...")

		zlistener, err := net.Listen("tcp", nodeAddr)
		if err != nil {
			log.Printf("[ZCoreNode] Failed to listen: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		nodeServer := grpc.NewServer()
		srv := &znode.ZCoreNode{EventChannel: eventChannel}
		zcoreproto.RegisterZCoreServiceServer(nodeServer, srv)
		zcoreproto.RegisterZNodeControllerServer(nodeServer, srv)

		log.Printf("[ZCoreNode] Node Server listening on %s", nodeAddr)
		if err := nodeServer.Serve(zlistener); err != nil {
			log.Printf("[ZCoreNode] Server exited: %v", err)
		}

		log.Println("[ZCoreNodeManager] Server shutting down...")
		break
	}
}

func ConnectZCoreHub(nodeID, nodeAddr, hubAddr string, retry bool) (zcoreproto.ZCoreServiceClient, error) {
	kpc := keepalive.ClientParameters{
		Time:                30 * time.Second,
		Timeout:             5 * time.Second,
		PermitWithoutStream: true,
	}

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kpc),
	}

	conn, err := grpc.NewClient(hubAddr, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	client := zcoreproto.NewZCoreServiceClient(conn)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		_, err := client.Register(ctx, &zcoreproto.RegisterRequest{
			NodeId:   nodeID,
			NodeAddr: nodeAddr,
		})
		cancel()

		if err == nil {
			log.Println("[ZCoreNode] Successfully registered with ZCoreHub!")
			return client, nil
		}

		if !retry {
			log.Fatalf("[ZCoreNode] Registration failed, retry disabled: %v", err)
			return nil, err
		}

		log.Printf("[ZCoreNode] Registration failed: %v. Retrying in 5s...", err)
		time.Sleep(5 * time.Second)
	}
}
