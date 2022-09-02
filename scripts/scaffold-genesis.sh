#!/bin/bash
set -e # exit on first error

CHAIN_DATA_PATH="./data"
KEYRING_BACKEND="test"
CHAIN_ID="mantrachain-1"
VALIDATOR1_NAME="validator-1"
VALIDATOR2_NAME="validator-2"

# clear previous data
rm -rf $CHAIN_DATA_PATH

mkdir $CHAIN_DATA_PATH
mkdir $CHAIN_DATA_PATH/validator-1
mkdir $CHAIN_DATA_PATH/validator-2


# Generate accounts and node ids for our validators
# validator-1
##################################################################
{
  ./build/mantrachaind keys add "$VALIDATOR1_NAME" --keyring-backend=$KEYRING_BACKEND --home=./$CHAIN_DATA_PATH/validator-1
} 2>&1 | tee $CHAIN_DATA_PATH/validator-1/validator-1.txt

VALIDATOR1_ADDRESS=$(./build/mantrachaind keys show "$VALIDATOR1_NAME" -a --keyring-backend=$KEYRING_BACKEND --home=./$CHAIN_DATA_PATH/validator-1)


./build/mantrachaind init "$VALIDATOR1_NAME" --chain-id=$CHAIN_ID --home=./$CHAIN_DATA_PATH/validator-1

{
  ./build/mantrachaind tendermint show-node-id --home=./$CHAIN_DATA_PATH/validator-1
} 2>&1 | tee $CHAIN_DATA_PATH/validator-1/validator-node-1-id.txt

./build/mantrachaind add-genesis-account "$VALIDATOR1_ADDRESS" 100000000000stake --home=./$CHAIN_DATA_PATH/validator-1

# validator-2
##################################################################
{
  ./build/mantrachaind keys add "$VALIDATOR2_NAME" --keyring-backend=$KEYRING_BACKEND --home=./$CHAIN_DATA_PATH/validator-2
} 2>&1 | tee $CHAIN_DATA_PATH/validator-2/validator-2.txt

VALIDATOR2_ADDRESS=$(./build/mantrachaind keys show $VALIDATOR2_NAME -a --keyring-backend=$KEYRING_BACKEND --home=./$CHAIN_DATA_PATH/validator-2)

./build/mantrachaind init $VALIDATOR2_NAME --chain-id=$CHAIN_ID --home=./$CHAIN_DATA_PATH/validator-2

{
  ./build/mantrachaind tendermint show-node-id --home=./$CHAIN_DATA_PATH/validator-2
} 2>&1 | tee $CHAIN_DATA_PATH/validator-2/validator-node-2-id.txt


# Add params to genesis
##################################################################
cat $CHAIN_DATA_PATH/validator-1/config/genesis.json | jq '.app_state["vault"]["params"]["staking_validator_address"]="'"$VALIDATOR1_ADDRESS"'"' > $CHAIN_DATA_PATH/validator-1/config/tmp_genesis.json && mv $CHAIN_DATA_PATH/validator-1/config/tmp_genesis.json $CHAIN_DATA_PATH/validator-1/config/genesis.json

echo "update staking genesis"
cat $CHAIN_DATA_PATH/validator-1/config/genesis.json | jq '.app_state["staking"]["params"]["unbonding_time"]="240s"' > $CHAIN_DATA_PATH/validator-1/config/tmp_genesis.json && mv $CHAIN_DATA_PATH/validator-1/config/tmp_genesis.json $CHAIN_DATA_PATH/validator-1/config/genesis.json

echo "udpate gov genesis"
cat $CHAIN_DATA_PATH/validator-1/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' > $CHAIN_DATA_PATH/validator-1/config/tmp_genesis.json && mv $CHAIN_DATA_PATH/validator-1/config/tmp_genesis.json $CHAIN_DATA_PATH/validator-1/config/genesis.json

#Add validator-2 as genesis account
##################################################################
./build/mantrachaind add-genesis-account "$VALIDATOR2_ADDRESS" 100000000000stake --home=./$CHAIN_DATA_PATH/validator-1

#Send the genesis to validator-2 for review
##################################################################
cp $CHAIN_DATA_PATH/validator-1/config/genesis.json $CHAIN_DATA_PATH/validator-2/config/genesis.json

./build/mantrachaind gentx $VALIDATOR2_NAME 100000000stake --keyring-backend=$KEYRING_BACKEND --chain-id=$CHAIN_ID --home=./$CHAIN_DATA_PATH/validator-2


#Recieve back its genesis transactions
##################################################################
mkdir $CHAIN_DATA_PATH/validator-1/config/gentx

cp $CHAIN_DATA_PATH/validator-2/config/gentx/gentx-* $CHAIN_DATA_PATH/validator-1/config/gentx

#Add validator-1 genesis transaction
##################################################################
./build/mantrachaind gentx $VALIDATOR1_NAME 100000000stake --keyring-backend=$KEYRING_BACKEND --chain-id=$CHAIN_ID --home=./$CHAIN_DATA_PATH/validator-1

#Execute the transactions against the genesis
##################################################################
./build/mantrachaind collect-gentxs --home=./$CHAIN_DATA_PATH/validator-1


cp $CHAIN_DATA_PATH/validator-1/config/genesis.json $CHAIN_DATA_PATH/genesis.json

#Send back the genesis
##################################################################
cp $CHAIN_DATA_PATH/validator-1/config/genesis.json $CHAIN_DATA_PATH/validator-2/config/genesis.json