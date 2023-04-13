./build/mantrachaind keys show admin --keyring-backend test
mantrachain1ecu6uftl9qykha92fw4xk493g5ptlgnv0lpf30

##### Mint Tokens #####

./build/mantrachaind tx coinfactory create-denom foo --from admin

./build/mantrachaind tx coinfactory create-denom bar --from admin

./build/mantrachaind tx coinfactory mint "1000000000000000factory/mantrachain1rydgmtzjf3nl8xkv6j8hal6xw64u6fykwksa07/foo" --from admin

./build/mantrachaind tx coinfactory mint "1000000000000000factory/mantrachain1rydgmtzjf3nl8xkv6j8hal6xw64u6fykwksa07/bar" --from admin

./build/mantrachaind query bank balances $(./build/mantrachaind keys show admin -a)


##### Create Pair and Pool #####

./build/mantrachaind tx liquidity create-pair "factory/mantrachain1rydgmtzjf3nl8xkv6j8hal6xw64u6fykwksa07/foo" "factory/mantrachain1rydgmtzjf3nl8xkv6j8hal6xw64u6fykwksa07/bar" --from admin

./build/mantrachaind query liquidity pairs


./build/mantrachaind tx liquidity create-pool 1 "10000000factory/mantrachain1rydgmtzjf3nl8xkv6j8hal6xw64u6fykwksa07/foo,10000000factory/mantrachain1rydgmtzjf3nl8xkv6j8hal6xw64u6fykwksa07/bar" --from admin

./build/mantrachaind tx liquidity create-ranged-pool 1 "10000000factory/mantrachain1tleemycal24r8skgk3dmj0vrvpdl274mfvr274/foo,10000000factory/mantrachain1tleemycal24r8skgk3dmj0vrvpdl274mfvr274/bar" 0.9 1.1 1 --from admin

./build/mantrachaind query liquidity pools

##### Make Trades/Swaps #####

./build/mantrachaind tx liquidity market-order 1 sell "10000factory/mantrachain1rydgmtzjf3nl8xkv6j8hal6xw64u6fykwksa07/foo" "factory/mantrachain1rydgmtzjf3nl8xkv6j8hal6xw64u6fykwksa07/bar" 10000 --from admin

./build/mantrachaind tx liquidity limit-order 1 sell "10000factory/mantrachain1rydgmtzjf3nl8xkv6j8hal6xw64u6fykwksa07/foo" "factory/mantrachain1rydgmtzjf3nl8xkv6j8hal6xw64u6fykwksa07/bar" 0.9 7777 --from admin

./build/mantrachaind tx liquidity mm-order 1 sell "10000factory/mantrachain1gaslnwv3xcdnc53qcqz7vc46hy258y53h20pl6/foo" "factory/mantrachain1gaslnwv3xcdnc53qcqz7vc46hy258y53h20pl6/bar" 0.9 10000 --from admin


##### Farming #####

./build/mantrachaind tx lpfarm create-private-plan \
"New Farming Plan" \
2023-03-17T07:00:00Z \
2023-04-17T07:00:00Z \
pool1:1000uaum pair1:1000uaum --from admin

./build/mantrachaind tx bank send mantrachain1tleemycal24r8skgk3dmj0vrvpdl274mfvr274 mantrachain15s9vfrppw07n3pdcnx7rnx0amdxn825vh5j4jg 10000uaum --from recipient

./build/mantrachaind tx lpfarm farm pool1 --from admin


###### Guard ######

./build/mantrachaind tx coinfactory create-denom foo --from admin

./build/mantrachaind tx coinfactory mint "1000000000000000factory/mantrachain1nrqqg8lf5whxqgx49xjmhkphyz38v3gvdrjpex/foo" --from admin

./build/mantrachaind tx bank send mantrachain15s9vfrppw07n3pdcnx7rnx0amdxn825vh5j4jg mantrachain1tleemycal24r8skgk3dmj0vrvpdl274mfvr274 1000000factory/mantrachain1nrqqg8lf5whxqgx49xjmhkphyz38v3gvdrjpex/foo --from admin

./build/mantrachaind tx token mint-nft "{'id':'mantrachain16ly6063z46zvjrk75hwcrdejxrwwahecae5u9e','title': 'AccountPrivileges','description': 'AccountPrivileges'}" --from admin --collection-creator mantrachain1639vunfvl5wpn8q65f7ddpjmdt9x0ymarvvsr3 --collection-id account_privileges_guard_nft_collection --receiver mantrachain16ly6063z46zvjrk75hwcrdejxrwwahecae5u9e

./build/mantrachaind q token all-collection-nfts mantrachain14x90jwjh9m5xlrfmteeg8f8cam3g6w42qy3mww account_privileges_guard_nft_collection