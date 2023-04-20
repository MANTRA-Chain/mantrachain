#!/bin/bash

CHAINID="mantrachain"
MONIKER="mantrachaindevnet"
HOMEDIR="$HOME/.mantrachain"

KEYS[0]="validator"
KEYS[1]="recipient"
KEYS[2]="admin"
ACCOUNT_PRIVILEGES_GUARD_NFT_COLLECTION_ID="account_privileges_guard_nft_collection"

KEYRING="test"
LOGLEVEL="info"

GAS_ADJ=2
GAS_PRICE=0.0001uaum

# Path variables
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
CLIENT_TOML=$HOMEDIR/config/client.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

# Ccv provider variables
PROV_NODE_DIR=~/.ccv-provider
PROV_NODE_MONIKER=ccv-provider
PROV_CHAIN_ID=provider
PROV_KEY=provider-key

command -v hermes >/dev/null 2>&1 || {
  echo >&2 "hermes not installed. More info: https://github.com/informalsystems/hermes"
  exit 1
}

source "$PWD"/scripts/common.sh

echo "Stop the consumer if any"
pkill -f mantrachain
pkill -f hermes

set -e # exit on first error

make build

# User prompt if an existing local node configuration is found.
if [ -d "$HOMEDIR" ]; then
  printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start a new local node with new keys from genesis. When declined, the existing local node is started. \n" "$HOMEDIR"
  cecho "RED" "CCV CONSUMER: Overwrite the existing configuration and start a new local node? [y/N]"
  read -r overwrite
else
  overwrite="Y"
fi

# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
  rm -rf $HOMEDIR

  mkdir $HOMEDIR

  cecho "CYAN" "Set client config"
  "$PWD"/build/mantrachaind config keyring-backend $KEYRING --home "$HOMEDIR"
  "$PWD"/build/mantrachaind config chain-id $CHAINID --home "$HOMEDIR"

  cecho "CYAN" "Create keys"
  for KEY in "${KEYS[@]}"; do
    "$PWD"/build/mantrachaind keys add $KEY --keyring-backend $KEYRING --home "$HOMEDIR" --output json >${HOMEDIR}/${KEY}.json 2>&1
  done

  VALIDATOR_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[0]} -a --keyring-backend $KEYRING --home "$HOMEDIR")
  RECIPIENT_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[1]} -a --keyring-backend $KEYRING --home "$HOMEDIR")
  ADMIN_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[2]} -a --keyring-backend $KEYRING --home "$HOMEDIR")

  cecho "CYAN" "Init the validator"
  "$PWD"/build/mantrachaind init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"

  cecho "CYAN" "Replace genesis denom with uaum"
  jq '.app_state["crisis"]["constant_fee"]["denom"]="uaum"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["lpfarm"]["params"]["private_plan_creation_fee"][0]["denom"]="uaum"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["liquidity"]["params"]["pair_creation_fee"][0]["denom"]="uaum"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["liquidity"]["params"]["pool_creation_fee"][0]["denom"]="uaum"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

  cecho "CYAN" "Update genesis"
  jq '.app_state["bank"]["denom_metadata"]=''[{"name":"aum","symbol":"AUM","description":"The native staking token of the Mantrachain.","denom_units":[{"denom":"uaum","exponent":0,"aliases":["microaum"]},{"denom":"maum","exponent":3,"aliases":["milliaum"]},{"denom":"aum","exponent":6}],"base":"uaum","display":"aum"}]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["guard"]["params"]["admin_account"]='\"$ADMIN_WALLET\" "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["guard"]["params"]["account_privileges_token_collection_creator"]='\"$ADMIN_WALLET\" "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["guard"]["params"]["account_privileges_token_collection_id"]='\"$ACCOUNT_PRIVILEGES_GUARD_NFT_COLLECTION_ID\" "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

  cecho "CYAN" "Create validator node with tokens to transfer to the node"
  "$PWD"/build/mantrachaind add-genesis-account $VALIDATOR_WALLET 10000000000000000000000000uaum --home "$HOMEDIR"

  cecho "CYAN" "Validate genesis"
  "$PWD"/build/mantrachaind validate-genesis --home "$HOMEDIR"

  cecho "CYAN" "Change app.toml values"
  sed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $APP_TOML
  sed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' $APP_TOML
  sed -i -E '112s/.*/enable = true/' $APP_TOML
  sed -i -E 's|minimum-gas-prices = \"0stake\"|minimum-gas-prices = \"'$GAS_PRICE'\"|g' $APP_TOML

  cecho "CYAN" "Change config.toml values"
  sed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $CONFIG
  sed -i -E 's|pprof_laddr = \"localhost:6060\"|pprof_laddr = \"localhost:6050\"|g' $CONFIG

  cecho "CYAN" "Change client.toml values"
  sed -i -E 's|chain-id = \"\"|chain-id = \"'$CHAINID'\"|g' $CLIENT_TOML
  sed -i -r "/node =/ s/= .*/= \"tcp:\/\/127.0.0.1:26648\"/" $CLIENT_TOML

  cecho "CYAN" "Create consumer proposal"
  tee ${PROV_NODE_DIR}/consumer-proposal.json <<EOF
{
  "title": "Create mantrachain",
  "description": "Mantrachain shared security",
  "chain_id": "mantrachain", 
  "initial_height": {
      "revision_height": 1
  },
  "genesis_hash": "Z2VuX2hhc2g=",
  "binary_hash": "YmluX2hhc2g=",
  "spawn_time": "2022-01-27T15:59:50.121607-08:00",
  "blocks_per_distribution_transmission": 1000,
  "consumer_redistribution_fraction": "0.75",
  "historical_entries": 10000,
  "transfer_timeout_period": 3600000000000,
  "ccv_timeout_period": 2419200000000000,
  "unbonding_period": 1728000000000000,
  "deposit": "10000stake"
}
EOF

  cecho "CYAN" "Wait 20 seconds for provider to init..."
  sleep 20

  cecho "CYAN" "Submit consumer proposal"
  interchain-security-pd tx gov submit-proposal \
    consumer-addition ${PROV_NODE_DIR}/consumer-proposal.json \
    --keyring-backend test \
    --chain-id $PROV_CHAIN_ID \
    --from $PROV_KEY \
    --home $PROV_NODE_DIR \
    -b block -y

  sleep 7

  cecho "CYAN" "Make deposit to consumer proposal"
  interchain-security-pd tx gov deposit 1 100000000stake \
    --keyring-backend test \
    --chain-id $PROV_CHAIN_ID \
    --from $PROV_KEY \
    --home $PROV_NODE_DIR \
    -b block -y

  sleep 7

  cecho "CYAN" "Vote on consumer proposal"
  interchain-security-pd tx gov vote 1 yes --from $PROV_KEY \
    --keyring-backend test --chain-id $PROV_CHAIN_ID --home $PROV_NODE_DIR -b block -y

  cecho "CYAN" "Wait 30 seconds for the consumer proposal to be processed..."
  sleep 30

  cecho "CYAN" "Get the genesis consumer chain state from the provider chain"
  interchain-security-pd query provider consumer-genesis $CHAINID --home $PROV_NODE_DIR -o json >ccvconsumer_genesis.json
  jq -s '.[0].app_state.ccvconsumer = .[1] | .[0]' ${HOMEDIR}/config/genesis.json ccvconsumer_genesis.json >${HOMEDIR}/edited_genesis.json
  mv ${HOMEDIR}/edited_genesis.json ${HOMEDIR}/config/genesis.json
  rm ccvconsumer_genesis.json

  cecho "CYAN" "Copy the validator keypair"
  echo '{"height": "0","round": 0,"step": 0}' >${HOMEDIR}/data/priv_validator_state.json
  cp ${PROV_NODE_DIR}/config/priv_validator_key.json ${HOMEDIR}/config/priv_validator_key.json
  cp ${PROV_NODE_DIR}/config/node_key.json ${HOMEDIR}/config/node_key.json
fi

COORDINATOR_P2P_ADDRESS=$(jq -r '.app_state.genutil.gen_txs[0].body.memo' ${PROV_NODE_DIR}/config/genesis.json)
CONS_P2P_ADDRESS=$(echo $COORDINATOR_P2P_ADDRESS | sed 's/@.*/@127.0.0.1:26646/')

cecho "CYAN" "Start the validator"
tmux new -s mantrachain -d "$PWD"/build/mantrachaind start \
  --pruning=nothing \
  --log_level=$LOGLEVEL \
  --home=$HOMEDIR \
  --rpc.laddr tcp://127.0.0.1:26648 \
  --grpc.address 127.0.0.1:9081 \
  --address tcp://127.0.0.1:26645 \
  --p2p.laddr tcp://127.0.0.1:26646 \
  --grpc-web.enable=true \
  --p2p.persistent_peers $CONS_P2P_ADDRESS

if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
  cecho "CYAN" "Setup IBC-Relayer"
  tee ~/.hermes/config.toml <<EOF
[global]
log_level = "info"

[[chains]]
account_prefix = "mantrachain"
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
grpc_addr = "tcp://127.0.0.1:9091"
id = "provider"
key_name = "relayer"
max_gas = 2000000
rpc_addr = "http://127.0.0.1:26658"
rpc_timeout = "10s"
store_prefix = "ibc"
trusting_period = "14days"
websocket_addr = "ws://127.0.0.1:26658/websocket"

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
  hermes keys add --key-file ${HOMEDIR}/${KEYS[0]}.json --chain mantrachain
  hermes keys add --key-file ${PROV_NODE_DIR}/${PROV_KEY}.json --chain provider

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
fi

cecho "CYAN" "Start Hermes"
tmux new -s hermes -d hermes --json start

if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
  cecho "CYAN" "Test the CCV protocol: Get provider delegations"
  DELEGATIONS=$(interchain-security-pd q staking delegations $(jq -r .address ${PROV_NODE_DIR}/${PROV_KEY}.json) --home $PROV_NODE_DIR -o json)

  cecho "CYAN" "Test the CCV protocol: Get provider operator address"
  OPERATOR_ADDR=$(echo $DELEGATIONS | jq -r '.delegation_responses[0].delegation.validator_address')

  sleep 7

  cecho "CYAN" "Test the CCV protocol: Delegate tokens"
  interchain-security-pd tx staking delegate $OPERATOR_ADDR 1000000stake \
    --from $PROV_KEY \
    --keyring-backend test \
    --home $PROV_NODE_DIR \
    --chain-id $PROV_CHAIN_ID \
    -y -b block

  cecho "CYAN" "Verify the chains validator-set: Query provider chain valset"
  sleep 7
  interchain-security-pd q tendermint-validator-set --home $PROV_NODE_DIR
  cecho "CYAN" "Verify the chains validator-set: Query consumer chain valset"
  sleep 14
  "$PWD"/build/mantrachaind q tendermint-validator-set --home $HOMEDIR
fi

echo
cecho "YELLOW" "To init sample:"
cecho "GREEN" "./scripts/init-sample.sh"
echo
cecho "YELLOW" "To track hermes:"
cecho "GREEN" "tmux a -t hermes"
echo
cecho "YELLOW" "To track mantrachain:"
cecho "GREEN" "tmux a -t mantrachain"
echo
cecho "YELLOW" "To track provider:"
cecho "GREEN" "tmux a -t provider"
echo
