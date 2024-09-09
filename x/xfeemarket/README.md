# `x/xfeemarket`

## Abstract

This document specifies the `xfeemarket` module of mantrachain.

The `xfeemarket` module is to supplement the `feemarket` module by skip-mev.

It provides an implementation for the `DenomResolver` interface, which is 
required to support multiple different fee tokens. It also has a `PostHandler`
that Replaces the default feemarket `PostHandler`  to eliminate tips and burn 
all fees paid by users.

## Contents

* [DenomResolver](#denomresolver)
* [PostHandler](#posthandler)
* [State](#state)
* [Messages](#messages)

## DenomResolver

The current `DenomResolver` implementation allows an authority account to update a 
map of denom => multiplier. This enables the module to calculate gas fees for 
different fee tokens based on a multiplier applied to the base denom's gas fee.

Future versions may leverage the skip-mev `connect` module to dynamically fetch 
oracle prices for the base fee token and the supplied token, enabling more 
accurate gas fee calculations.

## PostHandler

The custom `PostHandler` refunds tips to fee payer and burns/locks the required 
minimum fees. The minimum fees is burned if it is the default fee denom and locked
otherwise.

## State

The `x/xfeemarket` module keeps state of the following primary objects:

1. Denomination Multipliers

In addition, the `x/xfeemarket` module keeps the following indexes to manage the
aforementioned state:

* Denomination Multipliers: `0x1 | byte(denom) -> sdk.DecProto(multiplier)`

## Messages

### MsgUpsertFeeDenom

Upsert a multiplier for denom.

```protobuf reference
https://github.com/MANTRA-Chain/mantrachain/proto/mantrachain/xfeemarket/v1/tx.proto#L41-L46
```

The message handling can fail if:

* signer is not the gov module account address.
* denom Metadata does not exists.

### MsgRemoveFeeDenom

Remove a multiplier for denom.

```protobuf reference
https://github.com/MANTRA-Chain/mantrachain/proto/mantrachain/xfeemarket/v1/tx.proto#L50-L54
```

The message handling can fail if:

* signer is not the gov module account address.
* denom is not stored in the current Denomination Multipliers