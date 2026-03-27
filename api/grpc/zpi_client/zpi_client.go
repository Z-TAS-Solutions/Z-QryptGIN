package zpi_client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zscanproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunZPiClient(ip string, ipcClient *ipc.ZPIPCClient) (zscanproto.ZPiControllerClient, zscanproto.ZPiController_ToFEventStreamClient) {
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	client := zscanproto.NewZPiControllerClient(conn)
	log.Print("Connected to ZPi GRPC Host!")

	tofConfigResp, _ := client.ConfigureToF(context.Background(), &zscanproto.ToFConfig{
		Threshold: 320,
	})
	fmt.Println("ConfigureToF:", tofConfigResp.Message)

	stream, err := client.ToFEventStream(context.Background())
	if err != nil {
		log.Fatalf("Failed To Start ToF Event Stream: %v", err)
	}

	func() {
		for {
			evt, err := stream.Recv()
			if err != nil {
				log.Println("Server disconnected or error:", err)
				return
			}

			if evt.Type == zscanproto.ToFEvent_TRIGGER {
				log.Println("ToF trigger received!")

				time.Sleep(2 * time.Second)

				if ipcClient != nil {
					log.Println("Sending dummy match request to Rust IPC...")
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					isMatch, score, err := ipcClient.MatchTemplate(ctx, "zischl", []byte("bleh..."))
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
				if err := stream.Send(pendingState); err != nil {
					log.Println("Failed to send pending state:", err)
					continue
				}

				time.Sleep(3 * time.Second)

				response := &zscanproto.ToFEvent{
					Type:      zscanproto.ToFEvent_RESULT,
					LedStatus: zscanproto.LEDStatus_SUCCESS,
				}
				if err := stream.Send(response); err != nil {
					log.Println("Failed to send response:", err)
					continue
				}
				log.Println("Response sent !")
			}
		}
	}()

	return client, stream
}
