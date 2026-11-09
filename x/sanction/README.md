# `x/sanction`

## Abstract

This document specifies the `sanction` module of mantrachain.

The `sanction` module keeps a list of sanctioned accounts and prevents them from doing any transactions on the chain.

Governance and **sanctioners** change the list with the same two messages, `MsgSanctionAddresses` and `MsgUnsanctionAddresses`. Governance sends them through proposals. Sanctioners are addresses governance appoints to send them directly, so an attack can be stopped in one transaction without waiting for a proposal to pass.

## Contents

* [AnteHandler](#antehandler)
* [Sanctioners](#sanctioners)
* [State](#state)
* [Params](#params)
* [Messages](#messages)
* [Events](#events)
* [Queries](#queries)
* [CLI](#cli)
* [Gas](#gas)

## AnteHandler

The custom `AnteHandler` checks if the signer of a transaction, its fee granter or any authz granter is sanctioned.
If so, the transaction is rejected. EVM transactions are checked by sender.
A change to the list applies to every transaction checked after the change commits.

## Sanctioners

Governance owns the sanctioner list with `MsgUpdateSanctioners`, which replaces the whole list. Each sanctioner has a role:

| Actor | Add an address | Remove an address | Revoke an add-only sanctioner | Appointed by |
| --- | --- | --- | --- | --- |
| Governance (signer is the gov authority) | Yes | Yes, any entry | Yes, via `MsgUpdateSanctioners` | Token holders |
| Full sanctioner | Yes | Yes, entries added by sanctioners | Yes | Governance |
| Add-only sanctioner | Yes | No | No | Governance |

Add-only is meant for automated monitoring, so a stolen key can't be used to unblock anyone. Multisig accounts need no special handling.

Protections on the sanctioner path:

* A sanctioner can't sanction another sanctioner, the gov authority, or any module account. Module accounts are recognised both from `x/auth` and from the app's module list, so an address is protected before its module account exists.
* A sanctioner can't remove an entry added by governance, including entries that predate sanctioners. Only governance can.
* `max_addresses_per_msg` caps addresses per message and `max_reason_length` caps the reason. Governance isn't capped and can sanction protected addresses.
* Every message needs a non-empty reason, including governance's.
* Adding an address that's already listed, or removing one that isn't, is a no-op for that address, so batches don't fail on duplicates. A protected address fails the whole message.

An empty sanctioner list leaves the module governance-only, which is how it starts after the v8.9.0 upgrade.

## State

* Params: `p_sanction`
* Blacklist accounts: `0x01 | address`. Unchanged from consensus version 1.
* Sanctioners: `0x02 | address -> Sanctioner`
* Sanction records: `0x03 | address -> SanctionRecord`, holding the actor, its role, the reason and the block height. Entries that predate consensus version 2 have no record and are treated as governance entries.

### Migration from version 1

`Migrate1to2` sets params to their defaults and leaves existing blacklist entries untouched. The sanctioner list starts empty.

## Params

| Param | Default | Meaning |
| --- | --- | --- |
| `max_addresses_per_msg` | 50 | Most addresses one `MsgSanctionAddresses` or `MsgUnsanctionAddresses` from a sanctioner can carry |
| `max_reason_length` | 256 | Most bytes in the reason of a sanctioner's message |

## Messages

### MsgUpdateParams

Governance updates params. Both params must be positive.

### MsgUpdateSanctioners

Governance replaces the sanctioner list. Fails on an invalid address, a duplicate, or an unspecified role.

### MsgSanctionAddresses

The gov authority or any sanctioner adds addresses. Addresses already listed are skipped. For a sanctioner it fails if the caps are exceeded or any address is protected. It fails for any other signer, or for an empty reason. Add-only sanctioners should put the detection rule and the triggering transaction hash in the reason.

### MsgUnsanctionAddresses

The gov authority removes any entry. A full sanctioner removes entries added by sanctioners, and the message fails if any address was added by governance. Addresses not listed are skipped. Fails for an add-only signer.

### MsgAddBlacklistAccounts / MsgRemoveBlacklistAccounts (deprecated)

Removed in consensus version 2: there is no handler, so they can't be submitted. Use `MsgSanctionAddresses` and `MsgUnsanctionAddresses` signed by the gov authority. The types stay registered so past proposals that contain them still decode.

### MsgRevokeSanctioner

A full sanctioner removes an add-only sanctioner from the list. Only governance can add it back or remove a full sanctioner.

## Events

* `EventSanctionListChanged`: once per address added or removed, from either path, with the address, `sanctioned`, actor, actor role (`UNSPECIFIED` for governance), reason and height. Alerting can watch for `actor_role = SANCTIONER_ROLE_ADD_ONLY` so a human reviews every automated entry.
* `EventSanctionersUpdated`: when governance replaces the list, with the new list.
* `EventSanctionerRevoked`: when a full sanctioner revokes an add-only sanctioner.

## Queries

| Query | REST |
| --- | --- |
| Params | `/MANTRA-Chain/mantrachain/sanction/v1/params` |
| Blacklist | `/MANTRA-Chain/mantrachain/sanction/v1/blacklist` |
| Sanctioners, with roles | `/MANTRA-Chain/mantrachain/sanction/v1/sanctioners` |
| Sanctioned addresses, with records | `/MANTRA-Chain/mantrachain/sanction/v1/sanctioned_addresses` |
| Whether one address is sanctioned | `/MANTRA-Chain/mantrachain/sanction/v1/is_sanctioned/{address}` |

## CLI

```sh
# sanctioners
mantrachaind tx sanction sanction-addresses mantra1a...,mantra1b... --reason "rule drain-v1, tx 0xabc" --from bot
mantrachaind tx sanction unsanction-addresses mantra1a... --reason "false positive, MM rebalance" --from manager
mantrachaind tx sanction revoke-sanctioner mantra1bot... --reason "key compromised" --from manager

# queries
mantrachaind q sanction sanctioners
mantrachaind q sanction sanctioned-addresses
mantrachaind q sanction is-sanctioned mantra1a...
```

Governance sends any of the messages through `mantrachaind tx gov submit-proposal`, with the gov module address as `signer` or `authority`. For example, to sanction an address:

```json
{
  "messages": [{
    "@type": "/mantrachain.sanction.v1.MsgSanctionAddresses",
    "signer": "<gov module address>",
    "addresses": ["mantra1a..."],
    "reason": "Proposal: block exploiter of 2026-08-20"
  }]
}
```

And to appoint sanctioners:

```json
{
  "messages": [{
    "@type": "/mantrachain.sanction.v1.MsgUpdateSanctioners",
    "authority": "<gov module address>",
    "sanctioners": [
      {"address": "mantra1manager...", "role": "SANCTIONER_ROLE_FULL"},
      {"address": "mantra1bot...", "role": "SANCTIONER_ROLE_ADD_ONLY"}
    ]
  }]
}
```

## Gas

Measured by `TestSanctionGasAtCap` in `app/sanction_integration_test.go`, with a 256-byte reason and the default store gas config. Tx-size gas uses mainnet's `tx_size_cost_per_byte` of 10 and includes about 150 bytes for the signature and fee.

| Addresses | Execution gas | Tx-size gas | Total | Share of mainnet's 75M block |
| --- | --- | --- | --- | --- |
| 1 | 22,453 | ~5,600 | ~28,000 | 0.04% |
| 10 | 194,560 | ~9,900 | ~204,000 | 0.3% |
| 50 | 959,480 | ~28,700 | ~988,000 | 1.3% |

About 19,200 gas per address, of which about 7,700 is storing the reason with each entry. Ante handler costs such as signature verification come on top.
