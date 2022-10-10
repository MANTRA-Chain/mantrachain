#!/bin/bash

echo "stop all the validators if any"
tmux kill-session -t validator1
tmux kill-session -t validator2
tmux kill-session -t validator3

kill -9 $(lsof -t -i:26650) $(lsof -t -i:26653) $(lsof -t -i:26656)

set -e # exit on first error

rm -rf $HOME/.mantrachain/

echo "make four mantrachain directories"
mkdir $HOME/.mantrachain
mkdir $HOME/.mantrachain/validator1
mkdir $HOME/.mantrachain/validator2
mkdir $HOME/.mantrachain/validator3

echo "create keys for all three validators"
{
./build/mantrachaind keys add validator1 --keyring-backend=test --home=$HOME/.mantrachain/validator1
./build/mantrachaind keys add recipient1 --keyring-backend=test --home=$HOME/.mantrachain/validator1
./build/mantrachaind keys add admin1 --keyring-backend=test --home=$HOME/.mantrachain/validator1
} 2>&1 | tee accounts1.txt
{
./build/mantrachaind keys add validator2 --keyring-backend=test --home=$HOME/.mantrachain/validator2
} 2>&1 | tee accounts2.txt
{
./build/mantrachaind keys add validator3 --keyring-backend=test --home=$HOME/.mantrachain/validator3
} 2>&1 | tee accounts3.txt

echo "init all three validators"
./build/mantrachaind init --chain-id=mantrachain validator1 --home=$HOME/.mantrachain/validator1
./build/mantrachaind init --chain-id=mantrachain validator2 --home=$HOME/.mantrachain/validator2
./build/mantrachaind init --chain-id=mantrachain validator3 --home=$HOME/.mantrachain/validator3

echo "change staking denom to ustake"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="ustake"' > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

echo "create validator node with tokens to transfer to the three other nodes"
./build/mantrachaind add-genesis-account $(./build/mantrachaind keys show validator1 -a --keyring-backend=test --home=$HOME/.mantrachain/validator1) 100000000000000000ustake --home=$HOME/.mantrachain/validator1
./build/mantrachaind gentx validator1 100000000000000ustake --keyring-backend=test --home=$HOME/.mantrachain/validator1 --chain-id=mantrachain

VALIDATOR1_ADDRESS=$(cat $HOME/.mantrachain/validator1/config/gentx/$(ls $HOME/.mantrachain/validator1/config/gentx -AU | head -1) | jq '.body["messages"][0].validator_address')

ADMIN1=$(./build/mantrachaind keys show admin1 -a --keyring-backend=test --home=$HOME/.mantrachain/validator1)

echo "update vault genesis"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["vault"]["params"]["staking_validator_address"]='$VALIDATOR1_ADDRESS > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

echo "update bridge genesis"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["bridge"]["params"]["admin_account"]='\"$ADMIN1\" >$HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

echo "update denom metadata"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["bank"]["denom_metadata"]=''[{"description":"The native staking token of the Mantrachain.","denom_units":[{"denom":"ustake","exponent":0,"aliases":["microstake"]},{"denom":"mstake","exponent":3,"aliases":["millistake"]},{"denom":"stake","exponent":6}],"base":"ustake","display":"stake"}]' > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

echo "update crisis variable to ustake"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="ustake"' > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

./build/mantrachaind collect-gentxs --home=$HOME/.mantrachain/validator1

echo "update staking genesis"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["staking"]["params"]["unbonding_time"]="240s"' > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

echo "udpate gov genesis"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

echo "port key (validator1 uses default ports)"
echo "validator1 1317, 9090, 9091, 26658, 26657, 26656, 6060"
echo "validator2 1316, 9088, 9089, 26655, 26654, 26653, 6061"
echo "validator3 1315, 9086, 9087, 26652, 26651, 26650, 6062"

echo "change app.toml values"

echo "validator1"
gsed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $HOME/.mantrachain/validator1/config/app.toml
gsed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' $HOME/.mantrachain/validator1/config/app.toml
gsed -i -E '1,/enable = false/s|enable = false|enable = true|g' $HOME/.mantrachain/validator1/config/app.toml
gsed -i -E 's|minimum-gas-prices = \"\"|minimum-gas-prices = \"0.0001ustake\"|g' $HOME/.mantrachain/validator1/config/app.toml

echo "validator2"
gsed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $HOME/.mantrachain/validator2/config/app.toml
gsed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' $HOME/.mantrachain/validator2/config/app.toml
gsed -i -E '1,/enable = false/s|enable = false|enable = true|g' $HOME/.mantrachain/validator2/config/app.toml
gsed -i -E 's|minimum-gas-prices = \"\"|minimum-gas-prices = \"0.0001ustake\"|g' $HOME/.mantrachain/validator2/config/app.toml
gsed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1316|g' $HOME/.mantrachain/validator2/config/app.toml
gsed -i -E 's|0.0.0.0:9090|0.0.0.0:9088|g' $HOME/.mantrachain/validator2/config/app.toml
gsed -i -E 's|0.0.0.0:9091|0.0.0.0:9089|g' $HOME/.mantrachain/validator2/config/app.toml

echo "validator3"
gsed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $HOME/.mantrachain/validator3/config/app.toml
gsed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' $HOME/.mantrachain/validator3/config/app.toml
gsed -i -E '1,/enable = false/s|enable = false|enable = true|g' $HOME/.mantrachain/validator3/config/app.toml
gsed -i -E 's|minimum-gas-prices = \"\"|minimum-gas-prices = \"0.0001ustake\"|g' $HOME/.mantrachain/validator3/config/app.toml
gsed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1315|g' $HOME/.mantrachain/validator3/config/app.toml
gsed -i -E 's|0.0.0.0:9090|0.0.0.0:9086|g' $HOME/.mantrachain/validator3/config/app.toml
gsed -i -E 's|0.0.0.0:9091|0.0.0.0:9087|g' $HOME/.mantrachain/validator3/config/app.toml

echo "change config.toml values"

echo "validator1"
gsed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME/.mantrachain/validator1/config/config.toml
gsed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $HOME/.mantrachain/validator1/config/config.toml
echo "validator2"
gsed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $HOME/.mantrachain/validator2/config/config.toml
gsed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26655|g' $HOME/.mantrachain/validator2/config/config.toml
gsed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26654|g' $HOME/.mantrachain/validator2/config/config.toml
gsed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26653|g' $HOME/.mantrachain/validator2/config/config.toml
gsed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME/.mantrachain/validator2/config/config.toml
echo "validator3"
gsed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' $HOME/.mantrachain/validator3/config/config.toml
gsed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26652|g' $HOME/.mantrachain/validator3/config/config.toml
gsed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26651|g' $HOME/.mantrachain/validator3/config/config.toml
gsed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26650|g' $HOME/.mantrachain/validator3/config/config.toml
gsed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME/.mantrachain/validator3/config/config.toml

echo "change client.toml values"

echo "validator1"
gsed -i -E 's|chain-id = \"\"|chain-id = \"mantrachain\"|g' $HOME/.mantrachain/validator1/config/client.toml
echo "validator2"
gsed -i -E 's|chain-id = \"\"|chain-id = \"mantrachain\"|g' $HOME/.mantrachain/validator2/config/client.toml
gsed -i -E 's|26657|26654|g' $HOME/.mantrachain/validator2/config/client.toml
echo "validator3"
gsed -i -E 's|chain-id = \"\"|chain-id = \"mantrachain\"|g' $HOME/.mantrachain/validator3/config/client.toml
gsed -i -E 's|26657|26651|g' $HOME/.mantrachain/validator3/config/client.toml

echo "copy validator1 genesis file to validator2-3"
cp $HOME/.mantrachain/validator1/config/genesis.json $HOME/.mantrachain/validator2/config/genesis.json
cp $HOME/.mantrachain/validator1/config/genesis.json $HOME/.mantrachain/validator3/config/genesis.json

echo "copy tendermint node id of validator1 to persistent peers of validator2-3"
gsed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$(./build/mantrachaind tendermint show-node-id --home=$HOME/.mantrachain/validator1)@localhost:26656\"|g" $HOME/.mantrachain/validator2/config/config.toml
gsed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$(./build/mantrachaind tendermint show-node-id --home=$HOME/.mantrachain/validator1)@localhost:26656\"|g" $HOME/.mantrachain/validator3/config/config.toml

echo "start all three validators"
tmux new -s validator1 -d ./build/mantrachaind start --home=$HOME/.mantrachain/validator1
tmux new -s validator2 -d ./build/mantrachaind start --home=$HOME/.mantrachain/validator2
tmux new -s validator3 -d ./build/mantrachaind start --home=$HOME/.mantrachain/validator3

echo "send stake from first validator to second validator"
sleep 7
./build/mantrachaind tx bank send validator1 $(./build/mantrachaind keys show validator2 -a --keyring-backend=test --home=$HOME/.mantrachain/validator2) 500000000000ustake --keyring-backend=test --home=$HOME/.mantrachain/validator1 --chain-id=mantrachain --gas=auto --gas-adjustment="1.3" --gas-prices="0.0001ustake" --yes
echo "send stake from first validator to third validator"
sleep 7
./build/mantrachaind tx bank send validator1 $(./build/mantrachaind keys show validator3 -a --keyring-backend=test --home=$HOME/.mantrachain/validator3) 400000000000ustake --keyring-backend=test --home=$HOME/.mantrachain/validator1 --chain-id=mantrachain --gas=auto --gas-adjustment="1.3" --gas-prices="0.0001ustake" --yes

echo "create second validator"
sleep 7
./build/mantrachaind tx staking create-validator --amount=500000000000ustake --from=validator2 --pubkey=$(./build/mantrachaind tendermint show-validator --home=$HOME/.mantrachain/validator2) --moniker="validator2" --chain-id="mantrachain" --commission-rate="0.1" --commission-max-rate="0.2" --commission-max-change-rate="0.05" --min-self-delegation="500000000" --keyring-backend=test --home=$HOME/.mantrachain/validator2 --gas=auto --gas-adjustment="1.3" --gas-prices="0.0001ustake" --yes
sleep 7
./build/mantrachaind tx staking create-validator --amount=400000000000ustake --from=validator3 --pubkey=$(./build/mantrachaind tendermint show-validator --home=$HOME/.mantrachain/validator3) --moniker="validator3" --chain-id="mantrachain" --commission-rate="0.1" --commission-max-rate="0.2" --commission-max-change-rate="0.05" --min-self-delegation="400000000" --keyring-backend=test --home=$HOME/.mantrachain/validator3 --gas=auto --gas-adjustment="1.3" --gas-prices="0.0001ustake" --yes

echo "send ustake from validator1 to recipient1"
sleep 7
./build/mantrachaind tx bank send validator1 $(./build/mantrachaind keys show recipient1 -a --keyring-backend=test --home=$HOME/.mantrachain/validator1) 100000000000000ustake --keyring-backend=test --home=$HOME/.mantrachain/validator1 --chain-id=mantrachain --gas=auto --gas-adjustment="1.3" --gas-prices="0.0001ustake" --yes