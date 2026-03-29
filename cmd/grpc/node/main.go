package main

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/ochestrator/znode_engine"
)

func main() {
	nodeID := "ZTAS@0001"
	nodeAddr := "localhost:50052"
	hubAddr := "104.43.91.57:50051"

	znode_engine.RunZCoreWHub(nodeID, nodeAddr, hubAddr)
}
