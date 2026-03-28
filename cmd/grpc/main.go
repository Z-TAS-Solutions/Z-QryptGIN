package main

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zfusion_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zpi_client"
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

	session_count := 0

	log.Println("Dialing ZIPC Crypt Core...")
	zipcClient, err := ipc.DialIPC()
	if err != nil {
		log.Println("Cannot initialize IPC dialer: %v", err)
	}
	defer zipcClient.Close()

	log.Println("Dialing ZPiScanner...")
	zpiClient, err := zpi_client.RunZPiClient("192.168.1.229:50051")
	if err != nil {
		log.Println("Cannot Connect To ZPiScanner: %v", err)
	}
	_, errr = zpi_client.InitializeZPiClient(zpiClient, 320)
	tofEventStream, err := zpi_client.StartToFStream(zpiClient)

	log.Println("Dialing ZFusionCore...")
	zfusionClient, errrr := zfusion_client.RunZFusionClient("")

	func() {
		for {
			evt, err := tofEventStream.Recv()
			if err != nil {
				log.Println("[Orchestrator] ToF Stream lost:", err)
				return
			}

			if evt.Type == zscanproto.ToFEvent_TRIGGER {
				log.Println("ToF trigger received! Initializing ZFusion...")
				session_count++

				ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)

				zFusionStream, err := zfusionClient.FusionCapture(ctx, &zfusionproto.ZFusionRequest{
					SessionId: strconv.Itoa(session_count),
				})

				if err != nil {
					log.Println("[Error] Could not start Fusion Capture:", err)
					cancel()
					continue
				}

				for {
					zfusionResponse, err := zFusionStream.Recv()
					if err == io.EOF {
						log.Println("[ZFusion] Hardware closed session normally.")
						break
					}
					if err != nil {
						log.Printf("[Error] Fusion Stream interrupted: %v", err)
						break
					}

					switch zfusionResponse.CompletionPhase {
					case zfusionproto.ZFusionResponse_PHASE_ROI:
						log.Println("[ZFusion] ROI Phase: Hand Detected.")
						tofEventStream.Send(&zscanproto.ToFEvent{Type: zscanproto.ToFEvent_PENDING})

					case zfusionproto.ZFusionResponse_PHASE_FUSION:
						log.Println("[ZFusion] Fusion Phase: Extracting Bitstream...")

						if zfusionResponse.StatusMessage == "ROI_FAIL" {
							log.Println("[Warning] Hardware failed to lock ROI.")
							break
						}

						if zipcClient != nil {
							matchCtx, matchCancel := context.WithTimeout(context.Background(), 5*time.Second)
							matchResult, score, err := zipcClient.MatchTemplate(matchCtx, "zischl", zfusionResponse.FusionBitstream)
							matchCancel()

							var ledStatus zscanproto.LEDStatus
							if err != nil || !matchResult {
								log.Printf("Match Denied (Score: %f, Error: %v)", score, err)
								ledStatus = zscanproto.LEDStatus_FAILED
							} else {
								log.Printf("Match SUCCESS (Score: %f)", score)
								ledStatus = zscanproto.LEDStatus_SUCCESS
							}

							tofEventStream.Send(&zscanproto.ToFEvent{
								Type:      zscanproto.ToFEvent_RESULT,
								LedStatus: ledStatus,
							})
						}
					}
				}
				cancel()
			}
		}
	}()

}
