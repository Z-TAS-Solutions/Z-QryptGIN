package znode_remote

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
)

type ZCoreHub struct {
	zcoreproto.UnimplementedZCoreServiceServer
	mu    sync.RWMutex
	nodes map[string]string
}

func (s *ZCoreHub) Ping(ctx context.Context, req *zcoreproto.PingRequest) (*zcoreproto.PingResponse, error) {
	log.Printf("Received ping from peer: %s", req.Message)
	return &zcoreproto.PingResponse{Reply: "Hub received: " + req.Message}, nil
}

func (s *ZCoreHub) Register(ctx context.Context, req *zcoreproto.RegisterRequest) (*zcoreproto.RegisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nodes[req.NodeId] = req.NodeAddr
	log.Printf("[ZCoreHub] Registered Peer: %s at %s", req.NodeId, req.NodeAddr)

	return &zcoreproto.RegisterResponse{Success: true, Message: "[ZCoreHub] Registered on Hub"}, nil
}

func RunZCoreRemote() {
	for {
		log.Println("[ZCoreHub] Attempting to start Z-Qrypt GRPC Server on :50051...")

		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Printf("[ZCoreHub] Failed to listen: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		opts := []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				recovery.UnaryServerInterceptor(),
			),
		}

		zcoreprotoHub := grpc.NewServer(opts...)

		zcoreproto.RegisterZCoreServiceServer(zcoreprotoHub, &ZCoreHub{
			nodes: make(map[string]string),
		})

		log.Println("[ZCoreHub] Z-Qrypt GRPC Server Running...")

		if err := zcoreprotoHub.Serve(listener); err != nil {
			log.Printf("[ZCoreHub] Server crashed or stopped: %v", err)
		}

		log.Println("[ZCoreHub] Restarting server in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
