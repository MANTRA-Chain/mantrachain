#!/bin/bash

CHAINID="mantrachain_7000-1"
MONIKER="localtestnet"
HOMEDIR="$HOME/.mantrachain"

KEYS_VAL[0]="validator1"
KEYS_VAL[1]="validator2"
KEYS_VAL[2]="validator3"
KEY_RECP_1="recipient1"
KEY_ADM_1="admin1"

KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"

GAS_ADJ=1.4
GAS_PRICE=0.0001axom

# Path variables
GENESIS_1=$HOMEDIR/${KEYS_VAL[0]}/config/genesis.json
TMP_GENESIS_1=$HOMEDIR/${KEYS_VAL[0]}/config/tmp_genesis.json

source "$PWD"/scripts/common.sh

echo "Stop the validators if any"
pkill -f "mantrachain*"

set -e # exit on first error

make build

# User prompt if an existing local node configuration is found.
if [ -d "$HOMEDIR" ]; then
	printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start new local nodes with new keys from genesis. When declined, the existing local nodes are started. \n" "$HOMEDIR"
	cecho "RED" "Overwrite the existing configuration and start new local nodes? [y/N]"
	read -r overwrite
else
	overwrite="Y"
fi

# Setup local nodes if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
  rm -rf $HOMEDIR

  mkdir $HOMEDIR

  cecho "GREEN" "Make mantrachain directories"
  for KEY in "${KEYS_VAL[@]}"; do
    mkdir $HOMEDIR/$KEY
  done

  cecho "GREEN" "Set clients config"
  for KEY in "${KEYS_VAL[@]}"; do
    "$PWD"/build/mantrachaind config keyring-backend $KEYRING --home "$HOMEDIR/$KEY"
    "$PWD"/build/mantrachaind config chain-id $CHAINID --home "$HOMEDIR/$KEY"
  done

  cecho "GREEN" "Create keys for all three validators"
  {
    "$PWD"/build/mantrachaind keys add "${KEYS_VAL[0]}" --keyring-backend=$KEYRING --algo $KEYALGO --home=$HOMEDIR/${KEYS_VAL[0]}
    "$PWD"/build/mantrachaind keys add "$KEY_RECP_1" --keyring-backend=$KEYRING --algo $KEYALGO --home=$HOMEDIR/${KEYS_VAL[0]}
    "$PWD"/build/mantrachaind keys add "$KEY_ADM_1" --keyring-backend=$KEYRING --algo $KEYALGO --home=$HOMEDIR/${KEYS_VAL[0]}
  } 2>&1 | tee accounts1.txt
  {
    "$PWD"/build/mantrachaind keys add "${KEYS_VAL[1]}" --keyring-backend=$KEYRING --algo $KEYALGO --home=$HOMEDIR/${KEYS_VAL[1]}
  } 2>&1 | tee accounts2.txt
  {
    "$PWD"/build/mantrachaind keys add "${KEYS_VAL[2]}" --keyring-backend=$KEYRING --algo $KEYALGO --home=$HOMEDIR/${KEYS_VAL[2]}
  } 2>&1 | tee accounts3.txt

  VALIDATOR_1_WALLET=$("$PWD"/build/mantrachaind keys show "${KEYS_VAL[0]}" -a --keyring-backend $KEYRING --home "$HOMEDIR/${KEYS_VAL[0]}")
  VALIDATOR_2_WALLET=$("$PWD"/build/mantrachaind keys show "${KEYS_VAL[1]}" -a --keyring-backend $KEYRING --home "$HOMEDIR/${KEYS_VAL[1]}")
  VALIDATOR_3_WALLET=$("$PWD"/build/mantrachaind keys show "${KEYS_VAL[2]}" -a --keyring-backend $KEYRING --home "$HOMEDIR/${KEYS_VAL[2]}")
  RECIPIENT_1_WALLET=$("$PWD"/build/mantrachaind keys show "$KEY_RECP_1" -a --keyring-backend $KEYRING --home "$HOMEDIR/${KEYS_VAL[0]}")
  ADMIN_1_WALLET=$("$PWD"/build/mantrachaind keys show "$KEY_ADM_1" -a --keyring-backend $KEYRING --home "$HOMEDIR/${KEYS_VAL[0]}")

  cecho "GREEN" "Init all three validators"
  for KEY in "${KEYS_VAL[@]}"; do
    "$PWD"/build/mantrachaind init $MONIKER-$KEY -o --chain-id=$CHAINID --home "$HOMEDIR/$KEY"
  done

  cecho "GREEN" "Replace genesis denom with xom"
  jq '.app_state["staking"]["params"]["bond_denom"]="axom"' "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"
  jq '.app_state["evm"]["params"]["evm_denom"]="axom"' "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"
  jq '.app_state["mint"]["params"]["mint_denom"]="axom"' "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"
  jq '.app_state["crisis"]["constant_fee"]["denom"]="axom"' "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"
  jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="axom"' "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"

  cecho "GREEN" "Update denom metadata"
  jq '.app_state["bank"]["denom_metadata"]=''[{"name":"xom","symbol":"XOM","description":"The native staking token of the Mantrachain.","denom_units":[{"denom":"axom","exponent":0,"aliases":["attoxom"]},{"denom":"uxom","exponent":12,"aliases":["microxom"]},{"denom":"xom","exponent":18}],"base":"axom","display":"xom"}]' "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"

  cecho "GREEN" "Create validator node with tokens to transfer to the node"
  "$PWD"/build/mantrachaind add-genesis-account $VALIDATOR_1_WALLET 10000000000000000000000000000axom --home "$HOMEDIR/${KEYS_VAL[0]}"
  "$PWD"/build/mantrachaind gentx ${KEYS_VAL[0]} 1000000000000000000000axom --keyring-backend=$KEYRING --chain-id=$CHAINID --home "$HOMEDIR/${KEYS_VAL[0]}"

  VALIDATOR1_ADDRESS=$(cat $HOMEDIR/${KEYS_VAL[0]}/config/gentx/$(ls $HOMEDIR/${KEYS_VAL[0]}/config/gentx -AU | head -1) | jq '.body["messages"][0].validator_address')

  cecho "GREEN" "Update genesis"
  jq '.app_state["vault"]["params"]["staking_validator_address"]='$VALIDATOR1_ADDRESS "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"
   jq '.app_state["vault"]["params"]["admin_account"]='\"$ADMIN_1_WALLET\" "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"
  jq '.app_state["bridge"]["params"]["admin_account"]='\"$ADMIN_1_WALLET\" "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"
  jq '.consensus_params["block"]["max_gas"]="10000000"' "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"
  jq '.app_state["feemarket"]["params"]["min_gas_price"]="0.000000000000000000"' "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"
  jq '.app_state["staking"]["params"]["unbonding_time"]="240s"' "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"
  jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' "$GENESIS_1" > "$TMP_GENESIS_1" && mv "$TMP_GENESIS_1" "$GENESIS_1"

  cecho "GREEN" "Validate genesis and collect genesis tx"
  "$PWD"/build/mantrachaind collect-gentxs --home "$HOMEDIR/${KEYS_VAL[0]}"
  "$PWD"/build/mantrachaind validate-genesis --home "$HOMEDIR/${KEYS_VAL[0]}"

  cecho "GREEN" "Copy validator1 genesis file to validator2-3"
  cp $HOMEDIR/${KEYS_VAL[0]}/config/genesis.json $HOMEDIR/${KEYS_VAL[1]}/config/genesis.json
  cp $HOMEDIR/${KEYS_VAL[0]}/config/genesis.json $HOMEDIR/${KEYS_VAL[2]}/config/genesis.json

  cecho "GREEN" "port key (validator1 uses default ports)"
  cecho "GREEN" "validator1 1317, 9090, 9091, 26658, 26657, 26656, 6060, 8545, 8546"
  cecho "GREEN" "validator2 1316, 9088, 9089, 26655, 26654, 26653, 6061, 8543, 85464"
  cecho "GREEN" "validator3 1315, 9086, 9087, 26652, 26651, 26650, 6062, 8541, 85462"

  cecho "GREEN" "Change app.toml values"
  for KEY in "${KEYS_VAL[@]}"; do
    sed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $HOMEDIR/$KEY/config/app.toml
    sed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' $HOMEDIR/$KEY/config/app.toml
    sed -i -E '1,/enable = false/s|enable = false|enable = true|g' $HOMEDIR/$KEY/config/app.toml
    sed -i -E 's|minimum-gas-prices = \"\"|minimum-gas-prices = \"'$GAS_PRICE'\"|g' $HOMEDIR/$KEY/config/app.toml
  done
  sed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1316|g' $HOMEDIR/${KEYS_VAL[1]}/config/app.toml
  sed -i -E 's|0.0.0.0:9090|0.0.0.0:9088|g' $HOMEDIR/${KEYS_VAL[1]}/config/app.toml
  sed -i -E 's|0.0.0.0:9091|0.0.0.0:9089|g' $HOMEDIR/${KEYS_VAL[1]}/config/app.toml
  sed -i -E 's|0.0.0.0:8545|0.0.0.0:8543|g' $HOMEDIR/${KEYS_VAL[1]}/config/app.toml
  sed -i -E 's|0.0.0.0:8546|0.0.0.0:8544|g' $HOMEDIR/${KEYS_VAL[1]}/config/app.toml
  sed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1315|g' $HOMEDIR/${KEYS_VAL[2]}/config/app.toml
  sed -i -E 's|0.0.0.0:9090|0.0.0.0:9086|g' $HOMEDIR/${KEYS_VAL[2]}/config/app.toml
  sed -i -E 's|0.0.0.0:9091|0.0.0.0:9087|g' $HOMEDIR/${KEYS_VAL[2]}/config/app.toml
  sed -i -E 's|0.0.0.0:8545|0.0.0.0:8541|g' $HOMEDIR/${KEYS_VAL[2]}/config/app.toml
  sed -i -E 's|0.0.0.0:8546|0.0.0.0:8542|g' $HOMEDIR/${KEYS_VAL[2]}/config/app.toml

  cecho "GREEN" "Change config.toml values"
  sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOMEDIR/${KEYS_VAL[0]}/config/config.toml
  sed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $HOMEDIR/${KEYS_VAL[0]}/config/config.toml
  sed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $HOMEDIR/${KEYS_VAL[1]}/config/config.toml
  sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26655|g' $HOMEDIR/${KEYS_VAL[1]}/config/config.toml
  sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26654|g' $HOMEDIR/${KEYS_VAL[1]}/config/config.toml
  sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26653|g' $HOMEDIR/${KEYS_VAL[1]}/config/config.toml
  sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOMEDIR/${KEYS_VAL[1]}/config/config.toml
  sed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $HOMEDIR/${KEYS_VAL[2]}/config/config.toml
  sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26652|g' $HOMEDIR/${KEYS_VAL[2]}/config/config.toml
  sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26651|g' $HOMEDIR/${KEYS_VAL[2]}/config/config.toml
  sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26650|g' $HOMEDIR/${KEYS_VAL[2]}/config/config.toml
  sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOMEDIR/${KEYS_VAL[2]}/config/config.toml
  sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$("$PWD"/build/mantrachaind tendermint show-node-id --home=$HOMEDIR/${KEYS_VAL[0]})@localhost:26656\"|g" $HOMEDIR/${KEYS_VAL[1]}/config/config.toml
  sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$("$PWD"/build/mantrachaind tendermint show-node-id --home=$HOMEDIR/${KEYS_VAL[0]})@localhost:26656\"|g" $HOMEDIR/${KEYS_VAL[2]}/config/config.toml

  cecho "GREEN" "Change client.toml values"
  sed -i -E 's|chain-id = \"\"|chain-id = \"mantrachain\"|g' $HOMEDIR/${KEYS_VAL[0]}/config/client.toml
  sed -i -E 's|chain-id = \"\"|chain-id = \"mantrachain\"|g' $HOMEDIR/${KEYS_VAL[1]}/config/client.toml
  sed -i -E 's|26657|26654|g' $HOMEDIR/${KEYS_VAL[1]}/config/client.toml
  sed -i -E 's|chain-id = \"\"|chain-id = \"mantrachain\"|g' $HOMEDIR/${KEYS_VAL[2]}/config/client.toml
  sed -i -E 's|26657|26651|g' $HOMEDIR/${KEYS_VAL[2]}/config/client.toml
fi

cecho "GREEN" "Start all three validators"
for KEY in "${KEYS_VAL[@]}"; do
  tmux new -s "$KEY" -d "$PWD"/build/mantrachaind start --metrics --pruning=nothing --log_level $LOGLEVEL --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --home "$HOMEDIR/$KEY"
done

if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
  cecho "GREEN" "Send axom from ${KEYS_VAL[0]} to ${KEYS_VAL[1]}"
  sleep 7
  "$PWD"/build/mantrachaind tx bank send ${KEYS_VAL[0]} $VALIDATOR_2_WALLET 100000000000000000000axom --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes  --home "$HOMEDIR/${KEYS_VAL[0]}"
  cecho "GREEN" "Send axom from ${KEYS_VAL[0]} to ${KEYS_VAL[2]}"
  sleep 7
  "$PWD"/build/mantrachaind tx bank send ${KEYS_VAL[0]} $VALIDATOR_3_WALLET 100000000000000000000axom --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes  --home "$HOMEDIR/${KEYS_VAL[0]}"

  cecho "GREEN" "Create second and third validator"
  sleep 7
  "$PWD"/build/mantrachaind tx staking create-validator --amount=500000000000axom --from="${KEYS_VAL[1]}" --pubkey=$("$PWD"/build/mantrachaind tendermint show-validator --home=$HOMEDIR/${KEYS_VAL[1]}) --moniker="$MONIKER-${KEYS_VAL[1]}" --chain-id="$CHAINID" --commission-rate="0.1" --commission-max-rate="0.2" --commission-max-change-rate="0.05" --min-self-delegation="5000000000" --keyring-backend=test --home=$HOMEDIR/${KEYS_VAL[1]} --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes
  sleep 7
  "$PWD"/build/mantrachaind tx staking create-validator --amount=500000000000axom --from=${KEYS_VAL[2]} --pubkey=$("$PWD"/build/mantrachaind tendermint show-validator --home=$HOMEDIR/${KEYS_VAL[2]}) --moniker="$MONIKER-${KEYS_VAL[2]}" --chain-id="$CHAINID" --commission-rate="0.1" --commission-max-rate="0.2" --commission-max-change-rate="0.05" --min-self-delegation="5000000000" --keyring-backend=test --home=$HOMEDIR/${KEYS_VAL[2]} --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes

  cecho "GREEN" "Send axom from ${KEYS_VAL[0]} to $KEY_RECP_1"
  sleep 7
  "$PWD"/build/mantrachaind tx bank send $VALIDATOR_1_WALLET $RECIPIENT_1_WALLET 10000000000000000000axom --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes  --home "$HOMEDIR/${KEYS_VAL[0]}"

  cecho "GREEN" "Send axom from ${KEYS_VAL[0]} to $KEY_ADM_1"
  sleep 7
  "$PWD"/build/mantrachaind tx bank send $VALIDATOR_1_WALLET $ADMIN_1_WALLET 10000000000000000000axom --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes  --home "$HOMEDIR/${KEYS_VAL[0]}"
fi