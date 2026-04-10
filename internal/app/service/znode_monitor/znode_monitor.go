package znodemonitor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type ZNodeMonitor struct {
	MonitorContext context.Context
	EndMonitor     context.CancelFunc
}

// this right here's gonna be the new state watcher for mainly logging along side grpc keep alives
func RunNodeMonitor(ctx context.Context, serviceName string, conn *grpc.ClientConn) {
	go func() {
		state := conn.GetState()
		for {
			if conn.WaitForStateChange(ctx, state) {
				state = conn.GetState()

				switch state {
				case connectivity.Ready:
					log.Printf("[%s] Service Ready", serviceName)
				case connectivity.TransientFailure:
					log.Printf("[%s] Connection Failure, commencing keep alive protocols!", serviceName)
				case connectivity.Connecting:
					log.Printf("[%s] Connecting...", serviceName)
				case connectivity.Idle:
					log.Printf("[%s] Idle, gimme some work!", serviceName)
				case connectivity.Shutdown:
					log.Printf("[%s] Shuting down, bye !", serviceName)
					return
				}

			} else {
				// Either the context was cancelled or the connection was closed
				log.Printf("[%s] Watcher shutting down...", serviceName)
				return
			}
		}
	}()
}
