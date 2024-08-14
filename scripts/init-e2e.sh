#!/bin/bash

CHAINID="mantrachain-devnet-9001"
HOMEDIR="$HOME/.mantrachain"
KEYRING="test"
GAS_ADJ=2
GAS_PRICE=0.0002uom
ADMIN_WALLET=$(./build/mantrachaind keys show admin -a --keyring-backend $KEYRING --home "$HOMEDIR")

source "$PWD"/scripts/common.sh

set -e # exit on first error

CW_PROJECT_NAME=e2e_tests
CODE_ID=1

"$PWD"/build/mantrachaind tx wasm store "$PWD"/tests/e2e/src/wasm/$CW_PROJECT_NAME.wasm --from admin --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes

INIT={\"count\":0}

sleep 3

"$PWD"/build/mantrachaind tx wasm instantiate $CODE_ID "$INIT" --label "e2e-test" --admin admin --from admin --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes
