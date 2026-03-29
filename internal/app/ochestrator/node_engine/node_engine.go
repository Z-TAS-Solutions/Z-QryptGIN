package node_engine

import (
	"log"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zclient_handler"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zcore"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zscanproto"
)

func RunZCoreWHub(nodeID, nodeAddr, hubAddr string) {

	eventQueue := make(chan zcore.ZEvent, 200)

	//go znode_controller.RunZCoreNode(nodeAddr, eventQueue)

	ZCoreService := &zcore.ZCoreService{}

	zclient_handler.RunZClientHandler(ZCoreService)

	go ZCoreService.ZCoreEngine(eventQueue)

	for {
		select {
		case event := <-eventQueue:
			log.Printf("Received gRPC event: %s", event.Type)

			switch event.Type {
			case zcore.EventType(0):
				stream, ok := event.Payload.(zscanproto.ZPiController_ToFEventStreamClient)
				if !ok {
					log.Println("Error: Payload was not a ToF stream client!")
					continue
				}

				go ZCoreService.HandleFusionSession(stream)

			case zcore.EventType(1):
				log.Println("Sensor lost! Attempting recovery...")

			case zcore.EventType(2):
				log.Println("Received Command From ZCoreHub")
			}

		default:
			time.Sleep(time.Microsecond * 10)
		}
	}

}
