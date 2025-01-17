#!/bin/bash

set -ex

# initialize Hermes relayer configuration
mkdir -p /root/.hermes/
touch /root/.hermes/config.toml

echo $MANTRA_B_E2E_RLY_MNEMONIC >/root/.hermes/MANTRA_B_E2E_RLY_MNEMONIC.txt
echo $MANTRA_A_E2E_RLY_MNEMONIC >/root/.hermes/MANTRA_A_E2E_RLY_MNEMONIC.txt

# setup Hermes relayer configuration with non-zero gas_price
tee /root/.hermes/config.toml <<EOF
[global]
log_level = 'info'

[mode]

[mode.clients]
enabled = true
refresh = true
misbehaviour = true

[mode.connections]
enabled = false

[mode.channels]
enabled = true

[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true
tx_confirmation = true

[rest]
enabled = true
host = '0.0.0.0'
port = 3031

[telemetry]
enabled = true
host = '127.0.0.1'
port = 3001

[[chains]]
id = '$MANTRA_A_E2E_CHAIN_ID'
rpc_addr = 'http://$MANTRA_A_E2E_VAL_HOST:26657'
grpc_addr = 'http://$MANTRA_A_E2E_VAL_HOST:9090'
event_source = { mode = 'pull', interval = '1s', max_retries = 4 }
rpc_timeout = '10s'
account_prefix = 'mantra'
key_name = 'rly01-mantra-a'
store_prefix = 'ibc'
max_gas = 6000000
gas_price = { price = 0.01, denom = 'uom' }
gas_multiplier = 1.5
clock_drift = '1m' # to accommodate docker containers
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }


[[chains]]
id = '$MANTRA_B_E2E_CHAIN_ID'
rpc_addr = 'http://$MANTRA_B_E2E_VAL_HOST:26657'
grpc_addr = 'http://$MANTRA_B_E2E_VAL_HOST:9090'
event_source = { mode = 'pull', interval = '1s', max_retries = 4 }
rpc_timeout = '10s'
account_prefix = 'mantra'
key_name = 'rly01-mantra-b'
store_prefix = 'ibc'
max_gas = 6000000
gas_price = { price = 0.01, denom = 'uom' }
gas_multiplier = 1.5
clock_drift = '1m' # to accommodate docker containers
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }
EOF

# import keys
hermes keys add --key-name rly01-mantra-b --chain $MANTRA_B_E2E_CHAIN_ID --mnemonic-file /root/.hermes/MANTRA_B_E2E_RLY_MNEMONIC.txt
sleep 5
hermes keys add --key-name rly01-mantra-a --chain $MANTRA_A_E2E_CHAIN_ID --mnemonic-file /root/.hermes/MANTRA_A_E2E_RLY_MNEMONIC.txt
