package rest

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/alexallah/ethereum-healthmon/internal/beacon"
	"github.com/alexallah/ethereum-healthmon/internal/common"
)

func StartUpdater(state *common.State, addr string, timeout int64, certFile string) {
	client := httpClient(certFile, timeout)
	// add protocol
	if certFile == "" {
		addr = "http://" + addr
	} else {
		addr = "https://" + addr
	}

	go update(state, addr, client)
}

func update(state *common.State, addr string, client *http.Client) {
	for {
		time.Sleep(time.Second)

		err := isReady(addr, client)

		if err != nil {
			state.Error(err)
		} else {
			state.SetHealthy()
		}
	}
}

func httpClient(certFile string, timeout int64) *http.Client {
	var tlsConfig *tls.Config
	if certFile != "" {
		var err error
		tlsConfig, err = beacon.GetTLSConfig(certFile)
		if err != nil {
			log.Panic("can not get dialOption", err)
		}
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}
}
