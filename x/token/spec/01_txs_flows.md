<!-- order: 1 -->

# Transactions flows

## Create Nft Collection

```mermaid
sequenceDiagram
Creator->>+Token module: Nft Collection Tx
Token module->>Token module: Is restricted nft collection?
alt Restricted nft collection
  Token module->>Guard module: Is chain admin?
  alt Chain admin
    Token module->>Token module: Set nft collection metadata
    Token module->>Nft module: Set nft collection
    Token module-->>Creator: Success
  else Not a chain admin
    Token module--xCreator: Error
  end
else Not restricted nft collection
  Token module->>Token module: Set nft collection metadata
  Token module->>Nft module: Set nft collection
  Token module-->>-Creator: Success
end
```

Creates a new nft collection.

## Mint Nft/Batch Mint Nfts

```mermaid
sequenceDiagram
Creator->>+Token module: Mint Nft Tx/Batch Mint Nfts Tx
Token module->>Token module: Is restricted nft collection?
alt Restricted nft collection
  Token module->>Guard module: Is chain admin?
  alt Chain admin
    Token module->>Token module: Set nft(s) metadata
    Token module->>Nft module: Set nft(s)
    Token module->> Did module: Set did(s)
    Note over Token module, Did module: Set nft(s) did(s) if the did flag is true
    Token module-->>Creator: Success
  else Not a chain admin
    Token module--xCreator: Error
  end
else Not restricted nft collection
  Token module->>Token module: Can mint nft(s)?
  alt Can mint nft(s)
    Token module->>Token module: Set nft(s) metadata
    Token module->>Nft module: Set nft(s)
    Token module-->>Creator: Success
  else Cannot mint nft(s)
    Token module--x-Creator: Error
  end
end
```

Mints a new nft.

## Burn Nft/Batch Burn Nfts

```mermaid
sequenceDiagram
Creator->>+Token module: Burn Nft Tx/Batch Burn Nfts Tx
Token module->>Token module: Is restricted nft collection?
alt Restricted nft collection
  Token module->>Guard module: Is chain admin?
  alt Chain admin
    Token module->>Token module: Delete nft(s) metadata
    Token module->>Nft module: Delete nft(s)
    Token module->> Did module: Delete did(s)
    Note over Token module, Did module: Delete nft(s) did(s) if the did(dids) exists
    Token module->>Token module: Delete nft(s) approval(s)
    Token module-->>Creator: Success
  else Not a chain admin
    Token module--xCreator: Error
  end
else Not restricted nft collection
  Token module->>Token module: Can burn nft(s)?
  alt Can burn nft(s)
    Token module->>Token module: Delete nft(s) metadata
    Token module->>Nft module: Delete nft(s)
    Token module->>Token module: Delete nft(s) approval(s)
    Token module-->>Creator: Success
  else Cannot burn nft(s)
    Token module--x-Creator: Error
  end
end
```

Burns an existing nft.

## Approve Nft/Batch Approve Nfts

```mermaid
sequenceDiagram
Creator->>+Token module: Approve(Remove Approval) Nft Tx/Batch Approve(remove approval) Nfts Tx
Token module->>Token module: Is restricted nft collection?
alt Restricted nft collection
  Token module->>Guard module: Is chain admin?
  alt Chain admin
    Token module->>Token module: Approve(remove approval) nft(s)
    Token module-->>Creator: Success
  else Not a chain admin
    Token module--xCreator: Error
  end
else Not restricted nft collection
  Token module->>Token module: Can approve nft(s)?
  alt Can approve nft(s)
    Token module->>Token module: Approve(remove approval) nft(s)
    Token module-->>Creator: Success
  else Cannot approve nft(s)
    Token module--x-Creator: Error
  end
end
```

Adds/removes an approval for an existing nft.

## Approve All Nfts

```mermaid
sequenceDiagram
Creator->>+Token module: Approve(Remove Approval) All Nfts Tx
Token module->>Token module: Approve(remove approval) all nfts
Token module-->>-Creator: Success
```

Adds/removes an approval for all nfts.

## Transfer Nft/Batch Transfer Nfts

```mermaid
sequenceDiagram
Creator->>+Token module: Transfer Nft/Batch Transfer Nfts Tx
Token module->>Token module: Is soul bonded nft collection?
alt Not soul bonded nft collection
  Token module->>Token module: Is restricted nft collection?
  alt Restricted nft collection
    Token module->>Guard module: Is chain admin?
    alt Chain admin
      Token module->>Nft module: Transfer nft(s)
      Token module->>Token module: Delete nft(s) approval(s)
      Token module-->>Creator: Success
    else Not a chain admin
      Token module--xCreator: Error
    end
  else Not restricted nft collection
    Token module->>Token module: Can transfer nft(s)?
    alt Can transfer nft(s)
      Token module->>Nft module: Transfer nft(s)
      Token module->>Token module: Delete nft(s) approval(s)
      Token module-->>Creator: Success
    else Cannot transfer nft(s)
      Token module--xCreator: Error
    end
  end
else Soul bonded nft collection
  Token module--x-Creator: Error
end
```

Transfers an existing nft.

## Update Guard Soul Bond Nft Image

```mermaid
sequenceDiagram
Creator->>+Token module: Update Guard Soul Bond Nft Image Tx
Token module->>Token module: Is restricted nft collection?
alt Restricted nft collection
  Token module->>Guard module: Is chain admin?
  alt Chain admin
    Token module->>Token module: Update guard soul bond nft image
    Token module-->>Creator: Success
  else Not a chain admin
    Token module--xCreator: Error
  end
else Not restricted nft collection
  Token module--x-Creator: Error
end
```

Updates a guard's soul bond nft image.
