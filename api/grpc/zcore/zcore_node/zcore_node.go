package zcore_node

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ZCoreNode struct {
	zcoreproto.UnimplementedZCoreServiceServer
}

func (s *ZCoreNode) Ping(ctx context.Context, req *zcoreproto.PingRequest) (*zcoreproto.PingResponse, error) {
	log.Printf("Node received Ping from Hub: %s", req.Message)
	return &zcoreproto.PingResponse{Reply: "Peer is active!"}, nil
}

func RunZCoreNode(nodeAddr string) {
	go func() {
		zlistener, err := net.Listen("tcp", nodeAddr)
		if err != nil {
			log.Printf("[ZCoreNode] Failed to listen: %v", err)
			return
		}

		s := grpc.NewServer()
		zcoreproto.RegisterZCoreServiceServer(s, &ZCoreNode{})

		log.Printf("[ZCoreNode] Node Server listening on %s", nodeAddr)
		if err := s.Serve(zlistener); err != nil {
			log.Printf("[ZCoreNode] Server exited: %v", err)
		}
	}()
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
				log.Println("[ZCoreNode] Successfully registered with Hub!")
				return remoteNode, nil
			}
		}

		if !retry {
			log.Fatalf("[ZCoreNode] Failed to connect and retry is disabled: %v", err)
			return nil, err
		}

		log.Printf("[ZCoreNode] Hub unreachable. Retrying in 5s...")
		time.Sleep(5 * time.Second)
	}
}
