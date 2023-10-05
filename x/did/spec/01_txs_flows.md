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

Creates a new did document and metadata.

**Note**: Only the `chain admin` is authorized to execute this type of transaction.

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

Updates an existing did document.

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

Adds a new verification method and related verification relationships to a did document.

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

Removes a verification method and related verification relationships from a did document.

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

Overwrites the verification relationships for a verification methods of a did document.

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

Adds a new service to a did document.

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

Removes a service from a did document.

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

Adds a new controller to a did document.

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

Removes a controller from a did document.
