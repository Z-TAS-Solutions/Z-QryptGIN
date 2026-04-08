package ipc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zipcproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ZIPCClient struct {
	conn       *grpc.ClientConn
	crypticSvc zipcproto.CrypticServiceClient
}

func RunZIPCClient() (*ZIPCClient, error) {
	target, dialer, options := getPlatformDialConfig()

	options = append(options,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
		grpc.WithBlock(),
	)

	log.Printf("[ZIPC] Attempting to connect to Rust ZIPC via %s...", target)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		conn, err := grpc.DialContext(ctx, target, options...)
		if err != nil {
			log.Printf("[ZIPC] Connection failed to %s: %v. Retrying in 7 seconds...", target, err)
			cancel()
			time.Sleep(7 * time.Second) // Fixed to match your log message
			continue
		}

		cancel()
		log.Printf("[ZIPC] Successfully connected to Rust ZIPC via %s", target)

		return &ZIPCClient{
			conn:       conn,
			crypticSvc: zipcproto.NewCrypticServiceClient(conn),
		}, nil
	}
}

func (c *ZIPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *ZIPCClient) Ping(ctx context.Context, message string) (string, error) {
	req := &zipcproto.PingRequest{Message: message}
	resp, err := c.crypticSvc.Ping(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ZIPC Ping remote error: %w", err)
	}
	return resp.Ping, nil
}

func (c *ZIPCClient) MatchTemplate(ctx context.Context, livePayload []byte) (bool, float32, error) {
	req := &zipcproto.MatchTemplateRequest{
		Payload: livePayload,
	}
	resp, err := c.crypticSvc.MatchTemplate(ctx, req)
	if err != nil {
		return false, 0, fmt.Errorf("ZIPC MatchTemplate remote error: %w", err)
	}
	if resp.ResponseMessage != "" {
		return resp.Matched, resp.MatchScore, fmt.Errorf("cryptic service error: %s", resp.ResponseMessage)
	}
	return resp.Matched, resp.MatchScore, nil
}

func (c *ZIPCClient) StoreEncryptedTemplate(ctx context.Context, userID, templateID, templateType string, rawTemplateData []byte) (bool, error) {
	req := &zipcproto.StoreTemplateRequest{
		UserId:       userID,
		TemplateType: templateType,
		Payload:      rawTemplateData,
	}
	resp, err := c.crypticSvc.StoreEncryptedTemplate(ctx, req)
	if err != nil {
		return false, fmt.Errorf("ZIPC StoreEncryptedTemplate remote error: %w", err)
	}
	if !resp.State || resp.ResponseMessage != "" {
		return false, fmt.Errorf("failed to store template: %s", resp.ResponseMessage)
	}
	return resp.State, nil
}
