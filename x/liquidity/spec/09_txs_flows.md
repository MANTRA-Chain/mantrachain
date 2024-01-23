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

CLI command:

```bash
mantrachaind tx liquidity create-pair [base-coin-denom] [quote-coin-denom] [flags]
```

Example:

```bash
mantrachaind tx liquidity create-pair factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/atom factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/usdc --chain-id mantrachain-9001 --from admin --keyring-backend test --gas auto --gas-adjustment 2 --gas-prices 0.0002uaum --home $HOME/.mantrachain
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

**Note**: Only the `chain admin` is authorized to execute this type of transaction.

CLI commands:

```bash
mantrachaind tx liquidity create-pool [pair-id] [deposit-coins] [flags]
```

```bash
mantrachaind tx liquidity create-ranged-pool [pair-id] [deposit-coins] [min-price] [max-price] [initial-price] [flags]
```

Examples:

```bash
mantrachaind tx liquidity create-pool 1 1000000factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/atom,1000000factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/usdc --chain-id mantrachain-9001 --from admin --keyring-backend test --gas auto --gas-adjustment 2 --gas-prices 0.0002uaum --home $HOME/.mantrachain
```

```bash
mantrachaind tx liquidity create-ranged-pool 1 1000000factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/atom,1000000factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/usdc 0.9 1.1 1 --chain-id mantrachain-9001 --from admin --keyring-backend test --gas auto --gas-adjustment 2 --gas-prices 0.0002uaum --home $HOME/.mantrachain
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

CLI command:

```bash
mantrachaind tx liquidity deposit [pool-id] [deposit-coins] [flags]
```

Example:

```bash
mantrachaind tx liquidity deposit 1 1000factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/atom,1000factory/mantra1axznhnm82lah8qqvp9hxdad49yx3s5dcj66qka/usdc --chain-id mantrachain-9001 --from admin --keyring-backend test --gas auto --gas-adjustment 2 --gas-prices 0.0002uaum --home $HOME/.mantrachain
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

CLI command:

```bash
mantrachaind tx liquidity withdraw [pool-id] [pool-coin] [flags]
```

Example:

```bash
mantrachaind tx liquidity withdraw 1 1000pool1 --chain-id mantrachain-9001 --from admin --keyring-backend test --gas auto --gas-adjustment 2 --gas-prices 0.0002uaum --home $HOME/.mantrachain
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

Create limit order.

CLI command:

```bash
mantrachaind tx liquidity limit-order [pair-id] [direction] [offer-coin] [demand-coin-denom] [price] [amount] [flags]
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

Create market order.

CLI command:

```bash
mantrachaind tx liquidity market-order [pair-id] [direction] [offer-coin] [demand-coin-denom] [amount] [flags]
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

Create mm order.

CLI command:

```bash
mantrachaind tx liquidity mm-order [pair-id] [max-sell-price] [min-sell-price] [sell-amount] [max-buy-price] [min-buy-price] [buy-amount] [flags]
```

## Cancel Order

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Cancel Order
Liquidity module->>Liquidity module: Finish order(includes multiple calls)
Liquidity module-->>-Creator: Success
```

Cancel order.

CLI command:

```bash
mantrachaind tx liquidity cancel-order [pair-id] [order-id] [flags]
```

## Cancel All Orders

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Cancel All Orders Tx
Liquidity module->>Liquidity module: Get pair
Liquidity module->>Liquidity module: Finish all orders(includes multiple calls)
Liquidity module-->>-Creator: Success
```

Cancel all orders.

CLI command:

```bash
mantrachaind tx liquidity cancel-all-orders [pair-ids] [flags]
```

## Cancel MM Order

```mermaid
sequenceDiagram
Creator->>+Liquidity module: Cancel MM Order Tx
Liquidity module->>Liquidity module: Get pair
Liquidity module->>Liquidity module: Finish mm order(includes multiple calls)
Liquidity module-->>-Creator: Success
```

Cancel mm order.

CLI command:

```bash
mantrachaind tx liquidity cancel-mm-order [pair-id] [flags]
```
