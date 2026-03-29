package znode

import (
	"context"
	"log"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
)

type ZCoreNode struct {
	zcoreproto.UnimplementedZCoreServiceServer
	zcoreproto.UnimplementedZNodeControllerServer
	EventChannel chan zcore.ZEvent
}

func (s *ZCoreNode) StartEnrollment(ctx context.Context, req *zcoreproto.EnrollmentRequest) (*zcoreproto.EnrollmentResponse, error) {
	log.Printf("[ZCoreNode] StartEnrollment requested for UserID: %s", req.UserId)

	s.EventChannel <- zcore.ZEvent{
		Type:    zcore.EventType(2),
		Payload: req.UserId,
	}

	return &zcoreproto.EnrollmentResponse{
		Accepted: true,
		Message:  "Node enrollment initiated for user " + req.UserId,
	}, nil
}

func (s *ZCoreNode) Ping(ctx context.Context, req *zcoreproto.PingRequest) (*zcoreproto.PingResponse, error) {
	log.Printf("Node received Ping from Hub: %s", req.Message)
	return &zcoreproto.PingResponse{Reply: "Node is active!"}, nil
}
