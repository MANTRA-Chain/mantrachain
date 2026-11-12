# `x/sanction`

## Abstract

This document specifies the `sanction` module of mantrachain.

The `sanction` module keeps a list of sanctioned accounts and prevents them from doing any transactions on the chain.

Governance and **sanctioners** change the list with the same two messages, `MsgSanctionAddresses` and `MsgUnsanctionAddresses`. Governance sends them through proposals. Sanctioners are addresses governance appoints to send them directly, so an attack can be stopped in one transaction without waiting for a proposal to pass.

## Contents

* [AnteHandler](#antehandler)
* [Send restriction](#send-restriction)
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
If so, the transaction is rejected. EVM transactions are checked by sender, and a type-4 (EIP-7702) transaction is also rejected if any of its authorizations is signed by a sanctioned account. Otherwise a relayer could submit a sanctioned account's authorization and run code as that account.

When an address is sanctioned, any EIP-7702 delegation on it is cleared, so code it authorized earlier can no longer run as it. Contract code is left alone.
A change to the list applies to every transaction checked after the change commits.

## Send restriction

The ante handler only stops a sanctioned account signing transactions. Its coins can still be moved by paths that skip the ante handler, such as an ICA host transaction, authz executed by a contract, or `transferFrom` on the ERC20 precompile using an allowance granted earlier. A bank send restriction closes these: any bank send out of a sanctioned address fails with `ErrAccountBlacklisted`, whatever started it. This includes EVM value transfers, which lower the account's bank balance when the state commits.

* A sanctioned account can still receive coins.
* Module accounts can't be sanctioned by anyone, governance included, or listed in genesis. Blocking one could stop fee distribution or halt the chain.
* The lookup is not charged gas, so transfer gas costs are unchanged.
* It covers bank coins only. Tokens kept in contract storage (EVM-only ERC20s, cw20) aren't bank coins, so they're only stopped if the token contract blocks them.
* Staking moves coins with delegate and undelegate calls, not sends, so they aren't blocked. Delegated coins stay owned by the account.
* Tokenfactory `ForceTransfer` out of a sanctioned account is blocked as well.

## Sanctioners

Governance owns the sanctioner list with `MsgUpdateSanctioners`, which replaces the whole list. Each sanctioner has a role:

| Actor | Add an address | Remove an address | Revoke an add-only sanctioner | Appointed by |
| --- | --- | --- | --- | --- |
| Governance (signer is the gov authority) | Yes | Yes, any entry | Yes, via `MsgUpdateSanctioners` | Token holders |
| Full sanctioner | Yes | Yes, any entry | Yes | Governance |
| Add-only sanctioner | Yes | No | No | Governance |

Add-only is meant for automated monitoring, so a stolen key can't be used to unblock anyone. Multisig accounts need no special handling.

Rules:

* A sanctioner can't sanction another sanctioner.
* Nobody, governance included, can sanction a module account (the gov authority among them), because the send restriction would stop it sending.
* Governance can sanction a sanctioner.
* Appointing an address that isn't a sanctioner yet lifts any sanction on it, so one placed while the proposal votes can't lock the new sanctioner out.
* `max_addresses_per_msg` caps addresses per message and `max_reason_length` caps the reason, for every signer including governance. A larger governance batch is split across several messages in one proposal, or governance raises the params first.
* Every message needs a non-empty reason, including governance's.
* Adding an address that's already listed, or removing one that isn't, is a no-op for that address, so batches don't fail on duplicates. A protected address fails the whole message.

An empty sanctioner list leaves the module governance-only, which is how it starts after the v8.9.0 upgrade.

## State

* Params: `p_sanction`
* Blacklist accounts: `0x01 | address`. Unchanged from consensus version 1.
* Sanctioners: `0x02 | address -> Sanctioner`

State doesn't record who added an entry or why. That lives only in `EventSanctionListChanged` and in the transaction itself.

### Migration from version 1

`Migrate1to2` sets params to their defaults, removes any module account from the blacklist, and clears any EIP-7702 delegation from the remaining entries. Those entries are otherwise left untouched. The sanctioner list starts empty.

## Params

| Param | Default | Meaning |
| --- | --- | --- |
| `max_addresses_per_msg` | 50 | Most addresses one `MsgSanctionAddresses` or `MsgUnsanctionAddresses` can carry |
| `max_reason_length` | 256 | Most bytes in the reason of either message |

## Messages

### MsgUpdateParams

Governance updates params. Both params must be positive.

### MsgUpdateSanctioners

Governance replaces the sanctioner list. Fails on an invalid address, a duplicate, or an unspecified role.

### MsgSanctionAddresses

The gov authority or any sanctioner adds addresses. Addresses already listed are skipped. It fails if the caps are exceeded, the reason is empty or any address is a module account, and for a sanctioner if any address is another sanctioner. It fails for any other signer. Add-only sanctioners should put the detection rule and the triggering transaction hash in the reason.

### MsgUnsanctionAddresses

The gov authority or a full sanctioner removes addresses, whoever added them, including entries that predate sanctioners. Addresses not listed are skipped. It fails if the caps are exceeded, and for an add-only or any other signer.

### MsgAddBlacklistAccounts / MsgRemoveBlacklistAccounts (deprecated)

Removed in consensus version 2: there is no handler, so they can't be submitted. Use `MsgSanctionAddresses` and `MsgUnsanctionAddresses` signed by the gov authority. The types stay registered so past proposals that contain them still decode.

### MsgRevokeSanctioner

A full sanctioner removes an add-only sanctioner from the list. Only governance can add it back or remove a full sanctioner.

## Events

* `EventSanctionListChanged`: once per `MsgSanctionAddresses` or `MsgUnsanctionAddresses` that changes the list, from either path, and when an appointment unsanctions an address, with the addresses actually added or removed, `sanctioned`, actor, actor role (`UNSPECIFIED` for governance) and reason. Not emitted if every address was a no-op. Alerting can watch for `actor_role = SANCTIONER_ROLE_ADD_ONLY` so a human reviews every automated entry.
* `EventSanctionersUpdated`: when governance replaces the list, with the new list.
* `EventSanctionerRevoked`: when a full sanctioner revokes an add-only sanctioner.

## Queries

| Query | REST |
| --- | --- |
| Params | `/MANTRA-Chain/mantrachain/sanction/v1/params` |
| Blacklist | `/MANTRA-Chain/mantrachain/sanction/v1/blacklist` |
| Sanctioners, with roles | `/MANTRA-Chain/mantrachain/sanction/v1/sanctioners` |
| Whether one address is sanctioned | `/MANTRA-Chain/mantrachain/sanction/v1/is_sanctioned/{address}` |

## CLI

```sh
# sanctioners
mantrachaind tx sanction sanction-addresses mantra1a...,mantra1b... --reason "rule drain-v1, tx 0xabc" --from bot
mantrachaind tx sanction unsanction-addresses mantra1a... --reason "false positive, MM rebalance" --from manager
mantrachaind tx sanction revoke-sanctioner mantra1bot... --reason "key compromised" --from manager

# queries
mantrachaind q sanction sanctioners
mantrachaind q sanction blacklist
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
| 1 | 12,121 | ~5,600 | ~18,000 | 0.02% |
| 10 | 79,675 | ~9,900 | ~90,000 | 0.1% |
| 50 | 379,915 | ~28,700 | ~409,000 | 0.5% |

About 7,600 gas per address, including the module account check and the check for an EIP-7702 delegation. Ante handler costs such as signature verification come on top.
