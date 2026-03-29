package znode_remote

import (
	"context"
	"log"
	"sync"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
)

type ZCoreHub struct {
	zcoreproto.UnimplementedZCoreServiceServer
	mu    sync.RWMutex
	Nodes map[string]string
}

func (s *ZCoreHub) Ping(ctx context.Context, req *zcoreproto.PingRequest) (*zcoreproto.PingResponse, error) {
	log.Printf("Received ping from peer: %s", req.Message)
	return &zcoreproto.PingResponse{Reply: "Hub received: " + req.Message}, nil
}

func (s *ZCoreHub) Register(ctx context.Context, req *zcoreproto.RegisterRequest) (*zcoreproto.RegisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Nodes[req.NodeId] = req.NodeAddr
	log.Printf("[ZCoreHub] Registered Peer: %s at %s", req.NodeId, req.NodeAddr)

	return &zcoreproto.RegisterResponse{Success: true, Message: "[ZCoreHub] Registered on Hub"}, nil
}
