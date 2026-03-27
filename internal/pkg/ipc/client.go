package ipc

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime"

	"github.com/Microsoft/go-winio"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zpipcproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type ZPIPCClient struct {
	conn     *grpc.ClientConn
	cryptSvc zpipcproto.CrypticServiceClient
	pingSvc  zpipcproto.PingServiceClient
}

func DialIPC() (*ZPIPCClient, error) {
	var conn *grpc.ClientConn
	var err error

	if runtime.GOOS == "windows" {
		pipeAddr := `\\.\pipe\zpipcproto`
		dialer := func(ctx context.Context, addr string) (net.Conn, error) {
			return winio.DialPipeContext(ctx, addr)
		}
		conn, err = grpc.Dial(pipeAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(dialer),
			grpc.WithAuthority("localhost"),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to dial Windows IPC %s: %w", pipeAddr, err)
		}
	} else {
		socketPath := "/tmp/zpipcproto.sock"
		dialer := func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", addr)
		}
		conn, err = grpc.Dial(socketPath,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(dialer),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to dial Unix IPC %s: %w", socketPath, err)
		}
	}

	client := &ZPIPCClient{
		conn:     conn,
		cryptSvc: zpipcproto.NewCrypticServiceClient(conn),
		pingSvc:  zpipcproto.NewPingServiceClient(conn),
	}

	go func() {
		conn.Connect()
		state := conn.GetState()
		wasReady := false
		for {
			if !conn.WaitForStateChange(context.Background(), state) {
				break
			}
			state = conn.GetState()

			if state == connectivity.Ready {
				if !wasReady {
					log.Println("Rust IPC Connection Established!")
					wasReady = true
				}
			} else if state == connectivity.TransientFailure || state == connectivity.Idle || state == connectivity.Connecting {
				if wasReady {
					log.Println("[WARNING] Rust IPC Connection Lost. Waiting for reconnect...")
					wasReady = false
				}
				
				if state == connectivity.Idle {
					conn.Connect()
				}
			}
		}
	}()

	return client, nil
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
		return "", fmt.Errorf("IPC Ping remote error: %w", err)
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
		return false, 0, fmt.Errorf("IPC MatchTemplate remote error: %w", err)
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
		return false, fmt.Errorf("IPC StoreEncryptedTemplate remote error: %w", err)
	}
	if !resp.Success || resp.ErrorMessage != "" {
		return false, fmt.Errorf("failed to store template: %s", resp.ErrorMessage)
	}
	return resp.Success, nil
}
