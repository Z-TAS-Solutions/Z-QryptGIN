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
	conn     *grpc.ClientConn
	cryptSvc zipcproto.CrypticServiceClient
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

		return &ZPIPCClient{
			conn:     conn,
			cryptSvc: zpipcproto.NewCrypticServiceClient(conn),
			pingSvc:  zpipcproto.NewPingServiceClient(conn),
		}, nil
	}
}

func (c *ZPIPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *ZPIPCClient) Ping(ctx context.Context, message string) (string, error) {
	req := &zpipcproto.PingRequest{Message: message}
	resp, err := c.pingSvc.Ping(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ZIPC Ping remote error: %w", err)
	}
	return resp.Reply, nil
}

func (c *ZPIPCClient) MatchTemplate(ctx context.Context, userID string, liveTemplateData []byte) (bool, float32, error) {
	req := &zpipcproto.MatchTemplateRequest{
		UserId:           userID,
		LiveTemplateData: liveTemplateData,
	}
	resp, err := c.cryptSvc.MatchTemplate(ctx, req)
	if err != nil {
		return false, 0, fmt.Errorf("ZIPC MatchTemplate remote error: %w", err)
	}
	if resp.ErrorMessage != "" {
		return resp.IsMatch, resp.ConfidenceScore, fmt.Errorf("cryptic service error: %s", resp.ErrorMessage)
	}
	return resp.IsMatch, resp.ConfidenceScore, nil
}

func (c *ZPIPCClient) StoreEncryptedTemplate(ctx context.Context, userID, templateID, templateType string, rawTemplateData []byte) (bool, error) {
	req := &zpipcproto.StoreTemplateRequest{
		UserId:          userID,
		TemplateId:      templateID,
		TemplateType:    templateType,
		RawTemplateData: rawTemplateData,
	}
	resp, err := c.cryptSvc.StoreEncryptedTemplate(ctx, req)
	if err != nil {
		return false, fmt.Errorf("ZIPC StoreEncryptedTemplate remote error: %w", err)
	}
	if !resp.Success || resp.ErrorMessage != "" {
		return false, fmt.Errorf("failed to store template: %s", resp.ErrorMessage)
	}
	return resp.Success, nil
}
