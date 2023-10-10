<!-- order: 8 -->

# Transactions flows

## Liquid Farm

```mermaid
sequenceDiagram
Creator->>+Liquid Farming module: Liquid Farm Tx
Liquid Farming module->>Bank module: Transfer farming coins from the creator
Note over Liquid Farming module, Bank module: The transfer IS restricted by the guard module
Liquid Farming module->>Bank module: Mint LP coins to creator
Note over Liquid Farming module, Bank module: The transfer IS NOT restricted by the guard module
Liquid Farming module-->>-Creator: Success
```

Staking pool coins to the `lpfarm` module and minting LFCoin.

CLI command:

```bash
mantrachaind tx liquidfarming liquid-farm [pool-id] [amount] [flags]
```

## Liquid Unfarm

```mermaid
sequenceDiagram
Creator->>+Liquid Farming module: Liquid Unfarm Tx
Liquid Farming module->>Bank module: Return farming coins back to the creator
Note over Liquid Farming module, Bank module: The transfer IS NOT restricted by the guard module
Liquid Farming module->>Bank module: Withdraw farming rewards to the creator
Note over Liquid Farming module, Bank module: The transfer IS NOT restricted by the guard module
Liquid Farming module->>Bank module: Burn LP coins from creator
Note over Liquid Farming module, Bank module: The transfer IS NOT restricted by the guard module
Liquid Farming module-->>-Creator: Success
```

Unstaking pool coins from the `lpfarm` module and burning LFCoin. Also, the tx claims farming rewards.

CLI command:

```bash
mantrachaind tx liquidfarming liquid-unfarm [pool-id] [amount] [flags]
```

## Liquid Unfarm And Withdraw

```mermaid
sequenceDiagram
Creator->>+Liquid Farming module: Liquid Unfarm And Withdraw Tx
Liquid Farming module->>Bank module: Return farming coins back to the creator
Note over Liquid Farming module, Bank module: The transfer IS NOT restricted by the guard module
Liquid Farming module->>Bank module: Withdraw pool coins(includes call to the liquidity module)
Note over Liquid Farming module, Bank module: The transfer IS NOT restricted by the guard module
Liquid Farming module->>Bank module: Withdraw farming rewards to the creator
Note over Liquid Farming module, Bank module: The transfer IS NOT restricted by the guard module
Liquid Farming module->>Bank module: Burn LP coins from creator
Note over Liquid Farming module, Bank module: The transfer IS NOT restricted by the guard module
Liquid Farming module-->>-Creator: Success
```

Unstaking pool coins from the `lpfarm` module, withdrawing pool coins from the `liquidity` module, and burning LFCoin. Also, the tx claims farming rewards.

CLI command:

```bash
mantrachaind tx liquidfarming liquid-unfarm-and-withdraw [pool-id] [amount] [flags]
```

## Place Bid

```mermaid
sequenceDiagram
Creator->>+Liquid Farming module: Place Bid Tx
Liquid Farming module->>Liquid Farming module: Get Rewards Auction
Liquid Farming module->>Bank module: Refund the previous bid if the bidder has placed bid before
Note over Liquid Farming module, Bank module: The transfer IS NOT restricted by the guard module
Liquid Farming module->>Liquid Farming module: Set Bid
Liquid Farming module->>Bank module: Transfer bidding coin from the creator
Note over Liquid Farming module, Bank module: The transfer IS restricted by the guard module
Liquid Farming module-->>-Creator: Success
```

Placing a bid for the rewards auction.

CLI command:

```bash
mantrachaind tx liquidfarming place-bid [auction-id] [pool-id] [amount] [flags]
```

## Refund Bid

```mermaid
sequenceDiagram
Creator->>+Liquid Farming module: Refund Bid Tx
Liquid Farming module->>Liquid Farming module: Get Rewards Auction
Liquid Farming module->>Bank module: Return bidding coin back to the creator
Note over Liquid Farming module, Bank module: The transfer IS NOT restricted by the guard module
Liquid Farming module->>Liquid Farming module: Delete Bid
Liquid Farming module-->>-Creator: Success
```

Refunding the bid for the rewards auction.

CLI command:

```bash
mantrachaind tx liquidfarming refund-bid [auction-id] [pool-id] [flags]
```

## Advance Auction

```mermaid
sequenceDiagram
Creator->>+Liquid Farming module: Advance Auction Tx
Note left of Liquid Farming module: For testing purposes
Liquid Farming module->>Guard module: Is chain admin?
alt Chain admin
  Liquid Farming module->>Liquid Farming module: Is advance auction enabled?
  alt Advance auction enabled
    Liquid Farming module->>Liquid Farming module: Advance auction(includes calls to LPFarm and Liquidity modules)
    Liquid Farming module-->>Creator: Success
  else Advance auction disabled
    Liquid Farming module--xCreator: Error
  end
else Not a chain admin
  Liquid Farming module--x-Creator: Error
end
```

**Note**: Only the `chain admin` is authorized to execute this type of transaction.

Advancing the rewards auction. This transaction is only enabled when the `EnableAdvanceAuction` flag is set to `true`.
