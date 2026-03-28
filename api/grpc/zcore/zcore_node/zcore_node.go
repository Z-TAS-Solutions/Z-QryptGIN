package zcore_node

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ZCoreNode struct {
	zcoreproto.UnimplementedZCoreServiceServer
	EventChannel chan zcore.ZEvent
}

func (s *ZCoreNode) Ping(ctx context.Context, req *zcoreproto.PingRequest) (*zcoreproto.PingResponse, error) {
	log.Printf("Node received Ping from Hub: %s", req.Message)
	return &zcoreproto.PingResponse{Reply: "Node is active!"}, nil
}

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
		zcoreproto.RegisterZCoreServiceServer(nodeServer, &ZCoreNode{EventChannel: eventChannel})

		log.Printf("[ZCoreNode] Node Server listening on %s", nodeAddr)
		if err := nodeServer.Serve(zlistener); err != nil {
			log.Printf("[ZCoreNode] Server exited: %v", err)
		}

		if err != nil {
			log.Printf("[ZCoreNodeManager] Server crashed: %v. Restarting in 3s...", err)
			time.Sleep(3 * time.Second)
			continue
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
