# ADR 004: Allow users to create permissionless coins

## Context

All coins created by the coinfactory will be prefixed with `factory/`. If the guard module is switched on, the creator of the coin will be the only EOA that can send and receive the newly created coin as the coin admin (apart from the mantrachain admin). We want to add another message that can create permisionless coins that any EOA can freely send and receive.

## Decision

On launch of mainnet, we want to migrate some coins created on testnet into mainnet. These coins should be permissionless. We also made the decision to allow verified users to be able to create permissionless denom that all other users can send and receive.

To achieve this, we design a new Msg type as follows:

```protobuf
service Msg {
    rpc CreatePermissionlessDenom(MsgCreatePermissionlessDenom) returns (MsgCreatePermissionlessDenomResponse);
}

message MsgCreatePermissionlessDenom {
  option (cosmos.msg.v1.signer) = "sender";
  option (amino.name) = "mantrachain/x/coinfactory/MsgCreatePermissionlessDenom";

  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];
  // subdenom can be up to 44 "alphanumeric" characters long.
  string subdenom = 2;
}

message MsgCreatePermissionlessDenomResponse {}
```

This is exactly the same as `MsgCreateDenom`. The transaction will create a new denom that is prefixed with another constant. Tentatively we can use `free/`. Coins with this prefixed denom will skip checks by the guard module.

## Status

Superseded by ADR-005

## Further Discussions

We should decide on a prefix for these permissionless coins.
