#!/bin/bash

CHAINID="mantrachain_7000-1"
MONIKER="localtestnet"
HOMEDIR="$HOME/.mantrachain"

KEYS[0]="validator"
KEYS[1]="recipient"
KEYS[2]="admin"
GUARD_NFT_COLLECTION_ID="guardnft"

KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"

GAS_ADJ=1.4
GAS_PRICE=0.0001axom

# Path variables
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
CLIENT_TOML=$HOMEDIR/config/client.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

source "$PWD"/scripts/common.sh

echo "Stop the validators if any"
pkill -f "mantrachain*"

set -e # exit on first error

make build

# User prompt if an existing local node configuration is found.
if [ -d "$HOMEDIR" ]; then
	printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start a new local node with new keys from genesis. When declined, the existing local node is started. \n" "$HOMEDIR"
	cecho "RED" "Overwrite the existing configuration and start a new local node? [y/N]"
	read -r overwrite
else
	overwrite="Y"
fi

# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
  rm -rf $HOMEDIR

  mkdir $HOMEDIR
  
  cecho "GREEN" "Set client config"
  "$PWD"/build/mantrachaind config keyring-backend $KEYRING --home "$HOMEDIR"
	"$PWD"/build/mantrachaind config chain-id $CHAINID --home "$HOMEDIR"

  cecho "GREEN" "Create keys"
  {
    for KEY in "${KEYS[@]}"; do
      "$PWD"/build/mantrachaind keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
    done
  } 2>&1 | tee accounts.txt

  VALIDATOR_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[0]} -a --keyring-backend $KEYRING --home "$HOMEDIR")
  RECIPIENT_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[1]} -a --keyring-backend $KEYRING --home "$HOMEDIR")
  ADMIN_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[2]} -a --keyring-backend $KEYRING --home "$HOMEDIR")

  cecho "GREEN" "Init the validator"
	"$PWD"/build/mantrachaind init $MONIKER-$KEY -o --chain-id $CHAINID --home "$HOMEDIR"

  cecho "GREEN" "Replace genesis denom with xom"
  jq '.app_state["staking"]["params"]["bond_denom"]="axom"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["evm"]["params"]["evm_denom"]="axom"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["mint"]["params"]["mint_denom"]="axom"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["crisis"]["constant_fee"]["denom"]="axom"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="axom"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

  cecho "GREEN" "Update denom metadata"
  jq '.app_state["bank"]["denom_metadata"]=''[{"name":"xom","symbol":"XOM","description":"The native staking token of the Mantrachain.","denom_units":[{"denom":"axom","exponent":0,"aliases":["attoxom"]},{"denom":"uxom","exponent":12,"aliases":["microxom"]},{"denom":"xom","exponent":18}],"base":"axom","display":"xom"}]' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

  cecho "GREEN" "Create validator node with tokens to transfer to the node"
  "$PWD"/build/mantrachaind add-genesis-account $VALIDATOR_WALLET 10000000000000000000000000000axom --home "$HOMEDIR"
  "$PWD"/build/mantrachaind gentx ${KEYS[0]} 1000000000000000000000axom --keyring-backend=$KEYRING --chain-id=$CHAINID --home "$HOMEDIR"

  VALIDATOR_ADDRESS=$(cat $HOMEDIR/config/gentx/$(ls $HOMEDIR/config/gentx | head -1) | jq '.body["messages"][0].validator_address')

  GUARD_ACC_PERM_LIST_JSON=$(echo '[{"cat":"0","creator":"{admin}","whl_curr":["*"]}]' | sed -e "s/{validator}/$VALIDATOR_WALLET/g" | sed -e "s/{admin}/$ADMIN_WALLET/g")

  GUARD_TRANSFER_JSON=$(echo '{"enabled":false,"creator":"{admin}"}' | sed -e "s/{admin}/$ADMIN_WALLET/g")

  cecho "GREEN" "Update genesis"
  jq '.consensus_params["block"]["max_gas"]="10000000"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["feemarket"]["params"]["min_gas_price"]="0.000000000000000000"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["vault"]["params"]["staking_validator_address"]='$VALIDATOR_ADDRESS "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["vault"]["params"]["admin_account"]='\"$ADMIN_WALLET\" "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["bridge"]["params"]["admin_account"]='\"$ADMIN_WALLET\" "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["guard"]["params"]["admin_account"]='\"$ADMIN_WALLET\" "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["staking"]["params"]["unbonding_time"]="240s"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["guard"]["params"]["token_collection_creator"]='\"$ADMIN_WALLET\" "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["guard"]["params"]["token_collection_id"]='\"$GUARD_NFT_COLLECTION_ID\" "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["guard"]["acc_perm_list"]='$GUARD_ACC_PERM_LIST_JSON "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["guard"]["guard_transfer"]='$GUARD_TRANSFER_JSON "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

  cecho "GREEN" "Collect genesis tx"
  "$PWD"/build/mantrachaind collect-gentxs --home "$HOMEDIR"

  cecho "GREEN" "Validate genesis"
  "$PWD"/build/mantrachaind validate-genesis --home "$HOMEDIR"

  cecho "GREEN" "Port key (validator uses default ports)"
  cecho "GREEN" "validator 1317, 9090, 9091, 26658, 26657, 26656, 6060"

  cecho "GREEN" "Change app.toml values"
  sed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $APP_TOML
  sed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' $APP_TOML
  sed -i -E '117s/.*/enable = true/' $APP_TOML
  sed -i -E 's|minimum-gas-prices = \"0axom\"|minimum-gas-prices = \"'$GAS_PRICE'\"|g' $APP_TOML

  cecho "GREEN" "Change config.toml values"
  sed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $CONFIG

  cecho "GREEN" "Change client.toml values"
  sed -i -E 's|chain-id = \"\"|chain-id = \"'$CHAINID'\"|g' $CLIENT_TOML
fi

cecho "GREEN" "Start the validator"
tmux new -s mantrachain -d "$PWD"/build/mantrachaind start --metrics --pruning=nothing --log_level $LOGLEVEL --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --home "$HOMEDIR"

if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
  sleep 7

  cecho "GREEN" "Send axom from ${KEYS[0]} to ${KEYS[1]}"
  "$PWD"/build/mantrachaind tx bank send $VALIDATOR_WALLET $RECIPIENT_WALLET 100000000000000000000axom --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes  --home "$HOMEDIR"
  
  sleep 7

  cecho "GREEN" "Send axom from ${KEYS[0]} to ${KEYS[2]}"
  "$PWD"/build/mantrachaind tx bank send $VALIDATOR_WALLET $ADMIN_WALLET 10000000000000000000axom --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes  --home "$HOMEDIR"
  
  sleep 7
  
  GUARD_COLL_JSON=$(echo '{"id":"{id}","soul_bonded": true, "name":"GuardCollection","description":"GuardCollection","category":"utility"}' | sed -e "s/{id}/$GUARD_NFT_COLLECTION_ID/g")

  cecho "GREEN" "Create guard nft collection"
  "$PWD"/build/mantrachaind tx token create-nft-collection "$(echo $GUARD_COLL_JSON)" --chain-id $CHAINID --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes  --home "$HOMEDIR"

  sleep 7

  GUARD_NFT_VALIDATOR_JSON=$(echo '{"id":"{id}","title":"GuardNft","description":"GuardNft","attributes":[{"type":"AccPerm","value":"0"}]}' | sed -e "s/{id}/$VALIDATOR_WALLET/g")

  GUARD_NFT_ADMIN_JSON=$(echo '{"id":"{id}","title":"GuardNft","description":"GuardNft","attributes":[{"type":"AccPerm","value":"0"}]}' | sed -e "s/{id}/$ADMIN_WALLET/g")

  cecho "GREEN" "Mint guard nfts"
  "$PWD"/build/mantrachaind tx token mint-nft "$(echo $GUARD_NFT_VALIDATOR_JSON)" --collection-creator $ADMIN_WALLET --collection-id $GUARD_NFT_COLLECTION_ID --chain-id $CHAINID --from ${KEYS[2]} --receiver $VALIDATOR_WALLET --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes

  sleep 7

  "$PWD"/build/mantrachaind tx token mint-nft "$(echo $GUARD_NFT_ADMIN_JSON)" --collection-creator $ADMIN_WALLET --collection-id $GUARD_NFT_COLLECTION_ID --chain-id $CHAINID --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes

  # sleep 7

  # cecho "GREEN" "Update guard transfer"
  # "$PWD"/build/mantrachaind tx guard update-guard-transfer true --chain-id $CHAINID --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes
fi

cecho "GREEN" "Track logs"
tmux a -t mantrachain