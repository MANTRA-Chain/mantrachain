#!/bin/bash

CHAINID="mantrachain_7000-1"
HOMEDIR="$HOME/.mantrachain"

KEY_VAL_1="validator1"
KEY_RECP_1="recipient1"
KEY_BRG_1="bridge1"
KEY_ADM_1="admin1"
GUARD_NFT_COLLECTION_ID="guardnft"

KEYRING="test"
KEYALGO="eth_secp256k1"

GAS_ADJ=1.4
GAS_PRICE=0.0001axom

./scripts/multinode-local-testnet.sh

source "$PWD"/scripts/common.sh

set -e # exit on first error

cecho "RED" "Do you want to mint guard nft soul-bond collection, some sample marketplace nft collections and WUSDC for a test usage? [y/N]"
read -r mint

if [[ $mint == "y" || $mint == "Y" ]]; then
  VALIDATOR_1_WALLET=$("$PWD"/build/mantrachaind keys show $KEY_VAL_1 -a --keyring-backend=$KEYRING --home=$HOMEDIR/$KEY_VAL_1)
  ADMIN_1_WALLET=$("$PWD"/build/mantrachaind keys show "$KEY_ADM_1" -a --keyring-backend $KEYRING --home "$HOMEDIR/$KEY_VAL_1")

  sleep 7
  
  cecho "GREEN" "Create first nft collection"
  "$PWD"/build/mantrachaind tx token create-nft-collection '{"id":"id1","name":"collection1","description":"collection1","category":"music","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]}' --chain-id $CHAINID --from $KEY_VAL_1 --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

  sleep 7

  cecho "GREEN" "Mint first nft collection nfts"
  "$PWD"/build/mantrachaind tx token mint-nfts '{"nfts":[{"id":"id1","title":"nft1","description":"nft1","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id2","title":"nft2","description":"nft2","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id3","title":"nft3","description":"nft3","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id4","title":"nft4","description":"nft4","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id5","title":"nft5","description":"nft5","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id6","title":"nft6","description":"nft6","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id7","title":"nft7","description":"nft7","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id8","title":"nft8","description":"nft8","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id9","title":"nft9","description":"nft9","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id10","title":"nft10","description":"nft10","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]}]}' --collection-creator $VALIDATOR_1_WALLET --collection-id id1 --chain-id $CHAINID --from $KEY_VAL_1 --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

  sleep 7

  cecho "GREEN" "create second nft collection"
  "$PWD"/build/mantrachaind tx token create-nft-collection '{"id":"id2","name":"collection2","description":"collection2","category":"music","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]}' --chain-id $CHAINID --from $KEY_VAL_1 --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

  sleep 7

  cecho "GREEN" "mint second nft collection nfts"
  "$PWD"/build/mantrachaind tx token mint-nfts '{"nfts":[{"id":"id1","title":"nft1","description":"nft1","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id2","title":"nft2","description":"nft2","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id3","title":"nft3","description":"nft3","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id4","title":"nft4","description":"nft4","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id5","title":"nft5","description":"nft5","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id6","title":"nft6","description":"nft6","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id7","title":"nft7","description":"nft7","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id8","title":"nft8","description":"nft8","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id9","title":"nft9","description":"nft9","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id10","title":"nft10","description":"nft10","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]}]}' --collection-creator $VALIDATOR_1_WALLET --collection-id id2 --chain-id $CHAINID --from $KEY_VAL_1 --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

  sleep 7

  cecho "GREEN" "register marketplace"
  "$PWD"/build/mantrachaind tx marketplace register-marketplace '{"id":"id1","name":"marketplace1","description":"marketplace1","url":"http://xxx.yyy"}' --chain-id $CHAINID --from $VALIDATOR_1_WALLET --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

  sleep 7

  MARKETPLACE=$(echo $("$PWD"/build/mantrachaind query marketplace address) | sed -e "s/address: /${_}/g")

  cecho "GREEN" "approve all nfts"
  "$PWD"/build/mantrachaind tx token approve-all-nfts $MARKETPLACE true --chain-id $CHAINID --from $KEY_VAL_1 --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes
fi

sleep 7

cecho "GREEN" "create $KEY_BRG_1 account"
{
  "$PWD"/build/mantrachaind keys add "$KEY_BRG_1" --keyring-backend=$KEYRING --algo $KEYALGO --home=$HOMEDIR/$KEY_VAL_1
} 2>&1 | tee $KEY_BRG_1.txt

BRIDGE_1_WALLET=$("$PWD"/build/mantrachaind keys show bridge1 -a --keyring-backend=$KEYRING --home=$HOMEDIR/$KEY_VAL_1)

sleep 7

cecho "GREEN" "send axom to bridge1"
"$PWD"/build/mantrachaind tx bank send $KEY_VAL_1 $BRIDGE_1_WALLET 10000000000axom --keyring-backend=$KEYRING --home=$HOMEDIR/$KEY_VAL_1 --chain-id=$CHAINID --gas=auto --gas-adjustment="1.3" --gas-prices="$GAS_PRICE" --yes

sleep 7

cecho "GREEN" "create WUSDC cw20 contract"
"$PWD"/build/mantrachaind tx bridge create-cw-20-contract 0 "1" "./contracts/artifacts/v1/cw20_base-aarch64.wasm" --from $KEY_ADM_1 --chain-id $CHAINID --keyring-backend $KEYRING --gas=auto --gas-adjustment="1.3" --gas-prices="$GAS_PRICE" --home $HOMEDIR/$KEY_VAL_1 --yes

BRIDGE_JSON=$(echo '{"id":"id1","bridge_account":"{address}","cw20_name":"Wrapped USDC","cw20_symbol":"WUSDC","cw20_decimals":6,"cw20_initial_balances":[{"address":"{address}","amount":"0"}],"cw20_mint":{"minter":"{address}"}}' | sed -e "s/{address}/$BRIDGE_1_WALLET/g")

sleep 7

cecho "GREEN" "register bridge"
"$PWD"/build/mantrachaind tx bridge register-bridge "$(echo $BRIDGE_JSON)" --chain-id $CHAINID --from $KEY_BRG_1 --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

sleep 7

RECIPIENT_1_WALLET=$("$PWD"/build/mantrachaind keys show $KEY_RECP_1 -a --keyring-backend=$KEYRING --home=$HOMEDIR/$KEY_VAL_1)

if [[ $mint == "y" || $mint == "Y" ]]; then
  MINT_JSON=$(echo '{"mint_list": [{"receiver":"{address}","amount":"1000000000","tx_hash":"123"}]}' | sed -e "s/{address}/$RECIPIENT_1_WALLET/g")

  cecho "GREEN" "mint WUSDC for $KEY_RECP_1"
  "$PWD"/build/mantrachaind tx bridge mint "$(echo $MINT_JSON)" --chain-id $CHAINID --from $KEY_BRG_1 --keyring-backend $KEYRING --gas auto --bridge-creator $BRIDGE_1_WALLET --bridge-id id1 --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes
fi

CW20_CONTRACT_ADDRESS=$(echo $("$PWD"/build/mantrachaind query bridge bridge $BRIDGE_1_WALLET id1 --output json) | jq '.cw20_contract_address')

sleep 7

if [[ $mint == "y" || $mint == "Y" ]]; then
  FIRST_NFT_COLLECTION=$(echo '{"initially_nft_collection_owner_nfts_for_sale":true,"initially_nft_collection_owner_nfts_min_price":"100000000wusdc","initially_nfts_vault_lock_percentage":"50","cw20_contract_address":{address}}' | sed -e "s/{address}/$CW20_CONTRACT_ADDRESS/g")

  cecho "GREEN" "import first nft collection"
  "$PWD"/build/mantrachaind tx marketplace import-nft-collection "$(echo $FIRST_NFT_COLLECTION)" --chain-id $CHAINID --from $KEY_VAL_1 --keyring-backend $KEYRING --collection-creator $VALIDATOR_1_WALLET --collection-id id1 --marketplace-creator $VALIDATOR_1_WALLET --marketplace-id id1 --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

  sleep 7

  cecho "GREEN" "import second nft collection"
  "$PWD"/build/mantrachaind tx marketplace import-nft-collection '{"initially_nft_collection_owner_nfts_for_sale":true,"initially_nft_collection_owner_nfts_min_price":"1000000axom","initially_nfts_vault_lock_percentage":"50"}' --chain-id $CHAINID --from $KEY_VAL_1 --keyring-backend $KEYRING --collection-creator $VALIDATOR_1_WALLET --collection-id id2 --marketplace-creator $VALIDATOR_1_WALLET --marketplace-id id1 --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

  sleep 7
  
  GUARD_COLL_JSON=$(echo '{"id":"{id}","soul_bonded": true, "name":"GuardCollection","description":"GuardCollection","category":"utility"}' | sed -e "s/{id}/$GUARD_NFT_COLLECTION_ID/g")

  cecho "GREEN" "Create guard nft collection"
  "$PWD"/build/mantrachaind tx token create-nft-collection "$(echo $GUARD_COLL_JSON)" --chain-id $CHAINID --from $KEY_ADM_1 --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --yes  --home $HOMEDIR/$KEY_VAL_1

  sleep 7

  GUARD_NFT_VALIDATOR_JSON=$(echo '{"id":"{id}","title":"GuardNft","description":"GuardNft",}' | sed -e "s/{id}/$VALIDATOR_1_WALLET/g")

  GUARD_NFT_ADMIN_JSON=$(echo '{"id":"{id}","title":"GuardNft","description":"GuardNft",}' | sed -e "s/{id}/$ADMIN_1_WALLET/g")

  cecho "GREEN" "Mint guard nfts"
  "$PWD"/build/mantrachaind tx token mint-nft "$(echo $GUARD_NFT_VALIDATOR_JSON)" --collection-creator $ADMIN_1_WALLET --collection-id $GUARD_NFT_COLLECTION_ID --chain-id $CHAINID --from $KEY_ADM_1 --receiver $VALIDATOR_1_WALLET --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

  sleep 7

  "$PWD"/build/mantrachaind tx token mint-nft "$(echo $GUARD_NFT_ADMIN_JSON)" --collection-creator $ADMIN_1_WALLET --collection-id $GUARD_NFT_COLLECTION_ID --chain-id $CHAINID --from $KEY_ADM_1 --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

  sleep 7

  # cecho "GREEN" "Update guard transfer"
  # "$PWD"/build/mantrachaind tx guard update-guard-transfer true --chain-id $CHAINID --from $KEY_ADM_1 --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

  # sleep 7
fi

cecho "GREEN" "create chain validator bridge"
"$PWD"/build/mantrachaind tx vault create-chain-validator-bridge polygon test $BRIDGE_1_WALLET id1 --chain-id $CHAINID --from $KEY_ADM_1 --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home $HOMEDIR/$KEY_VAL_1 --yes

cecho "CYAN" "CW20 CONTRACT ADDRESS: $CW20_CONTRACT_ADDRESS"
