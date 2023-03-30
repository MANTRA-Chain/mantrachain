#!/bin/bash

source "$PWD"/scripts/common.sh

ADMIN_ACCOUNT=$(./build/mantrachaind keys show admin -a --keyring-backend test)

cecho "YELLOW" "Create FOO"
./build/mantrachaind tx coinfactory create-denom foo --from admin -y

sleep 7

cecho "YELLOW" "Create BAR"
./build/mantrachaind tx coinfactory create-denom bar --from admin -y

sleep 7

cecho "YELLOW" "Mint 1000000000000000 FOO"
./build/mantrachaind tx coinfactory mint "1000000000000000factory/$ADMIN_ACCOUNT/foo" --from admin -y

sleep 7

cecho "YELLOW" "Mint 1000000000000000 BAR"
./build/mantrachaind tx coinfactory mint "1000000000000000factory/$ADMIN_ACCOUNT/bar" --from admin -y

sleep 7

cecho "YELLOW" "Create FOO/BAR pair"
./build/mantrachaind tx liquidity create-pair "factory/$ADMIN_ACCOUNT/foo" "factory/$ADMIN_ACCOUNT/bar" --from admin -y

sleep 7

cecho "YELLOW" "Create FOO/BAR liquidity pool"
./build/mantrachaind tx liquidity create-pool 1 "10000000factory/$ADMIN_ACCOUNT/foo,10000000factory/$ADMIN_ACCOUNT/bar" -y --from admin

sleep 7

cecho "YELLOW" "Trade 10000 FOO for 10000 bar with limit order to set pair last price"
./build/mantrachaind tx liquidity limit-order 1 sell "10000factory/$ADMIN_ACCOUNT/foo" "factory/$ADMIN_ACCOUNT/bar" 0.9 10000 --from admin -y