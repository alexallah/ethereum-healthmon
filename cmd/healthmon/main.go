package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	beaconGRPC "github.com/alexallah/ethereum-healthmon/internal/beacon/grpc"
	beaconREST "github.com/alexallah/ethereum-healthmon/internal/beacon/rest"
	"github.com/alexallah/ethereum-healthmon/internal/common"
	"github.com/alexallah/ethereum-healthmon/internal/execution"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Chain   string `long:"chain" description:"Ethereum chain" choice:"execution" choice:"beacon" required:"true"`
	Port    int    `long:"port" description:"Node port (default: 8545 for execution, 4000 for beacon)"`
	Addr    string `long:"addr" description:"Node address" default:"localhost"`
	Timeout int64  `long:"timeout" description:"Node connection timeout, seconds" default:"5"`

	Execution struct {
		Jwt string `long:"engine-jwt" description:"JWT hex secret path. Use only when connecting to the engine RPC endpoint."`
	} `group:"Execution chain" namespace:"execution"`

	Beacon struct {
		Certificate string `long:"certificate" description:"TLS root certificate path. Specify only if you have it configured for your node as well."`
		Prysm       struct {
			GRPC bool `long:"grpc" description:"Required if you're connecting to a Prysm node. If Prysm migrates to a JSON-RPC protocol in the future versions, this flag should be removed."`
		} `group:"Prysm" namespace:"prysm"`
	} `group:"Beacon chain" namespace:"beacon"`

	Service struct {
		Port int    `long:"port" description:"healthmon listening port" default:"21171"`
		Addr string `long:"addr" description:"healthmon listening address. Set to 0.0.0.0 to allow external access." default:"localhost"`
	} `group:"healthcheck service" namespace:"http"`
}

var state *common.State

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		os.Exit(0)
	}

	state = new(common.State)

	// default node port
	nodePort := 8545 // execution
	if opts.Chain == "beacon" {
		nodePort = 4000
	}
	// custom override
	if opts.Port != 0 {
		nodePort = opts.Port
	}
	nodeAddr := fmt.Sprintf("%s:%d", opts.Addr, nodePort)

	switch opts.Chain {
	case "beacon":
		if opts.Beacon.Prysm.GRPC {
			beaconGRPC.StartUpdater(state, nodeAddr, opts.Timeout, opts.Beacon.Certificate)
		} else {
			beaconREST.StartUpdater(state, nodeAddr, opts.Timeout, opts.Beacon.Certificate)
		}
	case "execution":
		execution.StartUpdater(state, nodeAddr, opts.Timeout, opts.Execution.Jwt)
	default:
		log.Fatalf("unknown chain: %s.\n", opts.Chain)
	}

	log.Printf("%s node address is %s", opts.Chain, nodeAddr)
	serviceAddr := fmt.Sprintf("%s:%d", opts.Service.Addr, opts.Service.Port)
	log.Printf("healthmon listenting on %s\n", serviceAddr)

	http.HandleFunc("/ready", statusHandler)
	http.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe(serviceAddr, nil)
}

func statusHandler(w http.ResponseWriter, req *http.Request) {
	if state.IsHealthy() {
		io.WriteString(w, "OK")
	} else {
		w.WriteHeader(503)
		io.WriteString(w, "NOT READY")
	}
}

func metricsHandler(w http.ResponseWriter, req *http.Request) {
	var ready int
	if state.IsHealthy() {
		ready = 1
	}
	io.WriteString(w, fmt.Sprintf("# TYPE ready gauge\nready %d\n", ready))
}
