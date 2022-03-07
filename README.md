# mantrachain
**mantrachain** is a blockchain built using Cosmos SDK and Tendermint.

## Chain scaffolding

- `starport scaffold chain github.com/LimeChain/mantrachain --no-module`
- remove config.yaml;
- create make file, build command based on cosmos-sdk;

```
#!/usr/bin/make -f

BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/
build-linux:
	GOOS=linux GOARCH=$(if $(findstring aarch64,$(shell uname -m)) || $(findstring arm64,$(shell uname -m)),arm64,amd64) $(MAKE) build

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

.PHONY: build build-linux
```

- `make build`

---

## ****Setting up the keyring****

The keyring holds the private/public keypairs used to interact with a node. For instance, a validator key needs to be set up before running the blockchain node, so that blocks can be correctly signed. The private key can be stored in different locations, called "backends", such as a file or the operating system's own key storage.

`./build/mantrachaind keys add validator --keyring-backend test`

- Put the generated address in a variable for later use.
    
    `MY_VALIDATOR_ADDRESS=$(./build/mantrachaind keys show validator -a --keyring-backend test)`
    
    `echo $MY_VALIDATOR_ADDRESS`
    

## ****Initialize the Chain****

Before actually running the node, we need to initialize the chain, and most importantly its genesis file. This is done with the `init` subcommand:

`./build/mantrachaind init mantrachain --chain-id mantrachain`

## ****Adding Genesis Accounts****

Before starting the chain, you need to populate the state with at least one account. To do so, first [create a new account in the keyring](https://docs.cosmos.network/master/run-node/keyring.html#adding-keys-to-the-keyring) named `validator` under the `test` keyring backend (feel free to choose another name and another backend).

Now that you have created a local account, go ahead and grant it some `stake`tokens in your chain's genesis file. Doing so will also make sure your chain is aware of this account's existence:

`./build/mantrachaind add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000stake`

## **Create a gentx.**

`./build/mantrachaind gentx validator 100000000stake --chain-id mantrachain --keyring-backend test`  

**Add the gentx to the genesis file.**

`./build/mantrachaind collect-gentxs`

A `gentx` does three things:

1. Registers the `validator` account you created as a validator operator account (i.e. the account that controls the validator).
2. Self-delegates the provided `amount` of staking tokens.
3. Link the operator account with a Tendermint node pubkey that will be used for signing blocks. If no `-pubkey` flag is provided, it defaults to the local node pubkey created via the `demod init` command above.

## **Run a Localnet**

Now that everything is set up, you can finally start your node:

`./build/mantrachaind start`

## ****Interacting with the Node****

- ****Using the CLI****
    - Get balance
        
        `./build/mantrachaind query bank balances $MY_VALIDATOR_ADDRESS --chain-id mantrachain`
        
- ****Using gRPC****
    - `grpcurl -plaintext localhost:9090 list` - You should see a list of gRPC services
- ****Using the REST Endpoints****
    - check `app.toml` if the api server is enabled, swagger documentation as well;
    
    ```
    curl -X GET -H "Content-Type: application/json" http://localhost:1317/cosmos/bank/v1beta1/balances/$MY_VALIDATOR_ADDRESS
    ```
    

## SDK modules demo

### Create a second account:

`./build/mantrachaind keys add recipient --keyring-backend test`

`RECIPIENT=$(./build/mantrachaind keys show recipient -a --keyring-backend test)`

The command above creates a local key-pair that is not yet registered on the chain. An account is created the first time it receives tokens from another account.

### Send tokens from the validator account to the `recipient` account using **bank module**

`./build/mantrachaind tx bank send $MY_VALIDATOR_ADDRESS $RECIPIENT 49900000000stake --chain-id mantrachain --keyring-backend test`

Check that the recipient account did receive the tokens.

`./build/mantrachaind query bank balances $RECIPIENT --chain-id mantrachain`

### Delegate some of the stake tokens sent to the `recipient`
 account to the validator using **staking module**

`./build/mantrachaind tx staking delegate $(./build/mantrachaind keys show validator --bech val -a --keyring-backend test) 20000000000stake --from recipient --chain-id mantrachain --keyring-backend test`

Query the total delegations to `validator`.

`./build/mantrachaind query staking delegations-to $(./build/mantrachaind keys show validator --bech val -a --keyring-backend test) --chain-id mantrachain`

You should see two delegations, the first one made from the `gentx`, and the second one you just performed from the `recipient` account.****

### **Governance**

- Check gov parameters:

`./build/mantrachaind query gov params`

```
max_deposit_period: "172800000000000"
  min_deposit:
  - amount: "10000000"
    denom: stake
tally_params:
  quorum: "0.334000000000000000"
  threshold: "0.500000000000000000"
  veto_threshold: "0.334000000000000000"
voting_params:
  voting_period: "172800000000000"
```

- Submit proposal with an insufficient deposit:

`./build/mantrachaind tx gov submit-proposal --title="Demo Proposal" --description="Demo, testing, 1, 2, 3" --type="Text" --deposit="5000000stake" --keyring-backend test --chain-id mantrachain --from recipient`

- Check proposal info:

`./build/mantrachaind query gov proposal 1`

- Deposit the remaining part of the min deposit for the proposal:

`./build/mantrachaind tx gov deposit 1 5000000stake --keyring-backend test --chain-id mantrachain --from recipient`

- Check the proposal status again:

`./build/mantrachaind query gov proposal 1`

- Vote for the proposal with the validator account:

`./build/mantrachaind tx gov vote 1 yes --keyring-backend test --chain-id mantrachain --from validator`

- Check the votes:

`./build/mantrachaind query gov tally 1`

## Demonstrate the increase of total supply

[http://localhost:1317/#/Query/CosmosBankV1Beta1TotalSupply](http://localhost:1317/#/Query/CosmosBankV1Beta1TotalSupply)

## Cleaning the demo environment

*Note: if you want to clear the state of the chain and demonstrate the entire flow again, run these commands:*

```
rm -rf build                  # from the project directory
rm -rf ~/.mantrachain                # default directory for config files and keyring
```