#!/bin/bash

CHAINID="mantrachain"

while getopts "n:" arg; do
  case $arg in
    n) Op=$OPTARG;;
  esac
done

cecho(){
    RED="\033[0;31m"
    GREEN="\033[0;32m"  # <-- [0 means not bold
    YELLOW="\033[1;33m" # <-- [1 means bold
    CYAN="\033[1;36m"
    # ... Add more colors if you like

    NC="\033[0m" # No Color

    # printf "${(P)1}${2} ${NC}\n" # <-- zsh
    printf "${!1}${2} ${NC}\n" # <-- bash
}

./scripts/multinode-local-testnet.mac.sh

set -e # exit on first error

VALIDATOR1=$(./build/mantrachaind keys show validator1 -a --keyring-backend=test --home=$HOME/.mantrachain/validator1)

sleep 7

if [[ $Op = "mint-sample" ]]
then
  cecho "GREEN" "create first nft collection"
  ./build/mantrachaind tx token create-nft-collection '{"id":"id1","name":"collection1","description":"collection1","category":"music","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]}' --chain-id mantrachain --from validator1 --keyring-backend test --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes

  sleep 7

  cecho "GREEN" "mint first nft collection nfts"
  ./build/mantrachaind tx token mint-nfts '{"nfts":[{"id":"id1","title":"nft1","description":"nft1","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id2","title":"nft2","description":"nft2","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id3","title":"nft3","description":"nft3","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id4","title":"nft4","description":"nft4","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id5","title":"nft5","description":"nft5","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id6","title":"nft6","description":"nft6","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id7","title":"nft7","description":"nft7","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id8","title":"nft8","description":"nft8","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id9","title":"nft9","description":"nft9","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id10","title":"nft10","description":"nft10","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]}]}' --collection-creator $VALIDATOR1 --collection-id id1 --chain-id mantrachain --from validator1 --keyring-backend test --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes

  sleep 7

  cecho "GREEN" "create second nft collection"
  ./build/mantrachaind tx token create-nft-collection '{"id":"id2","name":"collection2","description":"collection2","category":"music","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]}' --chain-id mantrachain --from validator1 --keyring-backend test --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes

  sleep 7

  cecho "GREEN" "mint second nft collection nfts"
  ./build/mantrachaind tx token mint-nfts '{"nfts":[{"id":"id1","title":"nft1","description":"nft1","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id2","title":"nft2","description":"nft2","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id3","title":"nft3","description":"nft3","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id4","title":"nft4","description":"nft4","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id5","title":"nft5","description":"nft5","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id6","title":"nft6","description":"nft6","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id7","title":"nft7","description":"nft7","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id8","title":"nft8","description":"nft8","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id9","title":"nft9","description":"nft9","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]},{"id":"id10","title":"nft10","description":"nft10","url":"http://xxx.yyy","images":[{"type":"main","url":"http://xxx.yyy"}]}]}' --collection-creator $VALIDATOR1 --collection-id id2 --chain-id mantrachain --from validator1 --keyring-backend test --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes

  sleep 7

  cecho "GREEN" "register marketplace"
  ./build/mantrachaind tx marketplace register-marketplace '{"id":"id1","name":"marketplace1","description":"marketplace1","url":"http://xxx.yyy"}' --chain-id mantrachain --from $VALIDATOR1 --keyring-backend test --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes

  sleep 7

  MARKETPLACE=$(echo $(./build/mantrachaind query marketplace address) | sed -e "s/address: /${_}/g")

  cecho "GREEN" "approve all nfts"
  ./build/mantrachaind tx token approve-all-nfts $MARKETPLACE true --chain-id mantrachain --from validator1 --keyring-backend test --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes

  sleep 7
fi

ADMIN1=$(./build/mantrachaind keys show admin1 -a --keyring-backend=test --home=$HOME/.mantrachain/validator1)

sleep 7

cecho "GREEN" "send ustake to admin1"
./build/mantrachaind tx bank send validator1 $ADMIN1 10000000000ustake --keyring-backend=test --home=$HOME/.mantrachain/validator1 --chain-id=$CHAINID --gas=auto --gas-adjustment="1.3" --gas-prices="0.0001ustake" --yes

cecho "GREEN" "create bridge1 account"
{
  ./build/mantrachaind keys add bridge1 --keyring-backend=test --home=$HOME/.mantrachain/validator1
} 2>&1 | tee bridge1.txt
BRIDGE1=$(./build/mantrachaind keys show bridge1 -a --keyring-backend=test --home=$HOME/.mantrachain/validator1)

sleep 7

cecho "GREEN" "send ustake to bridge1"
./build/mantrachaind tx bank send validator1 $BRIDGE1 10000000000ustake --keyring-backend=test --home=$HOME/.mantrachain/validator1 --chain-id=$CHAINID --gas=auto --gas-adjustment="1.3" --gas-prices="0.0001ustake" --yes

sleep 7

cecho "GREEN" "create WUSDC cw20 contract"
./build/mantrachaind tx bridge create-cw-20-contract 0 "1" "./contracts/artifacts/v1/cw20_base-aarch64.wasm" --from admin1 --chain-id mantrachain --keyring-backend test --gas=auto --gas-adjustment="1.3" --gas-prices="0.0001ustake" --home $HOME/.mantrachain/validator1 --yes

BRIDGE_JSON=$(echo '{"id":"id1","bridge_account":"{address}","cw20_name":"Wrapped USDC","cw20_symbol":"WUSDC","cw20_decimals":6,"cw20_initial_balances":[{"address":"{address}","amount":"0"}],"cw20_mint":{"minter":"{address}"}}' | sed -e "s/{address}/$BRIDGE1/g")

sleep 7

cecho "GREEN" "register bridge"
./build/mantrachaind tx bridge register-bridge "$(echo $BRIDGE_JSON)" --chain-id mantrachain --from bridge1 --keyring-backend test --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes

sleep 7

RECIPIENT1=$(./build/mantrachaind keys show recipient1 -a --keyring-backend=test --home=$HOME/.mantrachain/validator1)

if [[ $Op = "mint-sample" ]]
then
  MINT_JSON=$(echo '{"mint_list": [{"receiver":"{address}","amount":"1000000000","tx_hash":"123"}]}' | sed -e "s/{address}/$RECIPIENT1/g")

  cecho "GREEN" "mint WUSDC for recipient1"
  ./build/mantrachaind tx bridge mint "$(echo $MINT_JSON)" --chain-id mantrachain --from bridge1 --keyring-backend test --gas auto --bridge-creator $BRIDGE1 --bridge-id id1 --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes
fi

CW20_CONTRACT_ADDRESS=$(echo $(./build/mantrachaind query bridge bridge $BRIDGE1 id1 --output json) | jq '.cw20_contract_address')

sleep 7

if [[ $Op = "mint-sample" ]]
then
  FIRST_NFT_COLLECTION=$(echo '{"initially_nft_collection_owner_nfts_for_sale":true,"initially_nft_collection_owner_nfts_min_price":"100000000wusdc","initially_nfts_vault_lock_percentage":"50","cw20_contract_address":{address}}' | sed -e "s/{address}/$CW20_CONTRACT_ADDRESS/g")

  cecho "GREEN" "import first nft collection"
  ./build/mantrachaind tx marketplace import-nft-collection "$(echo $FIRST_NFT_COLLECTION)" --chain-id mantrachain --from validator1 --keyring-backend test --collection-creator $VALIDATOR1 --collection-id id1 --marketplace-creator $VALIDATOR1 --marketplace-id id1 --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes

  sleep 7

  cecho "GREEN" "import second nft collection"
  ./build/mantrachaind tx marketplace import-nft-collection '{"initially_nft_collection_owner_nfts_for_sale":true,"initially_nft_collection_owner_nfts_min_price":"1000000ustake","initially_nfts_vault_lock_percentage":"50"}' --chain-id mantrachain --from validator1 --keyring-backend test --collection-creator $VALIDATOR1 --collection-id id2 --marketplace-creator $VALIDATOR1 --marketplace-id id1 --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes

  sleep 7
fi

cecho "GREEN" "create chain validator bridge"
./build/mantrachaind tx vault create-chain-validator-bridge polygon test $BRIDGE1 id1 --chain-id mantrachain --from admin1 --keyring-backend test --gas auto --gas-adjustment 1.4 --gas-prices 0.0001ustake --home $HOME/.mantrachain/validator1 --yes

cecho "GREEN" "CW20 CONTRACT ADDRESS: $CW20_CONTRACT_ADDRESS"
