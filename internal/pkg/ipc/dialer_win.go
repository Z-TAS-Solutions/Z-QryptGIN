//go:build windows

package ipc

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
	"google.golang.org/grpc"
)

func getPlatformDialConfig() (string, func(context.Context, string) (net.Conn, error), []grpc.DialOption) {
	target := `\\.\pipe\zpipcproto`
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return winio.DialPipeContext(ctx, addr)
	}

	options := []grpc.DialOption{
		grpc.WithAuthority("localhost"),
	}

	return target, dialer, options
}
