package zcore

import (
	"context"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zpi_client"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zfusionproto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zscanproto"
)

type ZCoreService struct {
	ZIPCClient *ipc.ZPIPCClient
	ZPiClient  zscanproto.ZPiControllerClient
	ZFusion    zfusionproto.FusionCaptureServiceClient

	sessionCount int
}

func (z *ZCoreService) HandleFusionSession(tofEventStream zscanproto.ZPiController_ToFEventStreamClient) {
	log.Println("[ZCore] Initializing ZFusion...")

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	zFusionStream, err := z.ZFusion.FusionCapture(ctx, &zfusionproto.ZFusionRequest{
		SessionId: strconv.Itoa(z.sessionCount),
	})

	if err != nil {
		log.Println("[Error] Could not start Fusion Capture:", err)
		cancel()
		return
	}

	for {
		zFusionResponse, err := zFusionStream.Recv()
		if err == io.EOF {
			log.Println("[ZFusion] Session Closed.")
			break
		}
		if err != nil {
			log.Printf("[Error] Fusion Stream interrupted: %v", err)
			break
		}

		log.Println("[ZFusion] Fusion :", zFusionResponse.StatusMessage)

		switch zFusionResponse.CompletionPhase {
		case zfusionproto.ZFusionResponse_PHASE_ROI:
			log.Println("[ZFusion] ROI Phase: Hand Detected.")
			tofEventStream.Send(&zscanproto.ToFEvent{Type: zscanproto.ToFEvent_PENDING})

			if zFusionResponse.StatusMessage == "ROI_FAIL" {
				log.Println("[Warning] Failed to lock ROI.")
				tofEventStream.Send(&zscanproto.ToFEvent{
					Type:      zscanproto.ToFEvent_RESULT,
					LedStatus: zscanproto.LEDStatus_FAILED,
				})

				return
			}

		case zfusionproto.ZFusionResponse_PHASE_FUSION:
			log.Println("[ZFusion] Fusion Phase: Extracting Bitstream...")

			z.HandleFusionMatch(zFusionResponse, tofEventStream)

		}
	}
}

func (z *ZCoreService) HandleFusionMatch(zFusionResponse *zfusionproto.ZFusionResponse, tofEventStream zscanproto.ZPiController_ToFEventStreamClient) {

	if z.ZIPCClient != nil {
		matchCtx, matchCancel := context.WithTimeout(context.Background(), 5*time.Second)
		matchResult, score, err := z.ZIPCClient.MatchTemplate(matchCtx, "zischl", zFusionResponse.FusionBitstream)
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

func (z *ZCoreService) ZCoreEngine() {
	tofEventStream, _ := zpi_client.StartToFStream(z.ZPiClient)

	for {
		evt, err := tofEventStream.Recv()
		if err != nil {
			log.Println("ToF Stream lost:", err)
			return
		}

		if evt.Type == zscanproto.ToFEvent_TRIGGER {
			log.Println("[ZCore] ToF trigger received!")
			go z.HandleFusionSession(tofEventStream)
		}
	}
}
