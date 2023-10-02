<!-- order: 4 -->

# Transactions flows

## Update account privileges/Update account privileges batch/Update account privileges grouped batch

```mermaid
sequenceDiagram
Creator ->> Guard module: Update account privileges/Update account privileges batch/Update account privileges grouped batch
alt Chain admin
  Guard module->> Guard module: Update account privileges
else Not a chain admin
  Guard module-->> Creator: Error
end
```

## Update required privileges/Update required privileges batch/Update required privileges grouped batch

```mermaid
sequenceDiagram
Creator ->> Guard module: Update required privileges/Update required privileges batch/Update required privileges grouped batch
alt Chain admin
  Guard module->> Guard module: Update required privileges
else Not a chain admin
  Guard module-->> Creator: Error
end
```

## Update guard transfer coins

```mermaid
sequenceDiagram
Creator ->> Guard module: Update guard transfer coins
alt Chain admin
  Guard module->> Guard module: Update guard transfer coins
else Not a chain admin
  Guard module-->> Creator: Error
end
```
