package main

import "github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ipc"

func main() {
	ipc.StartServer(`\\.\pipe\Z-IPC`, ipc.HandleConnection)
}
