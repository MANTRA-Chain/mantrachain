#!/bin/bash

HOMEDIR=~/.ccv-provider
MONIKER=ccv-provider
CHAINID=provider
KEY=provider-key

source "$PWD"/scripts/common.sh

echo "Stop the provider if any"
pkill -f interchain-security-pd

set -e # exit on first error

# User prompt if an existing local node configuration is found.
if [ -d "$HOMEDIR" ]; then
  printf "\nAn ccv provider existing folder at '%s' was found. You can choose to delete this folder and start a new local provider node with new keys from genesis. When declined, the existing local node is started. \n" "$HOMEDIR"
  cecho "RED" "CCV PROVIDER: Overwrite the existing configuration and start a new local provider node? [y/N]"
  read -r overwrite
else
  overwrite="Y"
fi

# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
  rm -rf "$HOME/.ccv-provider"

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
  mkdir $HOMEDIR

  cecho "CYAN" "Init the provider"
  interchain-security-pd init $MONIKER -o --chain-id $CHAINID --home $HOMEDIR

  cecho "CYAN" "Create key"
  interchain-security-pd keys add $KEY --home $HOMEDIR --keyring-backend test --output json >${HOMEDIR}/${KEY}.json 2>&1

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

  cecho "CYAN" "Change client.toml values"
  sed -i -r "/node =/ s/= .*/= \"tcp:\/\/127.0.0.1:26658\"/" ${HOMEDIR}/config/client.toml
fi

cecho "CYAN" "Start the provider"
tmux new -s provider -d interchain-security-pd start --home $HOMEDIR \
  --rpc.laddr tcp://127.0.0.1:26658 \
  --grpc.address 127.0.0.1:9091 \
  --address tcp://127.0.0.1:26655 \
  --p2p.laddr tcp://127.0.0.1:26656 \
  --grpc-web.enable=false \
  &>${HOMEDIR}/logs &
