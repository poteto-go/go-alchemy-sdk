#!/bin/bash
set -euo pipefail

echo $(kurtosis enclave inspect gas-testnet | grep "rpc: 8545/tcp" | grep -oh "127.0.0.1\:[0-9]*" | cut -d':' -f2)
