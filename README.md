<img src="https://global-uploads.webflow.com/62ed98169164a3b640e4a87c/62ee38047fea4e239903f8be_m-chain.svg" loading="lazy" alt="MANTRA Chain logo" class="omni3">

<h1 align="center">
    MANTRACHAIN
</h1>
<div align="center">
    MANTRACHAIN is a blockchain built using Cosmos SDK and Tendermint and created with Ignite CLI
</div>

## Pre-requisites

- [Go](https://golang.org/doc/install) >= 1.19.0
- [Ignite](https://github.com/ignite/cli) = 0.27.1

## Get started

```bash
make build
ignite chain serve -v
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

## Configure

Your blockchain in development can be configured with [config.yml](./config.yml). To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Accounts

`validator` - used for bootstrapping the blockchain logic (the chain will be started with only one validator)

`admin` - Genesis Admin for MANTRACHAIN, set during the genesis initialization, holding the initial control of
guard module and admin for the soul bound nft collection

`recipient` - used for various transactions once the chain is started

### Unit tests

```bash
make test
```

### E2E Tests

Start the chain running locally:

```bash
ignite chain serve -v
```

Setup the chain:

```bash
./scripts/init-guard.sh
./scripts/init-e2e.sh
```

Execute the tests:

```bash
cd tests/e2e
yarn
yarn test
```

## Learn more

### Aumega

- [Docs](./spec/README.md)
- [SDK](https://github.com/MANTRA-Finance/aumega-sdk.git)

### Cosmos SDK

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI Docs](https://docs.ignite.com)
- [Cosmos SDK Docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)
