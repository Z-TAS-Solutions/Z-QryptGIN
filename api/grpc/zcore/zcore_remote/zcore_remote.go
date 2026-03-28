package zcore_remote

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/recovery"
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
	log.Printf("Registered Peer: %s at %s", req.NodeId, req.NodeAddr)

	return &zcoreproto.RegisterResponse{Success: true, Message: "Registered on Hub"}, nil
}

func RunZCoreRemote() {
	listener, _ := net.Listen("tcp", ":50051")

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
		),
	}

	zcoreprotoHub := grpc.NewServer(opts...)

	zcoreproto.RegisterZCoreServiceServer(zcoreprotoHub, &ZCoreHub{
		nodes: make(map[string]string),
	})
	log.Println("Z-Qrypt GRPC Server Starting On :50051...")

	if err := zcoreprotoHub.Serve(listener); err != nil {
		log.Fatalf("[ZCoreHub] Server failed to serve: %v", err)
	}
}
