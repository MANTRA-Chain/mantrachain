<!-- order: 8 -->

# Transactions flows

## Apply Market Maker

```mermaid
sequenceDiagram
Creator ->> Market Maker module: Apply Market Maker
Market Maker module->>Market Maker module: Set deposit for each pair
Market Maker module->>Market Maker module: Set market maker for each pair
Market Maker module->>Bank module: Deposit coins
Market Maker module-->> Creator: Success
```

## Claim Incentives

```mermaid
sequenceDiagram
Creator ->> Market Maker module: Claim Incentives
Market Maker module->>Bank module: Withdraw claimable
Market Maker module->>Market Maker module: Delete incentive
Market Maker module-->> Creator: Success
```
