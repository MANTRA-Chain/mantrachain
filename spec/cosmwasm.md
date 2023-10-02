# CosmWasm & Mantrachain Tutorial

## Setup

### Rust

Install Rust using rustup with the following command and follow the prompts:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

#### Contract environment

Set `stable` as the default release channel:

```bash
rustup default stable
```

Add WASM as the compilation target:

```bash
rustup target add wasm32-unknown-unknown
```

Install the following packages to generate the contract:

```bash
cargo install cargo-generate --features vendored-openssl
cargo install cargo-run-script
```

### Chain

Install Go:

- [Go](https://golang.org/doc/install) (**version 1.19** or higher)

Install [Ignite CLI](https://docs.ignite.com/welcome/install) (version 0.27.1):

```bash
curl https://get.ignite.com/cli!  |  bash
```

Clone the chain:

```bash
git clone https://github.com/MANTRA-Finance/mantrachain
```

Build the chain:

```bash
cd mantrachain
make build
```

Run the chain:

```bash
ignite chain serve -v
```

If you have issues with the chain not finding `libwasmvm.x86_64.so` you can install it manually with:

```bash
sudo wget -P /usr/lib https://github.com/CosmWasm/wasmvm/releases/download/v1.3.0/libwasmvm.x86_64.so
```

(Optionally) You can enable the mantrachain guard module:

```bash
./scripts/init-guard.sh
```

### CosmWasm

#### Create new repo from template

In another terminal execute where `CW_PROJECT_NAME` should be the name of your smart contract:

```bash
CW_PROJECT_NAME={CW_PROJECT_NAME} # e.g. CW_PROJECT_NAME=cw20-contract

cargo generate --git https://github.com/CosmWasm/cw-template.git --name $CW_PROJECT_NAME
```

You will now have a new folder called `CW_PROJECT_NAME` containing a simple working contract and build system that you can customize.

#### Compiling and running tests

```bash
cd $CW_PROJECT_NAME # enter the contract directory

# this will produce a wasm build in ./target/wasm32-unknown-unknown/release/{CW_PROJECT_NAME}.wasm
cargo wasm

# this runs unit tests with helpful backtraces
RUST_BACKTRACE=1 cargo unit-test

# auto-generate json schema
cargo schema
```

#### Preparing the Wasm bytecode

Before we upload it to a chain, we need to ensure the smallest output size possible, as this will be included in the body of a transaction. To achieve this we use `rust-optimizer`, a docker image to produce an extremely small build output in a consistent manner.

```bash
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/rust-optimizer:0.14.0
```

Or, If you're on an arm64 machine, you should use a docker image built with arm64.

```bash
docker run --rm -v "$(pwd)":/code \
 --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
 --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
 cosmwasm/rust-optimizer-arm64:0.14.0
```

This produces an `artifacts` directory with a `{CW_PROJECT_NAME}.wasm`.

#### Deploy the contract

Navigate back to the chain's directory:

```bash
cd ..
```

Set some environment variables:

```bash
CHAINID="mantrachain-9001"
HOMEDIR="$HOME/.mantrachain"
KEYRING="test"
GAS_ADJ=2
GAS_PRICE=0.0002uaum
ACCOUNT_PRIVILEGES_GUARD_NFT_COLLECTION_ID="account_privileges_guard_nft_collection"
ADMIN_WALLET=$(./build/mantrachaind keys show admin -a --keyring-backend $KEYRING --home "$HOMEDIR")
```

Deploy the contract code on chain:

```bash
./build/mantrachaind tx wasm store ./$CW_PROJECT_NAME/artifacts/$CW_PROJECT_NAME.wasm --from admin --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes
```

You should see the tx hash output, e.g.

```bash
...
txhash: 76D3AFAC4382D3ACA42854C859360FBCED8DFF58FFD062DBE07CE1A4BEACD60C
```

Store your contract code-id in an environment variable where `TX_HASH` should be the tx hash from the previous step:

```bash
CODE_ID=$(./build/mantrachain query tx {TX_HASH} --output json | jq -r '.logs[0].events[-1].attributes[0].value')
```

Set the init contract params:

```bash
INIT={\"count\":0}
```

Instantiate the contract where `SOME_LABEL` should be your smart contract label, e.g. cw-20:

```bash
./build/mantrachaind tx wasm instantiate $CODE_ID "$INIT" --label "{SOME_LABEL}" --admin admin --from admin --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes
```

Set the contract address to environment variable:

```bash
CONTRACT_ADDRESS=$(./build/mantrachain query wasm list-contract-by-code $CODE_ID --output json | jq -r '.contracts[0]')
```

## Send some amount of custom coin to the contract

Create a custom coin, mint some amount of it and transfer the amount to the contract address.

### Custom coin

Create a denom where `CUSTOM_COIN_SUBDENOM` should be the subdenom of your custom coin, e.g. usdc:

```bash
CUSTOM_COIN_SUBDENOM={CUSTOM_COIN_SUBDENOM} # e.g. CUSTOM_COIN_SUBDENOM=usdc

./build/mantrachaind tx coinfactory create-denom $CUSTOM_COIN_SUBDENOM --from admin --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes
```

Mint some amount of coins:

```bash
./build/mantrachaind tx coinfactory mint 1000000factory/$ADMIN_WALLET/$CUSTOM_COIN_SUBDENOM --from admin --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes
```

#### Give the smart contract guard privileges, so it can send/receive the custom coin **_(Note: skip this step if the guard module hasn't been enabled earlier in this tutorial)._**

First set environment variable with the soul-bond nft payload to create such an nft with the smart contract for an owner:

```bash
ACCOUNT_PRIVILEGES_GUARD_NFT_CONTRACT_JSON=$(echo '{"id":"{id}","title":"AccountPrivileges","description":"AccountPrivileges"}' | sed -e "s/{id}/$CONTRACT_ADDRESS/g")
```

Mint the soul-bond nft:

```bash
./build/mantrachaind tx token mint-nft "$(echo $ACCOUNT_PRIVILEGES_GUARD_NFT_CONTRACT_JSON)" --collection-creator $ADMIN_WALLET --collection-id $ACCOUNT_PRIVILEGES_GUARD_NFT_COLLECTION_ID --chain-id $CHAINID --from admin --receiver $CONTRACT_ADDRESS --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR" --yes
```

In the next examples we will use the mantrachain-sdk.

### Important: **_Skip them if the guard module hasn't been enabled earlier in this tutorial._**

#### Setup the node.js project directory

Install Node.js:

- [Node.js](https://https://nodejs.org) (**version 16** or higher)

Create new empty node.js project where `BE_PROJECT_NAME` is your node.js app name:

```bash
mkdir {BE_PROJECT_NAME} # e.g. mkdir cw20-be
cd {BE_PROJECT_NAME} # e.g. cd cw20-be
npm init -y
```

Add empty js file:

```bash
touch app.js
```

Install the needed dependencies:

```bash
npm i --save @mantrachain/sdk @cosmjs/proto-signing
```

Open the app.js in a code editor e.g. in a [Visual Studio Code](https://code.visualstudio.com/download):

```bash
code .
```

Add the following lines in `app.js` so you can initialize the sdk:

```js
const { Client } = require("@mantrachain/sdk");
const { DirectSecp256k1HdWallet } = require("@cosmjs/proto-signing");

(async () => {
  const mnemonic = "..."; // MNEMONIC
  const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
    prefix: "mantra",
  });

  const creator = (await wallet.getAccounts())[0];

  const client = new Client(
    {
      apiURL: "http://127.0.0.1:1317",
      rpcURL: "ws://127.0.0.1:26657",
    },
    wallet
  );
})();
```

You should replace `const mnemonic = "...";` with the `admin` mnemonic from the `config.yml` file from the chain directory.
Set some additional constants:

```js
const coinDenom = "factory/{ADMIN_WALLET}/{CUSTOM_COIN_SUBDENOM}"; // e.g. const coinDenom  =  "factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/usdc";
const contractAddress = { CONTRACT_ADDRESS }; // e.g. const contractAddress = "mantra1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrq8pey27";
```

In our case the admin wallet is the current wallet address so you can set the `coinDenom` to:

```js
const coinDenom = `factory/${creator.address}/{CUSTOM_COIN_SUBDENOM}`;
```

Next, add the required privilejes and account privileges resectively to the custom coin and the smart contract:

```js
const { Client, Privileges, utils } =  require("@mantrachain/sdk");
...
(async () => {
...
  const coinPrivileges = Privileges.Empty().set(64).toBuffer();
  await client.MantrachainGuardV1.tx.sendMsgUpdateRequiredPrivileges({
    value: {
      creator: creator.address,
      index: utils.strToIndex(coinDenom),
      privileges: coinPrivileges,
      kind: "coin",
    },
  });

  const res = await client.MantrachainGuardV1.query.queryParams();
  const defaultPrivileges = utils.base64ToBytes(
    res.data.params.default_privileges
  );
  const accountPrivileges = Privileges.fromBuffer(defaultPrivileges)
    .set(64)
    .toBuffer();

  await client.MantrachainGuardV1.tx.sendMsgUpdateAccountPrivileges({
    value: {
      creator: creator.address,
      account: contractAddress,
      privileges: accountPrivileges,
    },
  });
...
})();
```

Save the file and execute it in the cli:

```bash
node app.js
```

Press `Ctrl+C` to exit.

Now you can send some amount of the custom coin to the smart contract:

```bash
cd .. # navigate back to the chain directory

./build/mantrachaind tx bank send  $ADMIN_WALLET  $CONTRACT_ADDRESS 1000000factory/$ADMIN_WALLET/$CUSTOM_COIN_SUBDENOM --from admin --chain-id $CHAINID --keyring-backend $KEYRING --gas auto --gas-adjustment $GAS_ADJ --gas-prices $GAS_PRICE --home "$HOMEDIR"
```

Check the contract balance:

```bash
./build/mantrachaind q bank balances $CONTRACT_ADDRESS
```

You should see the following output in the console:

```bash
balances:
- amount: "100000"
  denom: factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/usdc
pagination:
  next_key: null
  total: "0"
```
