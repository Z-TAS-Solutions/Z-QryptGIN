package zcore

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/znode_monitor"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service/zpi_client"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zcoreproto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zfusionproto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zscanproto"
)

type EventType int

const (
	ZPiTrigger EventType = iota
	ZPiTriggerOffline
	ZHubTrigger
)

type ZEvent struct {
	Type    EventType
	Payload interface{}
}

type ZCoreService struct {
	ZIPCClient *ipc.ZIPCClient
	ZPiClient  zscanproto.ZPiControllerClient
	ZFusion    zfusionproto.FusionCaptureServiceClient
	ZCoreHub   zcoreproto.ZCoreServiceClient

	ServiceMonitor *znodemonitor.ZNodeMonitor

	EventChannel chan ZEvent
	sessionCount int

	Mu                sync.Mutex
	EnrollmentPending bool
	PendingUserID     string
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

func (z *ZCoreService) HandleEnrollSession(userID string, tofEventStream zscanproto.ZPiController_ToFEventStreamClient) {
	log.Printf("[ZCore] Initializing ZFusion Enrollment for User: %s", userID)

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	zFusionStream, err := z.ZFusion.FusionCapture(ctx, &zfusionproto.ZFusionRequest{
		SessionId: "ENROLL_" + strconv.Itoa(z.sessionCount) + "_" + strconv.FormatInt(time.Now().Unix(), 10),
	})

	if err != nil {
		log.Println("[Error] Could not start Fusion Capture for Enrollment:", err)
		return
	}

	for {
		zFusionResponse, err := zFusionStream.Recv()
		if err == io.EOF {
			log.Println("[ZFusion] Enrollment Session Closed.")
			break
		}
		if err != nil {
			log.Printf("[Error] Fusion Enrollment Stream interrupted: %v", err)
			break
		}

		switch zFusionResponse.CompletionPhase {
		case zfusionproto.ZFusionResponse_PHASE_ROI:
			log.Println("[ZFusion] ROI Phase: Hand Detected for Enrollment.")
			tofEventStream.Send(&zscanproto.ToFEvent{Type: zscanproto.ToFEvent_PENDING})

			if zFusionResponse.StatusMessage == "ROI_FAIL" {
				log.Println("[Warning] Failed to lock ROI for Enrollment.")
				tofEventStream.Send(&zscanproto.ToFEvent{
					Type:      zscanproto.ToFEvent_RESULT,
					LedStatus: zscanproto.LEDStatus_FAILED,
				})
				return
			}

		case zfusionproto.ZFusionResponse_PHASE_FUSION:
			log.Printf("[ZFusion] Fusion Phase: Extracting Bitstream for User %s...", userID)
			z.HandleEnrollStore(userID, zFusionResponse, tofEventStream)
		}
	}
}

func (z *ZCoreService) HandleEnrollStore(userID string, zFusionResponse *zfusionproto.ZFusionResponse, tofEventStream zscanproto.ZPiController_ToFEventStreamClient) {
	if z.ZIPCClient != nil {
		storeCtx, storeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer storeCancel()

		templateID := fmt.Sprintf("tmpl_%d", time.Now().Unix())
		success, err := z.ZIPCClient.StoreEncryptedTemplate(storeCtx, userID, templateID, "fusion", zFusionResponse.FusionBitstream)

		var ledStatus zscanproto.LEDStatus
		if err != nil || !success {
			log.Printf("[ZCore] Enrollment FAILED for User %s: %v", userID, err)
			ledStatus = zscanproto.LEDStatus_FAILED
		} else {
			log.Printf("[ZCore] Enrollment SUCCESS for User %s (TemplateID: %s)", userID, templateID)
			ledStatus = zscanproto.LEDStatus_SUCCESS
		}

		tofEventStream.Send(&zscanproto.ToFEvent{
			Type:      zscanproto.ToFEvent_RESULT,
			LedStatus: ledStatus,
		})
	}
}

func (z *ZCoreService) HandleFusionMatch(zFusionResponse *zfusionproto.ZFusionResponse, tofEventStream zscanproto.ZPiController_ToFEventStreamClient) {

	if z.ZIPCClient != nil {
		matchCtx, matchCancel := context.WithTimeout(context.Background(), 5*time.Second)
		matchResult, score, err := z.ZIPCClient.MatchTemplate(matchCtx, zFusionResponse.FusionBitstream)
		matchCancel()

		var ledStatus zscanproto.LEDStatus
		if err != nil || !matchResult {
			log.Printf("Match Denied (Score: %f, Error: %v)", score, err)
			ledStatus = zscanproto.LEDStatus_FAILED
		} else {
			log.Printf("Match SUCCESS (Score: %f)", score)
			ledStatus = zscanproto.LEDStatus_SUCCESS

			go func(uid string) {
				twoFACtx, twoFACancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer twoFACancel()

				log.Printf("[ZCore] Sending 2FA Request for User: %s", uid)
				resp, err := z.ZCoreHub.Request2FA(twoFACtx, &zcoreproto.TwoFARequest{UserId: uid})
				if err != nil {
					log.Printf("[ZCore] 2FA Request failed for User %s: %v", uid, err)
					return
				}
				log.Printf("[ZCore] 2FA Response for User %s: %s (Success: %v)", uid, resp.Message, resp.Success)
			}("zischl")
		}

		tofEventStream.Send(&zscanproto.ToFEvent{
			Type:      zscanproto.ToFEvent_RESULT,
			LedStatus: ledStatus,
		})
	}
}

func (z *ZCoreService) ZCoreEngine(eventQueue chan ZEvent) {
	for {
		log.Println("[ZCore] Initializing ToF Stream...")

		tofEventStream, err := zpi_client.StartToFStream(z.ZPiClient)
		if err != nil {
			log.Printf("[ZCore] Connection failed: %v. Retrying in 3s...", err)
			time.Sleep(3 * time.Second)
			continue
		}

		log.Println("[ZCore] ToF Stream connected and active.")

		for {
			tofEvent, err := tofEventStream.Recv()
			if err != nil {
				log.Printf("[ZCore] Stream lost: %v. Attempting reconnect...", err)
				eventQueue <- ZEvent{Type: EventType(1), Payload: err}
				break
			}

			if tofEvent.Type == zscanproto.ToFEvent_TRIGGER {
				eventQueue <- ZEvent{
					Type:    EventType(0),
					Payload: tofEventStream,
				}
			}

		}
	}
}
