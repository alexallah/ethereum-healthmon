package common

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func newHealthy() *BlockTrack {
	bt := &BlockTrack{}
	bt.AddBlock(1)
	bt.AddBlock(2)
	return bt
}

func captureLog(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func Test_BlockTrack(t *testing.T) {
	// can add
	bt := &BlockTrack{}
	if bt.canAdd(0) {
		t.Error("shouldn't be able to add 0")
	}

	if !bt.canAdd(1) {
		t.Error("can not add a block")
	}

	bt.AddBlock(1)

	if bt.canAdd(1) {
		t.Error("can add identical block")
	}

	// distance
	bt = newHealthy()
	distanceOut := captureLog(func() {
		bt.checkDistance()
	})
	if distanceOut != "" {
		t.Error("check distance is not supposed to print anything, got", distanceOut)
	}
	bt.AddBlock(3)
	distanceOut = captureLog(func() {
		bt.checkDistance()
	})
	if strings.HasSuffix(distanceOut, "block track distance 98: 2 > 100") {
		t.Error("unexpected distance message:", distanceOut)
	}

	// add block
	bt = newHealthy()
	bt.AddBlock(3)
	if bt.previous != 2 && bt.current != 3 {
		t.Error("incorrect blocks after AddBlock()")
	}
	bt.AddBlock(3) // add same
	if bt.previous != 2 && bt.current != 3 {
		t.Error("no blocks were supposed to be added")
	}

	// is healthy
	bt = newHealthy()
	if bt.HealthCheck() != nil || !bt.isHealthy() {
		t.Error("should be healthy")
	}
	bt.AddBlock(100) // add a far ahead blcok
	if bt.HealthCheck() == nil || bt.isHealthy() {
		t.Error("should be healthy")
	}
	// old block
	bt = newHealthy()
	bt.currentTime = time.Now().Add(-1 * time.Hour)
	if bt.isHealthy() {
		t.Error("should be unhealthy")
	}
	bt.currentTime = time.Now() // reset
	// no info
	bt = newHealthy()
	bt.previous = 0
	if bt.isHealthy() {
		t.Error("should be unhealthy")
	}

}
