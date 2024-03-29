package beacon

import (
	"testing"

	"github.com/alexallah/ethereum-healthmon/internal/common"
)

func newHealhtySyncInfo() *syncInfo {
	// make blocktrack healthy
	blockTrack = common.BlockTrack{}
	blockTrack.AddBlock(1)
	blockTrack.AddBlock(2)
	// generate a healthy info
	return &syncInfo{
		IsSyncing:    false,
		SyncDistance: 2,
	}
}
func Test_SyncInfo(t *testing.T) {
	info := &syncInfo{}

	if checkSyncInfo(info) == nil {
		t.Error("not ready by default")
	}

	// check healthy
	info = newHealhtySyncInfo()
	if checkSyncInfo(info) != nil {
		t.Error("should be healthy")
	}

	// syncing
	info = newHealhtySyncInfo()
	info.IsSyncing = true
	if checkSyncInfo(info) == nil {
		t.Error("should be unhealthy when syncing")
	}

	// distance
	info = newHealhtySyncInfo()
	info.SyncDistance = 10
	if checkSyncInfo(info) == nil {
		t.Error("should be unhealthy when big block distance")
	}

	// optimistic
	info = newHealhtySyncInfo()
	info.IsOptimistic = true
	if checkSyncInfo(info) == nil {
		t.Error("should be unhealthy when optimistic")
	}

	// blocktrack not healthy
	info = newHealhtySyncInfo()
	blockTrack = common.BlockTrack{}
	if checkSyncInfo(info) == nil {
		t.Error("should be unhealthy")
	}

}
