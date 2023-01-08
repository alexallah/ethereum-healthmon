package common

import (
	"fmt"
	"log"
	"time"
)

type BlockTrack struct {
	previous uint64

	current     uint64
	currentTime time.Time
}

func (b *BlockTrack) canAdd(blockNum uint64) bool {
	if blockNum == 0 {
		return false
	}

	if b.current == blockNum {
		// seen already
		return false
	}

	return true
}

func (b *BlockTrack) checkDistance() {
	if b.previous == 0 {
		return
	}
	distance := b.current - b.previous
	if distance > 1 {
		log.Printf("block track distance %d: %d > %d\n", distance, b.previous, b.current)
	}
}

func (b *BlockTrack) AddBlock(blockNum uint64) {
	if !b.canAdd(blockNum) {
		return
	}

	// store previous
	b.previous = b.current

	// add block
	b.current = blockNum
	b.currentTime = time.Now()

	b.checkDistance()
}

func (b *BlockTrack) HealthCheck() error {
	if !b.isHealthy() {
		return fmt.Errorf("waiting for new blocks, latest %d", b.current)
	}
	return nil
}

func (b *BlockTrack) isHealthy() bool {
	// no block info yet
	if b.current == 0 || b.previous == 0 {
		return false
	}
	// difference is too high, let's wait for a new block
	if b.current-b.previous > 10 {
		return false
	}
	// the last block was added long time ago
	if time.Since(b.currentTime).Minutes() > 5 {
		return false
	}

	return true
}
