package grpc

import (
	"log"
	"time"

	"github.com/alexallah/ethereum-healthmon/internal/beacon"
	"github.com/alexallah/ethereum-healthmon/internal/common"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func StartUpdater(state *common.State, addr string, timeout int64, certFile string) {
	dialOption := getDialOption(certFile)
	go update(state, addr, timeout, dialOption)
}

func update(state *common.State, addr string, timeout int64, dial grpc.DialOption) {
	for {
		time.Sleep(time.Second)

		err := isReady(addr, timeout, dial)

		if err != nil {
			state.Error(err)
		} else {
			state.SetHealthy()
		}
	}
}

func getDialOption(certFile string) grpc.DialOption {
	var creds credentials.TransportCredentials
	if certFile != "" {
		tlsConfig, err := beacon.GetTLSConfig(certFile)
		if err != nil {
			log.Panic("can not load certificate", err)
		}
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	return grpc.WithTransportCredentials(creds)
}
