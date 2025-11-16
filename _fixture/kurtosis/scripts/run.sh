#!/bin/bash
set -euo pipefail

kurtosis run --enclave gas-testnet github.com/ethpandaops/ethereum-package --args-file network_params.yaml
