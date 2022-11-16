#!/bin/bash

echo "stop the validators if any"
tmux kill-session -t validator

kill -9 $(lsof -t -i:26656)

set -e # exit on first error

rm -rf $HOME/.mantrachain/

echo "create keys for the validator"
{
./build/mantrachaind keys add validator --keyring-backend=test
./build/mantrachaind keys add recipient --keyring-backend=test
./build/mantrachaind keys add admin --keyring-backend=test
} 2>&1 | tee accounts.txt

VALIDATOR=$(./build/mantrachaind keys show validator -a --keyring-backend test)
RECIPIENT=$(./build/mantrachaind keys show recipient -a --keyring-backend test)
ADMIN=$(./build/mantrachaind keys show admin -a --keyring-backend test)

echo "init the validator"
./build/mantrachaind init mantrachain --chain-id=mantrachain

echo "change staking denom to ustake"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="ustake"' > $HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

echo "create validator node with tokens to transfer to the node"
./build/mantrachaind add-genesis-account $VALIDATOR 100000000000000000ustake
./build/mantrachaind gentx validator 100000000000000ustake --keyring-backend=test --chain-id=mantrachain

VALIDATOR_ADDRESS=$(cat $HOME/.mantrachain/config/gentx/$(ls $HOME/.mantrachain/config/gentx | head -1) | jq '.body["messages"][0].validator_address')

echo "update vault genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["vault"]["params"]["staking_validator_address"]='$VALIDATOR_ADDRESS >$HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["vault"]["params"]["admin_account"]='\"$ADMIN\" >$HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

echo "update bridge genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["bridge"]["params"]["admin_account"]='\"$ADMIN\" >$HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

echo "update denom metadata"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["bank"]["denom_metadata"]=''[{"description":"The native staking token of the Mantrachain.","denom_units":[{"denom":"ustake","exponent":0,"aliases":["microstake"]},{"denom":"mstake","exponent":3,"aliases":["millistake"]},{"denom":"stake","exponent":6}],"base":"ustake","display":"stake"}]' > $HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

echo "update mint genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="ustake"' >$HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

echo "update crisis variable to ustake"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="ustake"' > $HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

./build/mantrachaind collect-gentxs

echo "update staking genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["staking"]["params"]["unbonding_time"]="240s"' >$HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

echo "udpate gov genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' >$HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="ustake"' >$HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json


echo "port key (validator uses default ports)"
echo "validator 1317, 9090, 9091, 26658, 26657, 26656, 6060"

echo "change app.toml values"

echo "validator"
sed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $HOME/.mantrachain/config/app.toml
sed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' $HOME/.mantrachain/config/app.toml
sed -i -E '1,/enable = false/s|enable = false|enable = true|g' $HOME/.mantrachain/config/app.toml
sed -i -E 's|minimum-gas-prices = \"\"|minimum-gas-prices = \"0.0001ustake\"|g' $HOME/.mantrachain/config/app.toml

echo "change config.toml values"

echo "validator"
sed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $HOME/.mantrachain/config/config.toml

echo "change client.toml values"

echo "validator"
sed -i -E 's|chain-id = \"\"|chain-id = \"mantrachain\"|g' $HOME/.mantrachain/config/client.toml

echo "start the validator"
tmux new -s validator -d ./build/mantrachaind start

echo "send ustake from validator to recipient"
sleep 7
./build/mantrachaind tx bank send $VALIDATOR $RECIPIENT 100000000000000ustake --chain-id mantrachain --keyring-backend test --gas auto --gas-adjustment 1.3 --gas-prices 0.0001ustake --yes

echo "track logs"
tmux a -t validator
