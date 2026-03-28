package zcore_node

import (
	"context"
	"log"
	"net"

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

func main() {
	nodeID := "ZTAS@00001"
	nodeAddr := "localhost:50052"
	hubAddr := "localhost:50051"

	go func() {
		zlistener, _ := net.Listen("tcp", nodeAddr)
		s := grpc.NewServer()
		zcoreproto.RegisterZCoreServiceServer(s, &ZCoreNode{})
		log.Printf("Node Server listening on %s", nodeAddr)
		s.Serve(zlistener)
	}()

	conn, _ := grpc.Dial(hubAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := zcoreproto.NewZCoreServiceClient(conn)

	_, err := client.Register(context.Background(), &zcoreproto.RegisterRequest{
		NodeId:   nodeID,
		NodeAddr: nodeAddr,
	})
	if err != nil {
		log.Fatalf("Could not register: %v", err)
	}

	select {}
}
