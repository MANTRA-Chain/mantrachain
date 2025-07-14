# ADR 006: Allow Unprotected Transactions

## Status

Proposed

## Context

EIP-155 introduced chain-id as part of transaction signatures to prevent transaction replay across different chains. However, this also prevents some useful scenarios:

1. **Deterministic contract deployment**

   Maintaining the same contract address across multiple chains greatly simplifies integrations and applications.

   To deploy a contract to the same address across multiple chains in a permissionless way, the author publishes the pre-signed deploy transaction and the signer address, so on a new chain, anyone can deposit some gas tokens to the signer address, and execute the deploy transaction, without getting permission from the author.

   Tools like [create2 factory](https://github.com/Arachnid/deterministic-deployment-proxy) and [createx factory](https://github.com/pcaversaccio/createx) rely on this pattern.

2. **Contract factories**

   Contract factories are invented to simplify the deterministic contract deployment process. But the factories themselves need to be deployed at same address using the above method.

   For example, only with the same CreateX factory address can one deploy [Uniswap Permit2](https://github.com/Uniswap/permit2) to the same address using the same nonce on a new chain.

This practice is supported in most EVM-compatible chains but not supported in Cosmos EVM with default parameters.

## Decision

We propose setting the `AllowUnprotectedTxs` parameter to `true` to allow processing transactions without chain-id signatures.

### Technical Details

Non-EIP-155 transactions have the following characteristics:

- No chain-id in the signature
- Can be replayed on any chain

### Security Considerations

Wallets for end users still sign EIP-155 transactions by default, and only developers who know what they are doing will use this feature.

## Consequences

### Positive

- Supports deterministic contract deployment patterns
- Improves developer experience and tool compatibility

### Negative

- Theoretically increases replay attack possibilities

## Implementation

1. Set `AllowUnprotectedTxs=true` in network parameters through governance process.
2. Set `allow-unprotected-txs=true` in json-rpc configuration.

## References

- [EIP-155](https://eips.ethereum.org/EIPS/eip-155)
- [Deterministic Deployment Proxy](https://github.com/Arachnid/deterministic-deployment-proxy)
- [CreateX Factory](https://github.com/pcaversaccio/createx)
