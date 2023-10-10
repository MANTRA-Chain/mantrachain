<!-- order: 8 -->

# Transactions flows

## Apply Market Maker

```mermaid
sequenceDiagram
Creator->>+Market Maker module: Apply Market Maker Tx
Market Maker module->>Market Maker module: Set deposit for each pair
Market Maker module->>Market Maker module: Set market maker for each pair
Market Maker module->>Bank module: Deposit coins
Note over Market Maker module, Bank module: The transfer IS restricted by the guard module
Market Maker module-->>-Creator: Success
```

Apply market maker for a list of pairs.

CLI command:

```bash
mantrachaind tx marketmaker apply [pool-ids] [flags]
```

## Claim Incentives

```mermaid
sequenceDiagram
Creator->>+Market Maker module: Claim Incentives Tx
Market Maker module->>Bank module: Withdraw claimable(includes multiple calls)
Note over Market Maker module, Bank module: The transfer IS NOT restricted by the guard module
Market Maker module->>Market Maker module: Delete incentive
Market Maker module-->>-Creator: Success
```

Claim incentives.

CLI command:

```bash
mantrachaind tx liquidfarming claim [flags]
```
