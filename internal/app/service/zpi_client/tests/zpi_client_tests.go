package zpi_client_tests

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zscanproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunZPiClient(ip string) zscanproto.ZPiControllerClient {
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	//defer cancel()

	client := zscanproto.NewZPiControllerClient(conn)
	log.Print("Connected to ZPi GRPC Host!")

	statusResp, _ := client.SetLEDStatus(context.Background(), &zscanproto.LEDStatusRequest{
		Status: zscanproto.LEDStatus_SUCCESS,
	})
	fmt.Println("SetLEDStatus:", statusResp.Message)

	ledState, _ := client.GetLED(context.Background(), &zscanproto.Empty{})
	fmt.Printf("LED State - R:%d G:%d B:%d\n", ledState.Red, ledState.Green, ledState.Blue)

	tofConfigResp, _ := client.ConfigureToF(context.Background(), &zscanproto.ToFConfig{
		Threshold: 150,
	})
	fmt.Println("ConfigureToF:", tofConfigResp.Message)

	enableResp, _ := client.EnableToF(context.Background(), &zscanproto.Empty{})
	fmt.Println("EnableToF:", enableResp.Message)

	activateResp, _ := client.ActivateToF(context.Background(), &zscanproto.Empty{})
	fmt.Println("ActivateToF:", activateResp.Message)

	time.Sleep(500 * time.Millisecond)

	deactivateResp, _ := client.DeactivateToF(context.Background(), &zscanproto.Empty{})
	fmt.Println("DeactivateToF:", deactivateResp.Message)

	disableResp, _ := client.DisableToF(context.Background(), &zscanproto.Empty{})
	fmt.Println("DisableToF:", disableResp.Message)

	response, _ := client.EnableToF(context.Background(), &zscanproto.Empty{})
	fmt.Println(response)

	response, _ = client.ActivateToF(context.Background(), &zscanproto.Empty{})
	fmt.Println(response)

	reconfigResp, _ := client.ConfigureToF(context.Background(), &zscanproto.ToFConfig{
		Threshold: 200,
	})
	fmt.Println("ReconfigureToF while active:", reconfigResp.Message)

	return client
}
