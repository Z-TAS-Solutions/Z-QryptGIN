package znode_engine

import (
	"log"
	"net"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zclient_handler"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/znode_controller"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/znode_remote"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zscanproto"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
)

func RunZCoreWHub(nodeID, nodeAddr, hubAddr string) {

	eventQueue := make(chan zcore.ZEvent, 200)

	go znode_controller.RunZCoreNode(nodeAddr, eventQueue)

	ZCoreService := &zcore.ZCoreService{}

	zclient_handler.RunZClientHandlerEx(ZCoreService, nodeID, nodeAddr, hubAddr)

	go ZCoreService.ZCoreEngine(eventQueue)

	for {
		select {
		case event := <-eventQueue:
			log.Printf("Received gRPC event: %v", event.Type)

			switch event.Type {
			case zcore.EventType(0):
				stream, ok := event.Payload.(zscanproto.ZPiController_ToFEventStreamClient)
				if !ok {
					log.Println("Error: Payload was not a ToF stream client!")
					continue
				}

				ZCoreService.Mu.Lock()
				isEnrollment := ZCoreService.EnrollmentPending
				userID := ZCoreService.PendingUserID
				if isEnrollment {
					ZCoreService.EnrollmentPending = false
				}
				ZCoreService.Mu.Unlock()

				if isEnrollment {
					go ZCoreService.HandleEnrollSession(userID, stream)
				} else {
					go ZCoreService.HandleFusionSession(stream)
				}

			case zcore.EventType(1):
				log.Println("Sensor lost! Attempting recovery...")

			case zcore.EventType(2):
				userID, ok := event.Payload.(string)
				if !ok {
					log.Println("Error: Enrollment payload was not a string!")
					continue
				}
				log.Printf("Received Enrollment Command for User: %s", userID)

				ZCoreService.Mu.Lock()
				ZCoreService.EnrollmentPending = true
				ZCoreService.PendingUserID = userID
				ZCoreService.Mu.Unlock()
			}

		default:
			time.Sleep(time.Microsecond * 10)
		}
	}

}

func RunZCoreRemote() {
	for {
		log.Println("[ZCoreHub] Attempting to start Z-Qrypt GRPC Server on :50051...")

		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Printf("[ZCoreHub] Failed to listen: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		opts := []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				recovery.UnaryServerInterceptor(),
			),
		}

		zcoreprotoHub := grpc.NewServer(opts...)

		hub := &znode_remote.ZCoreHub{
			Nodes: make(map[string]string),
		}
		zcoreproto.RegisterZCoreServiceServer(zcoreprotoHub, hub)
		zcoreproto.RegisterZNodeControllerServer(zcoreprotoHub, hub)

		log.Println("[ZCoreHub] Z-Qrypt GRPC Server Running...")

		if err := zcoreprotoHub.Serve(listener); err != nil {
			log.Printf("[ZCoreHub] Server crashed or stopped: %v", err)
		}

		log.Println("[ZCoreHub] Restarting server in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
