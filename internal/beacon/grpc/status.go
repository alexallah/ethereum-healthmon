package grpc

import (
	"context"
	"time"

	"github.com/alexallah/ethereum-healthmon/internal/beacon"
	"github.com/golang/protobuf/ptypes/empty"
	eth "github.com/prysmaticlabs/prysm/proto/eth/service"
	v1 "github.com/prysmaticlabs/prysm/proto/eth/v1"
	grpc "google.golang.org/grpc"
)

func toSyncInfo(gSyncInfo *v1.SyncInfo) *beacon.SyncInfo {
	return &beacon.SyncInfo{
		IsSyncing:    gSyncInfo.IsSyncing,
		IsOptimistic: gSyncInfo.IsOptimistic,
		HeadSlot:     uint64(gSyncInfo.HeadSlot),
		SyncDistance: uint64(gSyncInfo.SyncDistance),
	}
}

func getSyncing(conn *grpc.ClientConn, timeout int64) (*beacon.SyncInfo, error) {

	c := eth.NewBeaconNodeClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	in := new(empty.Empty)
	response, err := c.GetSyncStatus(ctx, in)
	if err != nil {
		return nil, err
	}
	return toSyncInfo(response.Data), nil
}

func isHealthy(conn *grpc.ClientConn, timeout int64) error {
	c := eth.NewBeaconNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	in := new(empty.Empty)
	_, err := c.GetHealth(ctx, in)
	if err != nil {
		return err
	}
	return nil
}

func isReady(addr string, timeout int64, opts ...grpc.DialOption) error {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	// make sure we are not syncing
	syncInfo, err := getSyncing(conn, timeout)
	if err != nil {
		return err
	}

	if err := beacon.CheckSyncInfo(syncInfo); err != nil {
		return err
	}

	// health call
	err = isHealthy(conn, timeout)
	if err != nil {
		return err
	}

	return nil
}
