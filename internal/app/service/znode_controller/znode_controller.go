package znode_controller

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/znode"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	for {
		zCoreHubConn, err := grpc.NewClient(hubAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err == nil {
			remoteNode := zcoreproto.NewZCoreServiceClient(zCoreHubConn)

			_, err = remoteNode.Register(context.Background(), &zcoreproto.RegisterRequest{
				NodeId:   nodeID,
				NodeAddr: nodeAddr,
			})

			if err == nil {
				log.Println("[ZCoreNode] Successfully registered with ZCoreHub!")
				return remoteNode, nil
			}
		}

		if !retry {
			log.Fatalf("[ZCoreNode] Failed to connect to ZCoreHub and retry is disabled: %v", err)
			return nil, err
		}

		log.Printf("[ZCoreNode] ZCoreHub unreachable. Retrying in 5s...")
		time.Sleep(5 * time.Second)
	}
}
