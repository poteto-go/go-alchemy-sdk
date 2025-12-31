---
sidebar_position: 2
---

# Kurtosis Development

## Kurtosis

> The open source dev environment engine for blockchain infra — made by and for protocol, platform, and dev ops teams.
> https://www.kurtosis.com/

## Kurtosis Development

This tutorial connects to the kurtosis network launched on your localhost.

:::warning

- Please note that kurtosis is often a technique to aid development and is not a tool recommended for production environments.

:::

You can get full-code here:

- [link](https://github.com/poteto-go/go-alchemy-sdk/tree/main/examples/kurtosis-development)

### 1. Setup

```bash
make kurtosis-development
cd kurtosis-development
go mod init example.com/example/kurtosis-development

mkdir network
mkdir scripts
```

for example:

```
kurtosis-development
├── network
├── scripts
└── go.mod
```

### 2. Configure Network

EL: geth

CL: lighthouse

```yaml title="network/network_params.yaml"
participants:
  - el_type: geth
    cl_type: lighthouse
network_params:
  network_id: "585858"
```

### 3. Setup Scripts for Kurtosis

```sh title="scripts/run.sh"
#!/bin/bash
set -euo pipefail

SCRIPT_DIR=$(dirname "$0")
PARAMS_FILE_PATH="$SCRIPT_DIR/../network/network_params.yaml"

kurtosis run --enclave gas-testnet github.com/ethpandaops/ethereum-package --args-file "$PARAMS_FILE_PATH"

sleep 10  # Wait for network to be ready
```

```sh title="scripts/detect-rpc-port.sh"
#!/bin/bash
set -euo pipefail

echo $(kurtosis enclave inspect gas-testnet | grep "rpc: 8545/tcp" | grep -oh "127.0.0.1\:[0-9]*" | cut -d':' -f2)
```

### 4. Start Kurtosis

```bash
bash ./scripts/run.sh
bash ./scripts/detect-rpc-port.sh
# => 32769
```

### 5. Development with Kurtosis

```go title="main.go"
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/poteto-go/go-alchemy-sdk/gas"
)

var alchemy gas.Alchemy

func init() {
	port, err := strconv.Atoi(os.Getenv("RPC_PORT"))
	if err != nil {
		panic(err)
	}

	setting := gas.AlchemySetting{
		PrivateNetworkConfig: gas.PrivateNetworkConfig{
			Host: "127.0.0.1",
			Port: port,
		},
	}
	alchemy = gas.NewAlchemy(setting)
}

func main() {
	blockNumber, err := alchemy.Core.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println(blockNumber)
}
```

and run

```bash
RPC_PORT=32769 go run main.go
# => blockNumber > 0
```
