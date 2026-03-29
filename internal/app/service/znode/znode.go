package znode

import (
	"context"
	"log"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
)

type ZCoreNode struct {
	zcoreproto.UnimplementedZCoreServiceServer
	EventChannel chan zcore.ZEvent
}

func (s *ZCoreNode) Ping(ctx context.Context, req *zcoreproto.PingRequest) (*zcoreproto.PingResponse, error) {
	log.Printf("Node received Ping from Hub: %s", req.Message)
	return &zcoreproto.PingResponse{Reply: "Node is active!"}, nil
}
