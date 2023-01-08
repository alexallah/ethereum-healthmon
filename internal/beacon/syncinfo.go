package beacon

import (
	"fmt"
	"log"

	"github.com/alexallah/ethereum-healthmon/internal/common"
)

var blockTrack common.BlockTrack

type SyncInfo struct {
	IsSyncing    bool   `json:"is_syncing"`
	IsOptimistic bool   `json:"is_optimistic"`
	HeadSlot     uint64 `json:"head_slot,string"`
	SyncDistance uint64 `json:"sync_distance,string"`
}

func CheckSyncInfo(syncInfo *SyncInfo) error {
	// distance
	if syncInfo.IsSyncing || syncInfo.SyncDistance >= 5 {
		return fmt.Errorf("syncing, distance %d slots", syncInfo.SyncDistance)
	} else if syncInfo.SyncDistance > 1 {
		log.Printf("syncing distance %d\n", syncInfo.SyncDistance)
	}
	// optimistic doesn't work for validator
	if syncInfo.IsOptimistic {
		return fmt.Errorf("beacon node is optimistic")
	}
	// waiting for a new block
	blockTrack.AddBlock(uint64(syncInfo.HeadSlot))
	if err := blockTrack.HealthCheck(); err != nil {
		return err
	}

	return nil
}
