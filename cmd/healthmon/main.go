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
	Port    int    `long:"port" description:"Node port (default: auto, default for the chain)"`
	Addr    string `long:"addr" description:"Node address" default:"localhost"`
	Timeout int64  `long:"timeout" description:"Node connection timeout, seconds" default:"5"`

	Execution struct {
		Jwt string `long:"jwt" description:"JWT hex secret path"`
	} `group:"Execution chain" namespace:"execution"`

	Beacon struct {
		Certificate string `long:"certificate" description:"TLS root certificate path. Specify only if have it configured for your node as well."`
		Prysm       struct {
			GRPC bool `long:"grpc" description:"Required if you're connecting to a Prysm node. If Prysm migrates to a JSON-RPC protocol in the future versions, this flag should be removed."`
		} `group:"Prysm" namespace:"prysm"`
	} `group:"Beacon chain" namespace:"beacon"`

	Service struct {
		Port int    `long:"port" description:"healthmon service port" default:"21171"`
		Addr string `long:"addr" description:"healthmon service address" default:"0.0.0.0"`
	} `group:"healthcheck service" namespace:"service"`
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
	nodePort := 30303 // execution
	if opts.Chain == "beacon" {
		nodePort = 4000
	}
	// custom override
	if opts.Port != 0 {
		nodePort = opts.Port
	}
	fullNodeAddr := fmt.Sprintf("%s:%d", opts.Addr, nodePort)

	switch opts.Chain {
	case "beacon":
		if opts.Beacon.Prysm.GRPC {
			beaconGRPC.StartUpdater(state, fullNodeAddr, opts.Timeout, opts.Beacon.Certificate)
		} else {
			beaconREST.StartUpdater(state, fullNodeAddr, opts.Timeout, opts.Beacon.Certificate)
		}
	case "execution":
		if opts.Execution.Jwt == "" {
			log.Fatalln("JWT path is not specified")
		}
		execution.StartUpdater(state, fullNodeAddr, opts.Timeout, opts.Execution.Jwt)
	default:
		log.Fatalf("unknown chain: %s.\n", opts.Chain)
	}

	log.Printf("%s node address is %s", opts.Chain, fullNodeAddr)
	fullServiceAddr := fmt.Sprintf("%s:%d", opts.Service.Addr, opts.Service.Port)
	log.Printf("healthmon listenting on %s\n", fullServiceAddr)

	http.HandleFunc("/ready", statusHandler)
	http.ListenAndServe(fullServiceAddr, nil)
}

func statusHandler(w http.ResponseWriter, req *http.Request) {
	if state.IsHealthy() {
		io.WriteString(w, "OK")
	} else {
		w.WriteHeader(503)
		io.WriteString(w, "Down")
	}
}
