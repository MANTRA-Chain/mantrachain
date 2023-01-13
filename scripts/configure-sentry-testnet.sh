#!/bin/bash
set -e # exit on first error

GENESIS_FILE="./data/genesis.json"
PERSISTENT_PEERS="nodeId@ip:port,nodeId@ip:port"
PRIVATE_PEERS="nodeId,nodeId"
UNCONDITIONAL_PEER_IDS="nodeId,nodeId"

cp $GENESIS_FILE "$HOME"/.mantrachain/config

# Configure app.toml
######################################################################
sed -i -E 's|minimum-gas-prices = \"\"|minimum-gas-prices = \"0.0001axom\"|g' "$HOME"/.mantrachain/config/app.toml
sed -i -E 's|enable-unsafe-cors = false|enable-unsafe-cors = true|g' "$HOME"/.mantrachain/config/app.toml
sed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' "$HOME"/.mantrachain/config/app.toml
sed -i -E '117s/.*/enable = true/' "$HOME"/.mantrachain/config/app.toml
sed -i -E '120s/.*/swagger = true/' "$HOME"/.mantrachain/config/app.toml

# Configure config.toml
######################################################################
sed -i -E '91s/.*/laddr = \"tcp://0.0.0.0:26657\"/' "$HOME"/.mantrachain/config/config.toml
sed -i -E 's|cors_allowed_origins = \[\]|cors_allowed_origins = \[\"*\"\]|g' "$HOME"/.mantrachain/config/config.toml
sed -i -E 's|pex = \[\]|pex = \"true\"|g' "$HOME"/.mantrachain/config/config.toml
sed -i -E 's|persistent_peers = \[\]|persistent_peers = \"'$PERSISTENT_PEERS'\"|g' "$HOME"/.mantrachain/config/config.toml
sed -i -E 's|private_peer_ids = \[\]|private_peer_ids = \"'$PRIVATE_PEERS'\"|g' "$HOME"/.mantrachain/config/config.toml
sed -i -E 's|unconditional_peer_ids = \[\]|unconditional_peer_ids = \"'$UNCONDITIONAL_PEER_IDS'\"|g' "$HOME"/.mantrachain/config/config.toml
sed -i -E 's|addr_book_strict = \[\]|addr_book_strict = \"false\"|g' "$HOME"/.mantrachain/config/config.toml

# Configure client.toml
######################################################################
sed -i -E 's|chain-id = \"\"|chain-id = \"'"$CHAIN_ID"'"|g' "$HOME"/.mantrachain/config/client.toml

printf "Done! Start the node with mantrachaind start\n"