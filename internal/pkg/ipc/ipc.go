package ipc

import (
	"log"
	"net"
	"os"
	"runtime"

	"github.com/Microsoft/go-winio"
	"google.golang.org/grpc"
)

func GetIPCListener() (net.Listener, error) {
	if runtime.GOOS == "windows" {
		return winio.ListenPipe(`\\.\pipe\zproto`, nil)
	}

	socketPath := "/tmp/zproto.sock"
	os.Remove(socketPath)
	return net.Listen("unix", socketPath)
}

func RunZIPCHub(grpcServer *grpc.Server) {

	socketPath, error := GetIPCListener()
	if error != nil {
		log.Fatalf("Failed to start IPC: %v", error)
	}

	log.Println("Z-Hub is running on local IPC...")

	if error := grpcServer.Serve(socketPath); error != nil {
		log.Fatalf("gRPC server failed: %v", error)
	}

}
