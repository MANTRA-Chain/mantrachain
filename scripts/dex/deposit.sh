#!/bin/bash

CHAINID="mantrachain"
HOMEDIR="$HOME/.mantrachain"

KEYS[0]="validator"
KEYS[1]="recipient"
KEYS[2]="admin"
ACCOUNT_PRIVILEGES_GUARD_NFT_COLLECTION_ID="account_privileges_guard_nft_collection"

KEYRING="test"
LOGLEVEL="info"

GAS_ADJ=2
GAS_PRICE=0.0001uaum



source "$PWD"/scripts/common.sh

set -e # exit on first error

VALIDATOR_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[0]} -a --keyring-backend $KEYRING --home "$HOMEDIR")
RECIPIENT_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[1]} -a --keyring-backend $KEYRING --home "$HOMEDIR")
ADMIN_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[2]} -a --keyring-backend $KEYRING --home "$HOMEDIR")

DENOM_A="a"
DENOM_B="b"
DENOM_A_FULL="factory/$ADMIN_WALLET/$DENOM_A"
DENOM_B_FULL="factory/$ADMIN_WALLET/$DENOM_B"

POOL_ID=$("$PWD"/build/mantrachaind q liquidity pools --output json | jq -r '.pools[-1].id')

for index in {1..3}
do
  cecho "CYAN" "Create user user$index"
  "$PWD"/build/mantrachaind keys add user$index --keyring-backend $KEYRING --home "$HOMEDIR" --output json
  sleep 7
  cecho "CYAN" "Send uaum from admin to user$index"
  "$PWD"/build/mantrachaind tx bank send $ADMIN_WALLET $("$PWD"/build/mantrachaind keys show user$index -a --keyring-backend $KEYRING --home "$HOMEDIR") 1000000uaum --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --chain-id $CHAINID --home "$HOMEDIR" -y
  sleep 7
  cecho "CYAN" "Send $DENOM_A_FULL from admin to user$index"
  "$PWD"/build/mantrachaind tx bank send $ADMIN_WALLET $("$PWD"/build/mantrachaind keys show user$index -a --keyring-backend $KEYRING --home "$HOMEDIR") 1000000$DENOM_A_FULL --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --chain-id $CHAINID --home "$HOMEDIR" -y
  sleep 7
  cecho "CYAN" "Send $DENOM_B_FULL from admin to user$index"
  "$PWD"/build/mantrachaind tx bank send $ADMIN_WALLET $("$PWD"/build/mantrachaind keys show user$index -a --keyring-backend $KEYRING --home "$HOMEDIR") 1000000$DENOM_B_FULL --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --chain-id $CHAINID --home "$HOMEDIR" -y
  sleep 7
  cecho "CYAN" "User user$index deposits to pool $POOL_ID"
  "$PWD"/build/mantrachaind tx liquidity deposit $(($POOL_ID)) "2000$DENOM_A_FULL,2000$DENOM_B_FULL" --from user$index --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --chain-id $CHAINID --home "$HOMEDIR" -y
done

"$PWD"/build/mantrachaind q bank balances $(./build/mantrachaind keys show admin -a) --output json | jq -s '.balances[]'
"$PWD"/build/mantrachaind q bank balances $(./build/mantrachaind keys show user1 -a) --output json | jq -s '.balances[]'
"$PWD"/build/mantrachaind q bank balances $(./build/mantrachaind keys show user2 -a) --output json | jq -s '.balances[]'
"$PWD"/build/mantrachaind q bank balances $(./build/mantrachaind keys show user3 -a) --output json | jq -s '.balances[]'