#!/bin/bash
set -euo pipefail

echo "deb [trusted=yes] https://sdk.kurtosis.com/kurtosis-cli-release-artifacts/ /" | sudo tee /etc/apt/sources.list.d/kurtosis.list && sudo apt update
sudo apt update
sudo apt install kurtosis-cli
