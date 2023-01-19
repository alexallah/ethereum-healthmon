# ethereum-healthmon

[![Build](https://github.com/alexallah/ethereum-healthmon/actions/workflows/build.yaml/badge.svg)](https://github.com/alexallah/ethereum-healthmon/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/alexallah/ethereum-healthmon)](https://goreportcard.com/report/alexallah/ethereum-healthmon)

Health monitoring for beacon and execution (eth1) nodes.

## Main features

- Doesn't rely on any third-party services. Connects only to the specified node.
- Works with both beacon and execution chains.
- Can be easily integrated with a load balancer.
- Exposes metrics for Prometheus.
- Supports Prysm's gRPC protocol.
- Can connect to the Engine RPC execution endpoint with JWT-based security.
- Optional TLS encryption for beacon nodes.

## HTTP API

`/ready` Returns HTTP 200 when the node is healthy. It can be used by load balancers or any other healthcheck tools.

Example healthy

```
> curl -i localhost:21171/ready
HTTP/1.1 200 OK
Date: Sun, 08 Jan 2023 22:23:51 GMT
Content-Length: 2
Content-Type: text/plain; charset=utf-8

OK
```

Example not healthy

```
> curl -i localhost:21171/ready
HTTP/1.1 503 Service Unavailable
Date: Sun, 08 Jan 2023 22:26:57 GMT
Content-Length: 9
Content-Type: text/plain; charset=utf-8

NOT READY
```

`/metrics` - export the readiness status using Prometheus metrics format. It can be collected by Prometheus or other compatible clients. Example output:

```
# TYPE ready gauge
ready 1
```

## Build

### Compile from source

Requirements: Git, Go

```bash
# get the source code
git clone https://github.com/alexallah/ethereum-healthmon
cd ethereum-healthmon
# compile the binary
go build -o healthmon cmd/healthmon/main.go
# ready
./healthmon --help
```

### Build using Docker, only for Linux

Requirements: Docker

```bash
# build a docker image with the binary
docker build https://github.com/alexallah/ethereum-healthmon.git#main --tag healthmon
# copy the binary to the local file system
docker run --rm --entrypoint "" healthmon cat /usr/bin/ethereum-healthmon > healthmon
chmod +x healthmon
# ready
./healthmon --help
```

## Usage examples

#### Monitor to Prysm using a secure connection

```
./healthmon --chain=beacon --beacon.prysm.grpc --beacon.certificate=tls/ca.cert
```

The --beacon.prysm.grpc flag is required for Prysm at the moment.
The certificate needs to be configured on the Prysm node as well. Using it is optional but recommended if programs are running on different computers.

#### Connect to Geth with JWT

```
./healthmon --chain execution --port 8551 --execution.engine-jwt jwt.hex
```

We are connecting to an Engine RPC endpoint instead of a regular one on port 8545.
That's why we need to provide a JWT secret file.

### Docker Compose

Run alongside Besu in a Docker container.

```yaml
version: "3.7"
services:
  execution-node:
    image: hyperledger/besu
    command:
      - --rpc-http-enabled
      - --rpc-http-host=0.0.0.0
      - --host-allowlist=*
  healthmon:
    build: /home/alex/Projects/healthmon
    command:
      - --chain=execution
      - --addr=execution-node
      - --http.addr=0.0.0.0
```

## Notes

By default, the service is listening on localhost. If you are using Docker or need any external access to its APIs (very likely), set http.host value. `--http.host=0.0.0.0` would allow access from any network, for example.

### Execution

- It is possible to connect to a regular RPC or Engine RPC. Some clients have them on the same port.
- `--execution.engine-jwt` is required when connecting to an Engine RPC endpoint. Make sure it is the same one you have configured for your execution node.

### Beacon

Prysm uses a custom gRPC API to connect between beacon and validator nodes. They might move to the standard JSON API in the future, but for now, it is necessary to use `--beacon.prysm.grpc` when connecting to a Prysm node.
