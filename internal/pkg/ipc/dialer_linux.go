//go:build linux

package ipc

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

func getPlatformDialConfig() (string, func(context.Context, string) (net.Conn, error), []grpc.DialOption) {
	target := "/tmp/zpipcproto.sock"
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", addr)
	}

	return target, dialer, nil
}
