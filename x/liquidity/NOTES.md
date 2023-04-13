Get your local "ADMIN" account address:
```
./build/mantrachaind keys show admin --keyring-backend test
```



### Mint Tokens
```
./build/mantrachaind tx coinfactory create-denom foo --from admin
```
```
./build/mantrachaind tx coinfactory create-denom bar --from admin
```
```
./build/mantrachaind tx coinfactory mint "1000000000000000factory/$(./build/mantrachaind keys show admin -a)/foo" --from admin

```
```
./build/mantrachaind tx coinfactory mint "1000000000000000factory/$(./build/mantrachaind keys show admin -a)/bar" --from admin

```

## Create Pair and Pool
```
./build/mantrachaind tx liquidity create-pair "factory/$(./build/mantrachaind keys show admin -a)/foo" "factory/$(./build/mantrachaind keys show admin -a)/bar" --from admin
```
```
./build/mantrachaind query liquidity pairs
```
```
./build/mantrachaind tx liquidity create-pool 1 "10000000factory/$(./build/mantrachaind keys show admin -a)/foo,10000000factory/$(./build/mantrachaind keys show admin -a)/bar" --from admin
```
```
./build/mantrachaind tx liquidity create-ranged-pool 1 "10000000factory/$(./build/mantrachaind keys show admin -a)/foo,10000000factory/$(./build/mantrachaind keys show admin -a)/bar" 0.9 1.1 1 --from admin
```
```
./build/mantrachaind query liquidity pools
```

##### Make Trades/Swaps #####
```
./build/mantrachaind tx liquidity market-order 1 sell "10000factory/$(./build/mantrachaind keys show admin -a)/foo" "factory/$(./build/mantrachaind keys show admin -a)/bar" 10000 --from admin
```
```
./build/mantrachaind tx liquidity limit-order 1 sell "10000factory/$(./build/mantrachaind keys show admin -a)/foo" "factory/$(./build/mantrachaind keys show admin -a)/bar" 0.9 7777 --from admin
```
```
./build/mantrachaind tx liquidity mm-order 1 sell "10000factory/$(./build/mantrachaind keys show admin -a)/foo" "factory/$(./build/mantrachaind keys show admin -a)/bar" 0.9 10000 --from admin
```

## Farming
Edit the dates and create new plan:
```
./build/mantrachaind tx lpfarm create-private-plan \
"New Farming Plan" \
2023-03-17T07:00:00Z \
2023-04-17T07:00:00Z \
pool1:1000uaum pair1:1000uaum --from admin
```
```
./build/mantrachaind tx bank send mantrachain1tleemycal24r8skgk3dmj0vrvpdl274mfvr274 mantrachain15s9vfrppw07n3pdcnx7rnx0amdxn825vh5j4jg 10000uaum --from recipient
```
```
./build/mantrachaind tx lpfarm farm pool1 --from admin
```