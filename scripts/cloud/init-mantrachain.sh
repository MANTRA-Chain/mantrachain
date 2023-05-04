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

source "$PWD"/scripts/common.sh

set -e # exit on first error

make build

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

cecho "YELLOW" "You can start the chain with:
/build/mantrachaind start \
  --pruning=nothing \
  --log_level=$LOGLEVEL \
  --home=$HOMEDIR \
  --rpc.laddr tcp://127.0.0.1:26648 \
  --grpc.address 127.0.0.1:9081 \
  --address tcp://127.0.0.1:26645 \
  --p2p.laddr tcp://127.0.0.1:26646 \
  --grpc-web.enable=true \
  --p2p.persistent_peers <ID@IP>"