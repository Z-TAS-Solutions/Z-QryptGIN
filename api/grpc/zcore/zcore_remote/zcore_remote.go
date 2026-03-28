package zcore_remote

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
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

	s.nodes[req.PeerId] = req.ListenAddr
	log.Printf("Registered Peer: %s at %s", req.PeerId, req.ListenAddr)

	return &zcoreproto.RegisterResponse{Success: true, Message: "Registered on Hub"}, nil
}

func main() {
	listener, _ := net.Listen("tcp", ":50051")
	zcoreprotoHub := grpc.NewServer()
	zcoreproto.RegisterZCoreServiceServer(zcoreprotoHub, &ZCoreHub{
		nodes: make(map[string]string),
	})
	log.Println("Z-Qrypt GRPC Server Starting On :50051...")

	zcoreprotoHub.Serve(listener)
}
