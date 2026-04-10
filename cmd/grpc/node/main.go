package main

import (
	"log"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/ochestrator/znode_engine"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/utils"
)

func main() {
	NodeIP, err := utils.GetPublicIP()
	if err != nil {
		log.Fatalf("Failed to retrieve Node IP!")
	}

	nodeID := "ZTAS@0001"
	nodeAddr := "localhost:50052"
	hubAddr := "104.43.91.57:50051"

	znode_engine.RunZCoreWHub(nodeID, nodeAddr, NodeIP, hubAddr)
}
