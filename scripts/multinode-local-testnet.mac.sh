echo "install gsed, jq, tmux"
brew install gsed jq tmux

#!/bin/bash
rm -rf $HOME/.mantrachain/

echo "make four mantrachain directories"
mkdir $HOME/.mantrachain
mkdir $HOME/.mantrachain/validator1
mkdir $HOME/.mantrachain/validator2
mkdir $HOME/.mantrachain/validator3

echo "init all three validators"
mantrachaind init --chain-id=mantrachain validator1 --home=$HOME/.mantrachain/validator1
mantrachaind init --chain-id=mantrachain validator2 --home=$HOME/.mantrachain/validator2
mantrachaind init --chain-id=mantrachain validator3 --home=$HOME/.mantrachain/validator3

echo "create keys for all three validators"
mantrachaind keys add validator1 --keyring-backend=test --home=$HOME/.mantrachain/validator1
mantrachaind keys add validator2 --keyring-backend=test --home=$HOME/.mantrachain/validator2
mantrachaind keys add validator3 --keyring-backend=test --home=$HOME/.mantrachain/validator3

echo "change staking denom to uom"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="uom"' > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

echo "create validator node with tokens to transfer to the three other nodes"
mantrachaind add-genesis-account $(mantrachaind keys show validator1 -a --keyring-backend=test --home=$HOME/.mantrachain/validator1) 100000000000uom,100000000000stake --home=$HOME/.mantrachain/validator1
mantrachaind gentx validator1 500000000uom --keyring-backend=test --home=$HOME/.mantrachain/validator1 --chain-id=mantrachain

VALIDATOR_ADDRESS=$(cat $HOME/.mantrachain/config/gentx/$(ls $HOME/.mantrachain/config/gentx -AU | head -1) | jq '.body["messages"][0].validator_address')

echo "update vault genesis"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["vault"]["params"]["staking_validator_address"]='$VALIDATOR_ADDRESS > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

mantrachaind collect-gentxs --home=$HOME/.mantrachain/validator1


echo "update staking genesis"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["staking"]["params"]["unbonding_time"]="240s"' > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

echo "update crisis variable to uom"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="uom"' > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

echo "udpate gov genesis"
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json
cat $HOME/.mantrachain/validator1/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="uom"' > $HOME/.mantrachain/validator1/config/tmp_genesis.json && mv $HOME/.mantrachain/validator1/config/tmp_genesis.json $HOME/.mantrachain/validator1/config/genesis.json

echo "port key (validator1 uses default ports)"
echo "validator1 1317, 9090, 9091, 26658, 26657, 26656, 6060"
echo "validator2 1316, 9088, 9089, 26655, 26654, 26653, 6061"
echo "validator3 1315, 9086, 9087, 26652, 26651, 26650, 6062"
echo "validator4 1314, 9084, 9085, 26649, 26648, 26647, 6063"


echo "change app.toml values"

echo "validator2"
gsed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1316|g' $HOME/.mantrachain/validator2/config/app.toml
gsed -i -E 's|0.0.0.0:9090|0.0.0.0:9088|g' $HOME/.mantrachain/validator2/config/app.toml
gsed -i -E 's|0.0.0.0:9091|0.0.0.0:9089|g' $HOME/.mantrachain/validator2/config/app.toml

echo "validator3"
gsed -i -E 's|tcp://0.0.0.0:1317|tcp://0.0.0.0:1315|g' $HOME/.mantrachain/validator3/config/app.toml
gsed -i -E 's|0.0.0.0:9090|0.0.0.0:9086|g' $HOME/.mantrachain/validator3/config/app.toml
gsed -i -E 's|0.0.0.0:9091|0.0.0.0:9087|g' $HOME/.mantrachain/validator3/config/app.toml


echo "change config.toml values"

echo "validator1"
gsed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME/.mantrachain/validator1/config/config.toml
echo "validator2"
gsed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26655|g' $HOME/.mantrachain/validator2/config/config.toml
gsed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26654|g' $HOME/.mantrachain/validator2/config/config.toml
gsed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26653|g' $HOME/.mantrachain/validator2/config/config.toml
gsed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26650|g' $HOME/.mantrachain/validator3/config/config.toml
gsed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME/.mantrachain/validator2/config/config.toml
echo "validator3"
gsed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26652|g' $HOME/.mantrachain/validator3/config/config.toml
gsed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26651|g' $HOME/.mantrachain/validator3/config/config.toml
gsed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26650|g' $HOME/.mantrachain/validator3/config/config.toml
gsed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26650|g' $HOME/.mantrachain/validator3/config/config.toml
gsed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME/.mantrachain/validator3/config/config.toml


echo "copy validator1 genesis file to validator2-3"
cp $HOME/.mantrachain/validator1/config/genesis.json $HOME/.mantrachain/validator2/config/genesis.json
cp $HOME/.mantrachain/validator1/config/genesis.json $HOME/.mantrachain/validator3/config/genesis.json


echo "copy tendermint node id of validator1 to persistent peers of validator2-3"
gsed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$(mantrachaind tendermint show-node-id --home=$HOME/.mantrachain/validator1)@$(curl -4 icanhazip.com):26656\"|g" $HOME/.mantrachain/validator2/config/config.toml
gsed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$(mantrachaind tendermint show-node-id --home=$HOME/.mantrachain/validator1)@$(curl -4 icanhazip.com):26656\"|g" $HOME/.mantrachain/validator3/config/config.toml


echo "start all three validators"
tmux new -s validator1 -d mantrachaind start --home=$HOME/.mantrachain/validator1
tmux new -s validator2 -d mantrachaind start --home=$HOME/.mantrachain/validator2
tmux new -s validator3 -d mantrachaind start --home=$HOME/.mantrachain/validator3


echo "send uom from first validator to second validator"
sleep 7
mantrachaind tx bank send validator1 $(mantrachaind keys show validator2 -a --keyring-backend=test --home=$HOME/.mantrachain/validator2) 500000000uom --keyring-backend=test --home=$HOME/.mantrachain/validator1 --chain-id=mantrachain --yes
sleep 7
mantrachaind tx bank send validator1 $(mantrachaind keys show validator3 -a --keyring-backend=test --home=$HOME/.mantrachain/validator3) 400000000uom --keyring-backend=test --home=$HOME/.mantrachain/validator1 --chain-id=mantrachain --yes

echo "create second validator"
sleep 7
mantrachaind tx staking create-validator --amount=500000000uom --from=validator2 --pubkey=$(mantrachaind tendermint show-validator --home=$HOME/.mantrachain/validator2) --moniker="validator2" --chain-id="mantrachain" --commission-rate="0.1" --commission-max-rate="0.2" --commission-max-change-rate="0.05" --min-self-delegation="500000000" --keyring-backend=test --home=$HOME/.mantrachain/validator2 --yes
sleep 7
mantrachaind tx staking create-validator --amount=400000000uom --from=validator3 --pubkey=$(mantrachaind tendermint show-validator --home=$HOME/.mantrachain/validator3) --moniker="validator3" --chain-id="mantrachain" --commission-rate="0.1" --commission-max-rate="0.2" --commission-max-change-rate="0.05" --min-self-delegation="400000000" --keyring-backend=test --home=$HOME/.mantrachain/validator3 --yes
