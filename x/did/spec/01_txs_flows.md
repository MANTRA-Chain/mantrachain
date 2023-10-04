<!-- order: 1 -->

# Transactions flows

## Create Did Document

```mermaid
sequenceDiagram
Creator->>+Did module: Create Did Document Tx
Did module->>Guard module: Is chain admin?
alt Chain admin
  Did module->>Did module: Set did document
  Did module->>Did module: Set did metadata
  Did module-->>Creator: Success
else Not a chain admin
  Did module--x-Creator: Error
end
```

## Update Did Document

```mermaid
sequenceDiagram
Creator->>+Did module: Update Did Document Tx
Did module->>Did module: Can update did document?
alt Can update
  Did module->>Did module: Update did document
  Did module-->>Creator: Success
else Cannot update
  Did module--x-Creator: No
end
```

## Add Verification

```mermaid
sequenceDiagram
Creator->>+Did module: Add Verification Tx
Did module->>Did module: Can update did document?
alt Can update
  Did module->>Did module: Add verification
  Did module-->>Creator: Success
else Cannot update
  Did module--x-Creator: No
end
```

## Revoke Verification

```mermaid
sequenceDiagram
Creator->>+Did module: Revoke Verification Tx
Did module->>Did module: Can update did document?
alt Can update
  Did module->>Did module: Revoke verification
  Did module-->>Creator: Success
else Cannot update
  Did module--x-Creator: No
end
```

## Set Verification Relationships

```mermaid
sequenceDiagram
Creator->>+Did module: Set Verification Relationships Tx
Did module->>Did module: Can update did document?
alt Can update
  Did module->>Did module: Set verification relationships
  Did module-->>Creator: Success
else Cannot update
  Did module--x-Creator: No
end
```

## Add Service

```mermaid
sequenceDiagram
Creator->>+Did module: Add Service Tx
Did module->>Did module: Can update did document?
alt Can update
  Did module->>Did module: Add service
  Did module-->>Creator: Success
else Cannot update
  Did module--x-Creator: No
end
```

## Delete Service

```mermaid
sequenceDiagram
Creator->>+Did module: Delete Service Tx
Did module->>Did module: Can update did document?
alt Can update
  Did module->>Did module: Delete service
  Did module-->>Creator: Success
else Cannot update
  Did module--x-Creator: No
end
```

## Add Controller

```mermaid
sequenceDiagram
Creator->>+Did module: Add Controller Tx
Did module->>Did module: Can update did document?
alt Can update
  Did module->>Did module: Add controller
  Did module-->>Creator: Success
else Cannot update
  Did module--x-Creator: No
end
```

## Delete Controller

```mermaid
sequenceDiagram
Creator->>+Did module: Delete Controller Tx
Did module->>Did module: Can update did document?
alt Can update
  Did module->>Did module: Delete controller
  Did module-->>Creator: Success
else Cannot update
  Did module--x-Creator: No
end
```
