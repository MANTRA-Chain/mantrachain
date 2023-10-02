<!-- order: 8 -->

# Transactions flows

## Liquid Farm

```mermaid
sequenceDiagram
Creator ->> Liquid Farming module: Liquid Farm
Liquid Farming module->>Bank module: Transfer farming coins from the creator
Liquid Farming module->> Bank module: Mint LP coins to creator
Liquid Farming module-->> Creator: Success
```

## Liquid Unfarm

```mermaid
sequenceDiagram
Creator ->> Liquid Farming module: Liquid Unfarm
Liquid Farming module->>Bank module: Return farming coins back to the creator
Liquid Farming module->>Bank module: Transfer farming rewards to the creator
Liquid Farming module->> Bank module: Burn LP coins from creator
Liquid Farming module-->> Creator: Success
```

## Liquid Unfarm And Withdraw

```mermaid
sequenceDiagram
Creator ->> Liquid Farming module: Liquid Unfarm And Withdraw
Liquid Farming module->>Bank module: Return farming coins back to the creator
Liquid Farming module->>Liquidity module: Withdraw pool coins
Liquid Farming module->>Bank module: Transfer farming rewards to the creator
Liquid Farming module->> Bank module: Burn LP coins from creator
Liquid Farming module-->> Creator: Success
```

## Place Bid

```mermaid
sequenceDiagram
Creator ->> Liquid Farming module: Place Bid
Liquid Farming module->>Liquid Farming module: Get Rewards Auction
Liquid Farming module->>Bank module: Refund the previous bid if the bidder has placed bid before
Liquid Farming module->>Liquid Farming module: Set Bid
Liquid Farming module->>Bank module: Transfer bidding coin from the creator
Liquid Farming module-->> Creator: Success
```

## Refund Bid

```mermaid
sequenceDiagram
Creator ->> Liquid Farming module: Refund Bid
Liquid Farming module->>Liquid Farming module: Get Rewards Auction
Liquid Farming module->>Liquid Farming module: Delete Bid
Liquid Farming module->>Bank module: Return bidding coin back to the creator
Liquid Farming module-->> Creator: Success
```
