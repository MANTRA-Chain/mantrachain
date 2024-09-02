# ADR 006: Standardise our tokenfactory protobuf namespace

## Context

Our current protobuf implementation for the coinfactory is `mantrachain.coinfactory.v1beta1.*`. This module has the same functionalities as the `tokenfactory` in Osmosis but with permissions enforced by the `guard` module. Most Cosmos SDK chains with this module (Osmosis, Neutron) follows the protobuf of osmosis with `osmosis.tokenfactory.v1beta1.*`.

## Decision

In order to allow for CosmWasm developers that are active on chains such as Osmosis and Neutron to port over their contracts without much additional work, we should adopt the same protobuf `osmosis.tokenfactory.v1beta1.*`

## Status

Accepted
