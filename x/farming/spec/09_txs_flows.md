<!-- order: 9 -->

# Transactions flows

## Create Fixed Amount Plan

```mermaid
sequenceDiagram
Creator->>+Farming module: Create Fixed Amount Plan Tx
Farming module->>Guard module: Is chain admin?
alt Chain admin
  Farming module->>Bank module: Charge fee
  Farming module->>Farming module: Set fixed amount plan
  Farming module-->>Creator: Success
else Not a chain admin
  Farming module--x-Creator: Error
end
```

**Note**: Only the `chain admin` is authorized to execute this type of transaction.

Create a new fixed amount plan. The plan will be created with the `Pending` status. The plan will be activated after the `start_time` is reached. The plan will be deleted after the `end_time` is reached.
The plan will be terminated if the `termination_address` calls the `Remove Plan` transaction. The plan's termination address is set to the plan creator.

CLI command:

```bash
aumegad tx farming create-private-fixed-plan [plan-file] [flags]
```

## Create Ratio Plan

```mermaid
sequenceDiagram
Creator->>+Farming module: Create Ratio Plan Tx
Farming module->>Guard module: Is chain admin?
alt Chain admin
  Farming module->>Bank module: Charge fee
  Farming module->>Farming module: Set ratio plan
  Farming module-->>Creator: Success
else Not a chain admin
  Farming module--x-Creator: Error
end
```

**Note**: Only the `chain admin` is authorized to execute this type of transaction.

Create a new ratio plan. The plan will be created with the `Pending` status. The plan will be activated after the `start_time` is reached. The plan will be deleted after the `end_time` is reached.
The plan will be terminated if the `termination_address` calls the `Remove Plan` transaction. The plan's termination address is set to the plan creator.

This transaction is only enabled when the `EnableRatioPlan` flag is set to `true`.

## Stake

```mermaid
sequenceDiagram
Creator->>+Farming module: Stake Tx
Farming module->>Bank module: Transfer staking coins from the creator
Note over Farming module, Bank module: The transfer IS restricted by the guard module
Farming module->>Farming module: Set stake
Farming module-->>-Creator: Success
```

Stake coins to a farming plan.

CLI command:

```bash
aumegad tx farming stake [amount] [flags]
```

## Untake

```mermaid
sequenceDiagram
Creator->>+Farming module: Unstake Tx
Farming module->>Bank module: Return staking coins back to the creator
Note over Farming module, Bank module: The transfer IS NOT restricted by the guard module
Farming module->>Farming module: Set unstake
Farming module-->>-Creator: Success
```

Unstake coins from a farming plan.

CLI command:

```bash
aumegad tx farming unstake [amount] [flags]
```

## Harvest

```mermaid
sequenceDiagram
Creator->>+Farming module: Harvest Tx
Farming module->>Bank module: Withdraw farming rewards to the creator
Note over Farming module, Bank module: The transfer IS NOT restricted by the guard module
Farming module-->>-Creator: Success
```

Harvest farming rewards from a farming plan.

CLI command:

```bash
aumegad tx farming harvest [staking-coin-denoms] [flags]
```

## Remove Plan

```mermaid
sequenceDiagram
Creator->>+Farming module: Remove Plan Tx
Farming module-->>Bank module: Return plan creation fee back to the creator
Farming module->>Farming module: Delete plan
Farming module-->>-Creator: Success
```

Remove a farming plan.

CLI command:

```bash
aumegad tx farming remove-plan [plan-id] [flags]
```

## Advance Epoch

```mermaid
sequenceDiagram
Creator->>+Farming module: Advance Epoch Tx
Note left of Farming module: For testing purposes
Farming module->>Guard module: Is chain admin?
alt Chain admin
  Farming module->>Farming module: Is advance epoch enabled?
  alt Advance epoch enabled
    Farming module->>Farming module: Advance epoch
    Farming module-->>Creator: Success
  else Advance epoch disabled
    Farming module--xCreator: Error
  end
else Not a chain admin
  Farming module--x-Creator: Error
end
```

**Note**: Only the `chain admin` is authorized to execute this type of transaction.

Advance the farming module's epoch. This transaction is only enabled when the `EnableAdvanceEpoch` flag is set to `true`.
