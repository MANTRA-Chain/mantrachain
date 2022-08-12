echo "not working yet"
exit 1

echo "install gsed, jq, tmux"
brew install gsed jq tmux

echo "stop all the validators if any"
tmux kill-session -t validator

kill -9 $(lsof -t -i:26656)
kill -9 $(lsof -t -i:6060)
kill -9 $(lsof -t -i:1317)

#!/bin/bash
rm -rf $HOME/.mantrachain/

echo "make four mantrachain directories"
mkdir $HOME/.mantrachain
mkdir $HOME/.mantrachain

echo "init the validator"
mantrachaind init --chain-id=mantrachain validator --home=$HOME/.mantrachain

echo "create keys for the validator"
mantrachaind keys add validator --keyring-backend=test --home=$HOME/.mantrachain
mantrachaind keys add recipient --keyring-backend=test --home=$HOME/.mantrachain

VALIDATOR=$(./build/mantrachaind keys show validator -a --keyring-backend test)
RECIPIENT=$(./build/mantrachaind keys show recipient -a --keyring-backend test)

echo "create validator node with tokens to transfer to the three other nodes"
mantrachaind add-genesis-account $(mantrachaind keys show validator -a --keyring-backend=test --home=$HOME/.mantrachain) 100000000000stake --home=$HOME/.mantrachain
mantrachaind gentx validator 500000000stake --keyring-backend=test --home=$HOME/.mantrachain --chain-id=mantrachain

VALIDATOR_ADDRESS=$(cat $HOME/.mantrachain/config/gentx/$(ls $HOME/.mantrachain/config/gentx -AU | head -1) | jq '.body["messages"][0].validator_address')

echo "update vault genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["vault"]["params"]["staking_validator_address"]='$VALIDATOR_ADDRESS > $HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

mantrachaind collect-gentxs --home=$HOME/.mantrachain


echo "update staking genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["staking"]["params"]["unbonding_time"]="240s"' > $HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

echo "udpate gov genesis"
cat $HOME/.mantrachain/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' > $HOME/.mantrachain/config/tmp_genesis.json && mv $HOME/.mantrachain/config/tmp_genesis.json $HOME/.mantrachain/config/genesis.json

echo "port key (validator uses default ports)"
echo "validator 1317, 9090, 9091, 26658, 26657, 26656, 6060"


echo "change app.toml values"

echo "validator"
gsed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $HOME/.mantrachain/config/app.toml
gsed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' $HOME/.mantrachain/config/app.toml
gsed -i -E 's|minimum-gas-prices = \"\"|minimum-gas-prices = \"0.00001stake\"|g' $HOME/.mantrachain/config/app.toml

echo "change config.toml values"

echo "validator"
gsed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $HOME/.mantrachain/config/config.toml

echo "change client.toml values"

echo "validator"
gsed -i -E 's|chain-id = \"\"|chain-id = \"mantrachain\"|g' $HOME/.mantrachain/config/client.toml

echo "start the validator"
# tmux new -s validator -d mantrachaind start --home=$HOME/.mantrachain
# tmux a -t validator