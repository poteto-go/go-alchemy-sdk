#!/bin/bash
set -euo pipefail

SCRIPT_DIR=$(dirname "$0")
PARAMS_FILE_PATH="$SCRIPT_DIR/../network_params.yaml"

kurtosis run --enclave gas-testnet github.com/ethpandaops/ethereum-package --args-file "$PARAMS_FILE_PATH"

sleep 10  # Wait for network to be ready
