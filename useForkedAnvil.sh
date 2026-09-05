#!/bin/bash
source ./.secrets.sh
mkdir -p .anvil/
anvil --fork-url "${SEPOLIA_RPC_URL?:define SEPOLIA_RPC_URL on .secrets.sh}"
