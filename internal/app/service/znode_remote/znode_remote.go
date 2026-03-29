package znode_remote

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ZCoreHub struct {
	zcoreproto.UnimplementedZCoreServiceServer
	zcoreproto.UnimplementedZNodeControllerServer
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

func (s *ZCoreHub) Request2FA(ctx context.Context, req *zcoreproto.TwoFARequest) (*zcoreproto.TwoFAResponse, error) {
	log.Printf("[ZCoreHub] 2FA Triggered for User: %s", req.UserId)
	return &zcoreproto.TwoFAResponse{
		Success: true,
		Message: "2FA request received and logged",
	}, nil
}

func (s *ZCoreHub) StartEnrollment(ctx context.Context, req *zcoreproto.EnrollmentRequest) (*zcoreproto.EnrollmentResponse, error) {
	s.mu.RLock()
	nodes := make(map[string]string)
	for id, addr := range s.Nodes {
		nodes[id] = addr
	}
	s.mu.RUnlock()

	log.Printf("[ZCoreHub] Broadcasting Enrollment Request for UserID: %s to %d nodes", req.UserId, len(nodes))

	for id, addr := range nodes {
		go func(nodeID, nodeAddr string) {
			log.Printf("[ZCoreHub] Sending enrollment request to node %s at %s", nodeID, nodeAddr)

			conn, err := grpc.NewClient(nodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Printf("[ZCoreHub] Failed to connect to node %s: %v", nodeID, err)
				return
			}
			defer conn.Close()

			client := zcoreproto.NewZNodeControllerClient(conn)

			nodeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.StartEnrollment(nodeCtx, req)
			if err != nil {
				log.Printf("[ZCoreHub] Failed to start enrollment on node %s: %v", nodeID, err)
				return
			}

			log.Printf("[ZCoreHub] Node %s response: %s (Accepted: %v)", nodeID, resp.Message, resp.Accepted)
		}(id, addr)
	}

	return &zcoreproto.EnrollmentResponse{
		Accepted: true,
		Message:  "Enrollment broadcast initiated to all registered nodes",
	}, nil
}
