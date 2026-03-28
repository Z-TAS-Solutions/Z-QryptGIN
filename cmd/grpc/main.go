package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/api/grpc/zpi_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zfusionproto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zscanproto"
	"google.golang.org/grpc"
)

func RunLocalGRPC(compute bool) {

	grpcServer := grpc.NewServer()

	if compute {
		ipc.RunZIPCHub(grpcServer)
	}

	localServer, error := net.Listen("tcp", ":50051")
	if error != nil {
		log.Fatalf("failed to listen: %v", error)
	}

	grpcServer.Serve(localServer)

}

func RunRemoteGRPC(compute bool, remoteAddr string) {
	if !compute && remoteAddr == "" {
		log.Fatal("invalid config: no compute available")
	}

	grpcServer := grpc.NewServer()

	if compute {
		ipc.RunZIPCHub(grpcServer)
	}

	//localServer, err := grpc.Dial(remoteAddr, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("failed to dial remote: %v", err)
	//}

}

func main() {

	log.Println("Dialing ZIPC Crypt Core...")
	zipcClient, err := ipc.DialIPC()
	if err != nil {
		log.Fatalf("Cannot initialize IPC dialer: %v", err)
	}
	defer zipcClient.Close()

	log.Println("Dialing ZPiScanner...")
	zpiClient, err := zpi_client.RunZPiClient("192.168.1.229:50051")
	if err != nil {
		log.Fatalf("Cannot Connect To ZPiScanner: %v", err)
	}
	_, error := zpi_client.InitializeZPiClient(zpiClient, 320)
	tofEventStream, err := zpi_client.StartToFStream(zpiClient)

	log.Println("Dialing ZFusionCore...")
	zfusionClient := zfusionproto.NewFusionCaptureServiceClient()

	func() {
		for {
			evt, err := tofEventStream.Recv()
			if err != nil {
				log.Println("Server disconnected or error:", err)
				return
			}

			if evt.Type == zscanproto.ToFEvent_TRIGGER {
				log.Println("ToF trigger received!")

				if zipcClient != nil {
					log.Println("Sending dummy match request to Rust IPC...")
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					isMatch, score, err := zipcClient.MatchTemplate(ctx, "zischl", []byte("bleh..."))
					cancel()

					if err != nil {
						log.Printf("IPC MatchTemplate Error: %v\n", err)
					} else {
						log.Printf("IPC MatchTemplate Result -> match=%v, score=%f\n", isMatch, score)
					}
				}

				pendingState := &zscanproto.ToFEvent{
					Type: zscanproto.ToFEvent_PENDING,
				}
				if err := tofEventStream.Send(pendingState); err != nil {
					log.Println("Failed to send pending state:", err)
					continue
				}

				time.Sleep(3 * time.Second)

				response := &zscanproto.ToFEvent{
					Type:      zscanproto.ToFEvent_RESULT,
					LedStatus: zscanproto.LEDStatus_SUCCESS,
				}
				if err := tofEventStream.Send(response); err != nil {
					log.Println("Failed to send response:", err)
					continue
				}
				log.Println("Response sent !")
			}
		}
	}()

}
