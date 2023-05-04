#!/bin/bash

HOMEDIR=~/.ccv-provider
MONIKER=ccv-provider
CHAINID=provider
KEY=provider-key

source "$PWD"/scripts/common.sh



echo "Stop the provider if any"
pkill -f interchain-security-pd
rm -rf ./interchain-security

set -e # exit on first error

git clone \
  --filter=blob:none \
  --no-checkout \
  https://github.com/cosmos/interchain-security.git

cd interchain-security
git checkout tags/v1.1.0 -b v1.1.0

make install
cd ..

rm -rf ./interchain-security

rm -rf $HOMEDIR
mkdir -p $HOMEDIR

cecho "CYAN" "Init the provider"
interchain-security-pd init $MONIKER -o --chain-id $CHAINID --home $HOMEDIR

cecho "CYAN" "Create key"
interchain-security-pd keys add $KEY --home $HOMEDIR --keyring-backend test --home $HOMEDIR --output json > ${HOMEDIR}/${KEY}.json 2>&1

cecho "CYAN" "Update genesis"
jq ".app_state.gov.voting_params.voting_period = \"30s\"" ${HOMEDIR}/config/genesis.json >${HOMEDIR}/edited_genesis.json
mv ${HOMEDIR}/edited_genesis.json ${HOMEDIR}/config/genesis.json

PROV_ACCOUNT_ADDR=$(jq -r .address ${HOMEDIR}/${KEY}.json)

cecho "CYAN" "Create provider node with tokens to transfer to the node"
interchain-security-pd add-genesis-account $PROV_ACCOUNT_ADDR 1000000000000000stake --keyring-backend test --home $HOMEDIR

interchain-security-pd gentx $KEY 100000000000stake \
  --keyring-backend test \
  --moniker $MONIKER \
  --chain-id $CHAINID \
  --home $HOMEDIR

cecho "CYAN" "Collect genesis tx"
interchain-security-pd collect-gentxs --home $HOMEDIR --gentx-dir ${HOMEDIR}/config/gentx/

cecho "CYAN" "You can start the chain with:"
cecho "YELLOW" "interchain-security-pd start --home $HOMEDIR --grpc.address 127.0.0.1:9092 --grpc-web.enable=false"