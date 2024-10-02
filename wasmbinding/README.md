# CosmWasm support

This package contains CosmWasm integration points.

## NOTE

This work has been adapted from Osmosis, and was placed into `mantrachain` at commit: 10c8eab.

We've selected the Osmosis approach due to the fixes to non-determinism when working with SDK modules.


This package provides first class support for queries from the token factory. 

## Command line interface (CLI)

- Commands

```sh
  mantrachaind tx wasm -h
```

- Query

```sh
  mantrachaind query wasm -h
```

## Tests

This contains a few high level tests that `x/wasm` is properly
integrated.

Since the code tested is not in this repo, and we are just testing the
application integration (app.go), I figured this is the most suitable
location for it.
