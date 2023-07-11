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

DENOM_A="a"
DENOM_B="b"

source "$PWD"/scripts/common.sh

set -e # exit on first error

VALIDATOR_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[0]} -a --keyring-backend $KEYRING --home "$HOMEDIR")
RECIPIENT_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[1]} -a --keyring-backend $KEYRING --home "$HOMEDIR")
ADMIN_WALLET=$("$PWD"/build/mantrachaind keys show ${KEYS[2]} -a --keyring-backend $KEYRING --home "$HOMEDIR")

cecho "CYAN" "Send aum from ${KEYS[0]} to ${KEYS[1]}"
"$PWD"/build/mantrachaind tx bank send $VALIDATOR_WALLET $RECIPIENT_WALLET 100000000000000uaum --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" -y

sleep 7

cecho "CYAN" "Send uaum from ${KEYS[0]} to ${KEYS[2]}"
"$PWD"/build/mantrachaind tx bank send $VALIDATOR_WALLET $ADMIN_WALLET 100000000000000uaum --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" -y

sleep 7

cecho "CYAN" "Create denom $DENOM_A"
"$PWD"/build/mantrachaind tx coinfactory create-denom $DENOM_A --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --chain-id $CHAINID --home "$HOMEDIR" -y

sleep 7

cecho "CYAN" "Create denom $DENOM_B"
"$PWD"/build/mantrachaind tx coinfactory create-denom $DENOM_B --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --chain-id $CHAINID --home "$HOMEDIR" -y

DENOM_A_FULL="factory/$ADMIN_WALLET/$DENOM_A"
DENOM_B_FULL="factory/$ADMIN_WALLET/$DENOM_B"

sleep 7

cecho "CYAN" "Mint 100000000000000000 $DENOM_A"
"$PWD"/build/mantrachaind tx coinfactory mint 100000000000000000$DENOM_A_FULL --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --chain-id $CHAINID --home "$HOMEDIR" -y

sleep 7

cecho "CYAN" "Mint 100000000000000000 $DENOM_B"
"$PWD"/build/mantrachaind tx coinfactory mint 100000000000000000$DENOM_B_FULL --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --chain-id $CHAINID --home "$HOMEDIR" -y

sleep 7

cecho "CYAN" "Create pair $DENOM_A/$DENOM_B"
"$PWD"/build/mantrachaind tx liquidity create-pair $DENOM_A_FULL $DENOM_B_FULL --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --chain-id $CHAINID --home "$HOMEDIR" -y

sleep 7

PAIR_ID=$("$PWD"/build/mantrachaind q liquidity pairs --output json | jq -r '.pairs[-1].id')

cecho "CYAN" "Create Basic Pool with 1000000000000$DENOM_A_FULL 1000000000000$DENOM_B_FULL"
"$PWD"/build/mantrachaind tx liquidity create-pool $(($PAIR_ID)) "1000000000000$DENOM_A_FULL,1000000000000$DENOM_B_FULL" --from ${KEYS[2]} --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --chain-id $CHAINID --home "$HOMEDIR" -y

ADMIN_BALANCES=$("$PWD"/build/mantrachaind q bank balances $("$PWD"/build/mantrachaind keys show admin -a))

cecho "YELLOW" "$ADMIN_BALANCE"