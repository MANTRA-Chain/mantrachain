#!/bin/bash

source "$PWD"/scripts/common.sh

cecho "CYAN" "Setup IBC-Relayer"

pkill -f hermes

mkdir -p ~/.hermes

set -e # exit on first error

tee ~/.hermes/config.toml <<EOF
[global]
log_level = "info"

[[chains]]
account_prefix = "mantra"
clock_drift = "5s"
gas_multiplier = 2.0
grpc_addr = "tcp://127.0.0.1:9081"
id = "mantrachain"
key_name = "relayer"
max_gas = 2000000
rpc_addr = "http://127.0.0.1:26648"
rpc_timeout = "10s"
store_prefix = "ibc"
trusting_period = "14days"
websocket_addr = "ws://127.0.0.1:26648/websocket"

[chains.gas_price]
  denom = "uaum"
  price = 0.0001

[chains.trust_threshold]
  denominator = "3"
  numerator = "1"

[[chains]]
account_prefix = "cosmos"
clock_drift = "5s"
gas_multiplier = 1.1
grpc_addr = "tcp://127.0.0.1:9092"
id = "provider"
key_name = "relayer"
max_gas = 2000000
rpc_addr = "http://127.0.0.1:26657"
rpc_timeout = "10s"
store_prefix = "ibc"
trusting_period = "14days"
websocket_addr = "ws://127.0.0.1:26657/websocket"

[chains.gas_price]
  denom = "stake"
  price = 0.00

[chains.trust_threshold]
  denominator = "3"
  numerator = "1"
EOF

cecho "CYAN" "Import keypair accounts to the IBC-Relayer"
hermes keys delete --chain mantrachain --all
hermes keys delete --chain provider --all
hermes keys add --key-file consumer_key.json --chain mantrachain
hermes keys add --key-file provider_key.json --chain provider

sleep 7

cecho "CYAN" "Create IBC chanel"
hermes create connection \
  --a-chain mantrachain \
  --a-client 07-tendermint-0 \
  --b-client 07-tendermint-0

cecho "CYAN" "Create IBC chanel"
hermes create channel \
    --a-chain mantrachain \
    --a-port consumer \
    --b-port provider \
    --order ordered \
    --channel-version 1 \
    --a-connection connection-0


cecho "CYAN" "Start Hermes"
tmux new -s hermes -d hermes --json start