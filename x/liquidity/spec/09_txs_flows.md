<!-- order: 9 -->

# Transactions flows

## Create Pair

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Create Pair Tx
Liquidity module->>Guard module: Is chain admin?
alt Chain admin
  Liquidity module->>Liquidity module: Set pair
  Liquidity module->>Bank module: Charge fee
  Liquidity module-->>Creator: Success
else Not a chain admin
  Liquidity module--x-Creator: Error
end
```

## Create Basic Pool/Create Ranged Pool

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Create Basic Pool Tx/Create Ranged Pool Tx
Liquidity module->>Guard module: Is chain admin?
alt Chain admin
  Liquidity module->>Liquidity module: Get pair
  Liquidity module->>Bank module: Deposit pool coins
  Note over Liquidity module, Bank module: The transfer IS restricted by the guard module
  Liquidity module->>Liquidity module: Set basic/ranged pool
  Liquidity module->>Bank module: Charge fee
  Liquidity module->> Bank module: Mint LP coins to creator
  Note over Liquidity module, Bank module: The transfer IS NOT restricted by the guard module
  Liquidity module-->>Creator: Success
else Not a chain admin
  Liquidity module--x-Creator: Error
end
```

## Deposit

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Deposit Tx
Liquidity module->>Liquidity module: Get pool
Liquidity module->>Bank module: Deposit pool coins
Note over Liquidity module, Bank module: The transfer IS restricted by the guard module
Liquidity module->>Bank module: Mint LP coins to creator
Note over Liquidity module, Bank module: The transfer IS NOT restricted by the guard module
Liquidity module-->>-Creator: Success
```

## Withdraw

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Withdraw Tx
Liquidity module->>Liquidity module: Get pool
Liquidity module->>Bank module: Withdraw pool coins
Note over Liquidity module, Bank module: The transfer IS NOT restricted by the guard module
Liquidity module->>Bank module: Burn LP coins from creator
Note over Liquidity module, Bank module: The transfer IS NOT restricted by the guard module
Liquidity module-->>-Creator: Success
```

## Limit Order

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Limit Order Tx
Liquidity module->>Liquidity module: Get pair
Liquidity module->>Bank module: Transfer offer coin from the creator
Note over Liquidity module, Bank module: The transfer IS restricted by the guard module
Liquidity module->>Liquidity module: Create limit order(includes multiple calls)
Liquidity module-->>-Creator: Success
```

## Market Order

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Market Order Tx
Liquidity module->>Liquidity module: Get pair
Liquidity module->>Bank module: Transfer offer coin from the creator
Note over Liquidity module, Bank module: The transfer IS restricted by the guard module
Liquidity module->>Liquidity module: Create market order(includes multiple calls)
Liquidity module-->>-Creator: Success
```

## MM Order

```mermaid
sequenceDiagram
Creator->>+Liquidity module: MM Order Tx
Liquidity module->>Liquidity module: Get pair
Liquidity module->>Bank module: Transfer offer coin from the creator
Note over Liquidity module, Bank module: The transfer IS restricted by the guard module
Liquidity module->>Liquidity module: Create mm order(includes multiple calls)
Liquidity module-->>-Creator: Success
```

## Cancel Order

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Cancel Order
Liquidity module->>Liquidity module: Finish order(includes multiple calls)
Liquidity module-->>-Creator: Success
```

## Cancel All Orders

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Cancel All Orders Tx
Liquidity module->>Liquidity module: Get pair
Liquidity module->>Liquidity module: Finish all orders(includes multiple calls)
Liquidity module-->>-Creator: Success
```

## Cancel MM Order

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Cancel MM Order Tx
Liquidity module->>Liquidity module: Get pair
Liquidity module->>Liquidity module: Finish mm order(includes multiple calls)
Liquidity module-->>-Creator: Success
```
