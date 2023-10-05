<!-- order: 7 -->

# Transactions flows

## Create Private Plan

```mermaid
sequenceDiagram
Creator->>+LP Farm module: Create Private Plan Tx
LP Farm module->>Guard module: Is chain admin?
alt Chain admin
  Guard module-->>LP Farm module: Yes
  LP Farm module->>LP Farm module: Set private plan
  LP Farm module->>Bank module: Charge fee
  LP Farm module-->>Creator: Success
else Not a chain admin
  Guard module-->>LP Farm module: No
  LP Farm module--x-Creator: Error
end
```

**Note**: Only the `chain admin` is authorized to execute this type of transaction.

Create a new private farming plan.
The newly created plan's farming pool address is automatically generated and will have no balances in the account initially.
Manually send enough reward coins to the generated farming pool address to make sure that the rewards allocation happens.
The plan's termination address is set to the plan creator.

## Terminate Private Plan

```mermaid
sequenceDiagram
Creator->>+LP Farm module: Terminate Private Plan Tx
LP Farm module->>LP Farm module: Is the plan termination account?
alt Plan termination account
  LP Farm module->>LP Farm module: Terminate private plan
  LP Farm module-->>Creator: Success
else Not the plan termination account
  LP Farm module--x-Creator: Error
end
```

Terminate a private farming plan.

## Farm

```mermaid
sequenceDiagram
Creator->>+LP Farm module: Farm Tx
LP Farm module->>Bank module: Transfer farming coins from the creator
Note over LP Farm module, Bank module: The transfer IS restricted by the guard module
LP Farm module->>LP Farm module: Set farming position
LP Farm module-->>-Creator: Success
```

Add a new farming position to an existing farming plan.

## Unfarm

```mermaid
sequenceDiagram
Creator->>+LP Farm module: Unfarm Tx
LP Farm module->>Bank module: Return farming coins back to the creator
Note over LP Farm module, Bank module: The transfer IS NOT restricted by the guard module
LP Farm module->>Bank module: Withdraw farming rewards to the creator
Note over LP Farm module, Bank module: The transfer IS NOT restricted by the guard module
LP Farm module->>LP Farm module: Delete farming position
LP Farm module-->>-Creator: Success
```

Remove an existing farming position from an existing farming plan.

## Harvest

```mermaid
sequenceDiagram
Creator->>+LP Farm module: Harvest Tx
LP Farm module->>Bank module: Withdraw farming rewards to the creator
Note over LP Farm module, Bank module: The transfer IS NOT restricted by the guard module
LP Farm module-->>-Creator: Success
```

Withdraw farming rewards from an existing farming plan.
