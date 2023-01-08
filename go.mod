module github.com/alexallah/ethereum-healthmon

go 1.19

require (
	github.com/golang-jwt/jwt/v4 v4.4.3
	github.com/golang/protobuf v1.5.2
	github.com/jessevdk/go-flags v1.5.0
	github.com/prysmaticlabs/prysm v1.4.2-0.20220805185555-4e225fc667d8
	google.golang.org/grpc v1.51.0
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/ethereum/go-ethereum v1.10.26 // indirect
	github.com/ferranbt/fastssz v0.0.0-20210526181520-7df50c8568f8 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.3 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prysmaticlabs/eth2-types v0.0.0-20210303084904-c9735a06829d // indirect
	github.com/prysmaticlabs/ethereumapis v0.0.0-20220804141212-99a982d2918a // indirect
	github.com/prysmaticlabs/fastssz v0.0.0-20221107182844-78142813af44 // indirect
	github.com/prysmaticlabs/go-bitfield v0.0.0-20210809151128-385d8c5e3fb7 // indirect
	github.com/prysmaticlabs/gohashtree v0.0.2-alpha // indirect
	github.com/thomaso-mirodin/intmath v0.0.0-20160323211736-5dc6d854e46e // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230106154932-a12b697841d9 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/grpc-ecosystem/grpc-gateway/v2 => github.com/prysmaticlabs/grpc-gateway/v2 v2.3.1-0.20210622145107-ca3041e1b380
