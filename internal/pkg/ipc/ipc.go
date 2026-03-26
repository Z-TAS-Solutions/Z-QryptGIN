package ipc

import (
	"context"
	"log"
	"net"
	"os"
	"runtime"

	"github.com/Microsoft/go-winio"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zproto"
	"google.golang.org/grpc"
)

type PingServer struct {
	zproto.UnimplementedPingServiceServer
}

func (s *PingServer) Ping(ctx context.Context, in *zproto.PingRequest) (*zproto.PingResponse, error) {
	log.Printf("Received via gRPC-IPC: %s", in.GetMessage())
	return &zproto.PingResponse{Reply: "Z-Qrypt Hub Active"}, nil
}

func GetIPCListener() (net.Listener, error) {
	if runtime.GOOS == "windows" {
		return winio.ListenPipe(`\\.\pipe\zproto`, nil)
	}

	socketPath := "/tmp/zproto.sock"
	os.Remove(socketPath)
	return net.Listen("unix", socketPath)
}

func RunZIPCHub(grpcServer *grpc.Server) {

	zproto.RegisterPingServiceServer(grpcServer, &PingServer{})

	socketPath, error := GetIPCListener()
	if error != nil {
		log.Fatalf("Failed to start IPC: %v", error)
	}

	log.Println("Z-Hub is running on local IPC...")

	if error := grpcServer.Serve(socketPath); error != nil {
		log.Fatalf("gRPC server failed: %v", error)
	}

}
