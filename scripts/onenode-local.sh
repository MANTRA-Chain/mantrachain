#!/bin/bash

CHAINID="mantrachain-local"
MONIKER="localtestnet"
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
      "$PWD"/build/mantrachaind keys add $KEY --keyring-backend $KEYRING --home "$HOMEDIR"
    done
  } 2>&1 | tee accounts.txt

  VALIDATOR_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[0]} -a --keyring-backend $KEYRING --home "$HOMEDIR")
  RECIPIENT_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[1]} -a --keyring-backend $KEYRING --home "$HOMEDIR")
  ADMIN_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[2]} -a --keyring-backend $KEYRING --home "$HOMEDIR")

  cecho "GREEN" "Init the validator"
	"$PWD"/build/mantrachaind init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"
  
  cecho "GREEN" "Replace genesis denom with aum"
  jq '.app_state["staking"]["params"]["bond_denom"]="uaum"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["mint"]["params"]["mint_denom"]="uaum"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["crisis"]["constant_fee"]["denom"]="uaum"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="uaum"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["lpfarm"]["params"]["private_plan_creation_fee"]["denom"]="uaum"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["liquidity"]["params"]["pair_creation_fee"][0]["denom"]="uaum"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["liquidity"]["params"]["pool_creation_fee"][0]["denom"]="uaum"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

  cecho "GREEN" "Update denom metadata"
  jq '.app_state["bank"]["denom_metadata"]=''[{"name":"aum","symbol":"AUM","description":"The native staking token of the Mantrachain.","denom_units":[{"denom":"uaum","exponent":0,"aliases":["microaum"]},{"denom":"maum","exponent":3,"aliases":["milliaum"]},{"denom":"aum","exponent":6}],"base":"uaum","display":"aum"}]' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

  cecho "GREEN" "Create validator node with tokens to transfer to the node"
  "$PWD"/build/mantrachaind add-genesis-account $VALIDATOR_WALLET 100000000000000000uaum --home "$HOMEDIR"
  "$PWD"/build/mantrachaind gentx ${KEYS[0]} 100000000000000uaum --keyring-backend=$KEYRING --chain-id=$CHAINID --home "$HOMEDIR"

  cecho "GREEN" "Update genesis"
  jq '.app_state["guard"]["params"]["admin_account"]='\"$ADMIN_WALLET\" "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["staking"]["params"]["unbonding_time"]="240s"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["guard"]["params"]["account_privileges_token_collection_creator"]='\"$ADMIN_WALLET\" "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
  jq '.app_state["guard"]["params"]["account_privileges_token_collection_id"]='\"$ACCOUNT_PRIVILEGES_GUARD_NFT_COLLECTION_ID\" "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

  cecho "GREEN" "Collect genesis tx"
  "$PWD"/build/mantrachaind collect-gentxs --home "$HOMEDIR"

  cecho "GREEN" "Validate genesis"
  "$PWD"/build/mantrachaind validate-genesis --home "$HOMEDIR"

  cecho "GREEN" "Port key (validator uses default ports)"
  cecho "GREEN" "validator 1317, 9090, 9091, 26658, 26657, 26656, 6060"

  cecho "GREEN" "Change app.toml values"
  sed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $APP_TOML
  sed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' $APP_TOML
  sed -i -E '112s/.*/enable = true/' $APP_TOML
  sed -i -E 's|minimum-gas-prices = \"0uaum\"|minimum-gas-prices = \"'$GAS_PRICE'\"|g' $APP_TOML

  cecho "GREEN" "Change config.toml values"
  sed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $CONFIG

  cecho "GREEN" "Change client.toml values"
  sed -i -E 's|chain-id = \"\"|chain-id = \"'$CHAINID'\"|g' $CLIENT_TOML
fi

cecho "GREEN" "Start the validator"
tmux new -s mantrachain -d "$PWD"/build/mantrachaind start --pruning=nothing --log_level $LOGLEVEL --home "$HOMEDIR"

if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
  sleep 7

  cecho "GREEN" "Send uaum from ${KEYS[0]} to ${KEYS[1]}"
  "$PWD"/build/mantrachaind tx bank send $VALIDATOR_WALLET $RECIPIENT_WALLET 100000000000000uaum --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes --home "$HOMEDIR"
  
  sleep 7

  cecho "GREEN" "Send uaum from ${KEYS[0]} to ${KEYS[2]}"
  "$PWD"/build/mantrachaind tx bank send $VALIDATOR_WALLET $ADMIN_WALLET 100000000000000uaum --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes --home "$HOMEDIR"
  
  sleep 7
  
  ACCOUNT_PRIVILEGES_GUARD_COLL_JSON=$(echo '{"id":"{id}","soul_bonded_nfts": true,"restricted_nfts":true,"name":"AccountPrivilegesCollection","description":"AccountPrivilegesCollection","category":"utility"}' | sed -e "s/{id}/$ACCOUNT_PRIVILEGES_GUARD_NFT_COLLECTION_ID/g")

  cecho "GREEN" "Create account privileges guard nft collection"
  "$PWD"/build/mantrachaind tx token create-nft-collection "$(echo $ACCOUNT_PRIVILEGES_GUARD_COLL_JSON)" --chain-id $CHAINID --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes --home "$HOMEDIR"
  
  sleep 7

  ACCOUNT_PRIVILEGES_GUARD_NFT_VALIDATOR_JSON=$(echo '{"id":"{id}","title":"AccountPrivileges","description":"AccountPrivileges"}' | sed -e "s/{id}/$VALIDATOR_WALLET/g")

  ACCOUNT_PRIVILEGES_GUARD_NFT_ADMIN_JSON=$(echo '{"id":"{id}","title":"AccountPrivileges","description":"AccountPrivileges"}' | sed -e "s/{id}/$ADMIN_WALLET/g")

  cecho "GREEN" "Mint account privileges guard nfts"
  "$PWD"/build/mantrachaind tx token mint-nft "$(echo $ACCOUNT_PRIVILEGES_GUARD_NFT_VALIDATOR_JSON)" --collection-creator $ADMIN_WALLET --collection-id $ACCOUNT_PRIVILEGES_GUARD_NFT_COLLECTION_ID --chain-id $CHAINID --from ${KEYS[2]} --receiver $VALIDATOR_WALLET --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes

  sleep 7

  "$PWD"/build/mantrachaind tx token mint-nft "$(echo $ACCOUNT_PRIVILEGES_GUARD_NFT_ADMIN_JSON)" --collection-creator $ADMIN_WALLET --collection-id $ACCOUNT_PRIVILEGES_GUARD_NFT_COLLECTION_ID --chain-id $CHAINID --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes

  sleep 7

  cecho "GREEN" "Update guard transfer coins"
  "$PWD"/build/mantrachaind tx guard update-guard-transfer-coins true --chain-id $CHAINID --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes
fi

cecho "GREEN" "Track logs"
tmux a -t mantrachain