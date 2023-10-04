<!-- order: 4 -->

# Transactions flows

## Update Account Privileges/Update Account Privileges batch/Update Account Privileges Grouped Batch

```mermaid
sequenceDiagram
Creator->>+Guard module: Update Account Privileges Tx/Update Account Privileges batch Tx/Update Account Privileges Grouped Batch Tx
Guard module->>Guard module: Is chain admin?
alt Chain admin
  Guard module->>Guard module: Update account privileges
else Not a chain admin
  Guard module--x-Creator: Error
end
```

## Update Required Privileges/Update Required Privileges Batch/Update Required Privileges Grouped Batch

```mermaid
sequenceDiagram
Creator->>+Guard module: Update Required Privileges Tx/Update Required Privileges Batch Tx/Update Required Privileges Grouped Batch Tx
Guard module->>Guard module: Is chain admin?
alt Chain admin
  Guard module->>Guard module: Update required privileges
else Not a chain admin
  Guard module--x-Creator: Error
end
```

## Update Guard Transfer Coins

```mermaid
sequenceDiagram
Creator->>+Guard module: Update Guard Transfer Coins Tx
Guard module->>Guard module: Is chain admin?
alt Chain admin
  Guard module->>Guard module: Update guard transfer coins
else Not a chain admin
  Guard module--x-Creator: Error
end
```

## Update Authz Generic Grant Revoke Batch

```mermaid
sequenceDiagram
Creator->>+Guard module: Update Authz Generic Grant Revoke Batch Tx
Guard module->>Guard module: Is chain admin?
alt Chain admin
  Guard module->>Authz module: Add/remove generic authorization(s)
else Not a chain admin
  Guard module--x-Creator: Error
end
```
