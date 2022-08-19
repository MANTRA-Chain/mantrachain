#!/bin/bash
set -e # exit on first error

rm -rf $HOME/.mantrachain/

echo "create keys for the validator"
./build/mantrachaind keys add validator --keyring-backend=test
./build/mantrachaind keys add recipient --keyring-backend=test

VALIDATOR=$(./build/mantrachaind keys show validator -a --keyring-backend test)
RECIPIENT=$(./build/mantrachaind keys show recipient -a --keyring-backend test)

echo "init the validator"
./build/mantrachaind init mantrachain --chain-id=mantrachain

echo "create validator node with tokens to transfer to the node"
./build/mantrachaind add-genesis-account $VALIDATOR 100000000000stake
./build/mantrachaind gentx validator 100000000stake --keyring-backend=test --chain-id=mantrachain

VALIDATOR_ADDRESS=$(cat $HOME/.mantrachain/config/gentx/$(ls $HOME/.mantrachain/config/gentx | head -1) | jq '.body["messages"][0].validator_address')

echo "update vault genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["vault"]["params"]["staking_validator_address"]='$VALIDATOR_ADDRESS >$HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

./build/mantrachaind collect-gentxs

echo "update staking genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["staking"]["params"]["unbonding_time"]="240s"' >$HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

echo "udpate gov genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' >$HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

echo "port key (validator uses default ports)"
echo "validator 1317, 9090, 9091, 26658, 26657, 26656, 6060"

echo "change app.toml values"

echo "validator"
sed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $HOME/.mantrachain/config/app.toml
sed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' $HOME/.mantrachain/config/app.toml
sed -i -E 's|minimum-gas-prices = \"\"|minimum-gas-prices = \"0.00001stake\"|g' $HOME/.mantrachain/config/app.toml

echo "change config.toml values"

echo "validator"
sed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $HOME/.mantrachain/config/config.toml

echo "change client.toml values"

echo "validator"
sed -i -E 's|chain-id = \"\"|chain-id = \"mantrachain\"|g' $HOME/.mantrachain/config/client.toml

echo "start the validator"
tmux new -s validator -d ./build/mantrachaind start

echo "send stake from validator to recipient"
sleep 7
./build/mantrachaind tx bank send $VALIDATOR $RECIPIENT 10000000000stake --chain-id mantrachain --keyring-backend test --gas auto --gas-adjustment 1.25 --gas-prices 0.00001stake --yes

echo "track logs"
tmux a -t validator
