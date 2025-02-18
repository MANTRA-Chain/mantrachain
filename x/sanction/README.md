# `x/sanction`

## Abstract

This document specifies the `sanction` module of mantrachain.

The `sanction` module is to blacklist accounts preventing them from doing any transactions on the chain.

## Contents

* [AnteHandler](#antehandler)
* [State](#state)
* [Messages](#messages)

## AnteHandler

The custom `AnteHandler` checks if the signer of a transaction is blacklisted. 
If the signer is blacklisted, the transaction is rejected.

## State

The `x/sanction` module keeps state of the following primary objects:

1. Blacklist Accounts

In addition, the `x/sanction` module keeps a set of accounts to manage the
aforementioned state:

* Denomination Multipliers: `0x1 | byte(account)`

## Messages

### MsgAddBlacklistAccount

Add a blacklist account.

```protobuf reference
https://github.com/MANTRA-Chain/mantrachain/tree/main/proto/mantrachain/sanction/v1/tx.proto#L24-L30
```

The message handling can fail if:

* signer is not the authority module address, usually it is the gov module account address.
* blacklist account is not an account in Bech32 string.

### MsgRemoveBlacklistAccount

Remove a blacklist account.

```protobuf reference
https://github.com/MANTRA-Chain/mantrachain/tree/main/proto/mantrachain/sanction/v1/tx.proto#L35-L41
```

The message handling can fail if:

* signer is not the authority module address, usually it is the gov module account address.
* blacklist account is not an account in Bech32 string.