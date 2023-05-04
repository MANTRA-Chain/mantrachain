#!/bin/bash

PROV_CHAIN_ID=provider
PROV_KEY=provider-key
HOMEDIR=~/.ccv-provider
OUT_DIR=./out
CONSUMER_CHAIN_ID=mantrachain

source "$PWD"/scripts/common.sh

set -e # exit on first error

cecho "CYAN" "Create consumer proposal"
tee ./consumer-proposal.json <<EOF
{
  "title": "Create mantrachain",
  "description": "Mantrachain shared security",
  "chain_id": "mantrachain",
  "initial_height": {
      "revision_height": 1
  },
  "genesis_hash": "Z2VuX2hhc2g=",
  "binary_hash": "YmluX2hhc2g=",
  "spawn_time": "2022-01-27T15:59:50.121607-08:00",
  "blocks_per_distribution_transmission": 1000,
  "consumer_redistribution_fraction": "0.75",
  "historical_entries": 10000,
  "transfer_timeout_period": 3600000000000,
  "ccv_timeout_period": 2419200000000000,
  "unbonding_period": 1728000000000000,
  "deposit": "10000stake"
}
EOF

sleep 20

cecho "CYAN" "Submit consumer proposal"
interchain-security-pd tx gov submit-proposal \
  consumer-addition ./consumer-proposal.json \
  --keyring-backend test \
  --chain-id $PROV_CHAIN_ID \
  --from $PROV_KEY \
  --home $HOMEDIR \
  -b block -y

sleep 7

cecho "CYAN" "Make deposit to consumer proposal"
interchain-security-pd tx gov deposit 1 100000000stake \
  --keyring-backend test \
  --chain-id $PROV_CHAIN_ID \
  --from $PROV_KEY \
  --home $HOMEDIR \
  -b block -y

sleep 7

cecho "CYAN" "Vote on consumer proposal"
interchain-security-pd tx gov vote 1 yes --from $PROV_KEY \
  --keyring-backend test --chain-id $PROV_CHAIN_ID --home $HOMEDIR -b block -y

cecho "CYAN" "Wait 30 seconds for the consumer proposal to be processed..."
sleep 30

cecho "YELLOW" "Copy the consumer genesis to $HOMEDIR/consumer_genesis.json to continue [y]?: "

read -r confirm

mkdir -p $OUT_DIR

if [[ $confirm == "y" || $confirm == "Y" ]]; then
    sleep 7
    interchain-security-pd query provider consumer-genesis $CONSUMER_CHAIN_ID --home $HOMEDIR -o json > $HOMEDIR/ccvconsumer_genesis.json
    jq -s '.[0].app_state.ccvconsumer = .[1] | .[0]' $HOMEDIR/consumer_genesis.json $HOMEDIR/ccvconsumer_genesis.json > ${HOMEDIR}/edited_genesis.json
    mv ${HOMEDIR}/edited_genesis.json $OUT_DIR/consumer_genesis.json
    rm ${HOMEDIR}/consumer_genesis.json

    echo '{"height": "0","round": 0,"step": 0}' > $OUT_DIR/priv_validator_state.json
    cp ${HOMEDIR}/config/priv_validator_key.json $OUT_DIR/priv_validator_key.json
    cp ${HOMEDIR}/config/node_key.json $OUT_DIR/node_key.json

    COORDINATOR_P2P_ADDRESS=$(jq -r '.app_state.genutil.gen_txs[0].body.memo' $HOMEDIR/config/genesis.json)

    echo "$COORDINATOR_P2P_ADDRESS" > $OUT_DIR/cordinator_p2p_address.json
fi