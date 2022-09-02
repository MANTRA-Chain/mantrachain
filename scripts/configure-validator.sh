#!/bin/bash
set -e # exit on first error

GENESIS_FILE="./data/genesis.json"
PERSISTENT_PEERS="nodeId@ip:port,nodeId@ip:port"

./build/mantrachaind keys add validator-1 --recover

#ACCOUNT=$(./build/mantrachaind keys show validator-1 -a)

cp $GENESIS_FILE "$HOME"/.mantrachain/config

# Configure app.toml
######################################################################
sed -i -E 's|minimum-gas-prices = \"\"|minimum-gas-prices = \"0.00001stake\"|g' "$HOME"/.mantrachain/config/app.toml
sed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' "$HOME"/.mantrachain/config/app.toml
#sed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' "$HOME"/.mantrachain/config/app.toml
#sed -i -E '108s/.*/enable = true/' "$HOME"/.mantrachain/config/app.toml
#sed -i -E '111s/.*/swagger = true/' "$HOME"/.mantrachain/config/app.toml

# Configure config.toml
######################################################################
sed -i -E 's|pex = \[\]|pex = \"false\"|g' "$HOME"/.mantrachain/config/config.toml
sed -i -E 's|persistent_peers = \[\]|persistent_peers = \"'$PERSISTENT_PEERS'\"|g' "$HOME"/.mantrachain/config/config.toml
sed -i -E 's|addr_book_strict = \[\]|addr_book_strict = \"false\"|g' "$HOME"/.mantrachain/config/config.toml
sed -i -E 's|double_sign_check_height = \[\]|double_sign_check_height = \"10\"|g' "$HOME"/.mantrachain/config/config.toml

# Configure client.toml
######################################################################
sed -i -E 's|chain-id = \"\"|chain-id = \"'"$CHAIN_ID"'"|g' "$HOME"/.mantrachain/config/client.toml

printf "Done! Start the node with mantrachaind start\n"