# ADR 005: Default coinfactory coins to be permissionless

## Context

All coins created by the coinfactory will be prefixed with `factory/`. If the guard module is switched on, the creator of the coin will be the only EOA that can send and receive the newly created coin as the coin admin (apart from the mantrachain admin). We want to allow these tokens to be transferrable to any EOA by default. The creator need to set the permissions if they do not wish for such behavior after creation of the coin.

## Decision

Coins without permissions set will be transferrable by any internal wallet or EOA.

## Status

Accepted
