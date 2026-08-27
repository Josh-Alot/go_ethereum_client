#!/bin/bash
unset ETH_RPC_URL
mkdir -p .anvil/
anvil --block-time 15 --state .anvil/state.json
