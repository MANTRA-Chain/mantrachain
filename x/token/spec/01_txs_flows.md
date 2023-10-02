<!-- order: 1 -->

# Transactions flows

## Create nft collection

```mermaid
sequenceDiagram
Creator ->> Token module: Create nft collection
alt Restricted nft collection
  Token module->>Guard module: Is chain admin?
  Note left of Guard module: The guard module checks if the creator is the chain admin.
  alt Chain admin
    Guard module-->> Token module: Yes
    Token module->> Token module: Set nft collection metadata
    Token module->> Nft module: Set nft collection
    Token module-->> Creator: Success
  else Not a chain admin
    Guard module-->> Token module: No
    Token module-->> Creator: Error
  end
else Not restricted nft collection
  Token module->> Token module: Set nft collection metadata
  Token module->> Nft module: Set nft collection
  Token module-->> Creator: Success
end
```

## Mint nft/Batch mint nfts

```mermaid
sequenceDiagram
Creator ->> Token module: Mint nft/Batch mint nfts
alt Restricted nft collection
  Token module->>Guard module: Is chain admin?
  Note left of Guard module: The guard module checks if the creator is the chain admin.
  alt Chain admin
    Guard module-->> Token module: Yes
    Token module->> Token module: Set nft(s) metadata
    Token module->> Nft module: Set nft(s)
    Token module->> Did module: Set did(s)
    Note left of Token module: Set nft(s) did(s) if the did flag is true.
    Token module-->> Creator: Success
  else Not a chain admin
    Guard module-->> Token module: No
    Token module-->> Creator: Error
  end
else Not restricted nft collection
  Token module->> Token module: Set nft(s) metadata
  Token module->> Nft module: Set nft(s)
  Token module-->> Creator: Success
end
```

## Burn nft/Batch burn nfts

```mermaid
sequenceDiagram
Creator ->> Token module: Burn nft/Batch burn nfts
alt Restricted nft collection
  Token module->>Guard module: Is chain admin?
  Note left of Guard module: The guard module checks if the creator is the chain admin.
  alt Chain admin
    Guard module-->> Token module: Yes
    Token module->> Token module: Delete nft(s) metadata
    Token module->> Nft module: Delete nft(s)
    Token module->> Did module: Delete did(s)
    Note left of Token module: Delete nft(s) did(s) if the did(dids) exists.
    Token module-->> Creator: Success
  else Not a chain admin
    Guard module-->> Token module: No
    Token module-->> Creator: Error
  end
else Not restricted nft collection
  Token module->> Token module: Delete nft(s) metadata
  Token module->> Nft module: Delete nft(s)
  Token module-->> Creator: Success
end
```

## Approve nft/Batch approve nfts

```mermaid
sequenceDiagram
Creator ->> Token module: Approve(remove approval) nft/Batch approve(remove approval) nfts
alt Restricted nft collection
  Token module->>Guard module: Is chain admin?
  Note left of Guard module: The guard module checks if the creator is the chain admin.
  alt Chain admin
    Guard module-->> Token module: Yes
    Token module->> Token module: Approve(remove approval) nft(s)
    Token module-->> Creator: Success
  else Not a chain admin
    Guard module-->> Token module: No
    Token module-->> Creator: Error
  end
else Not restricted nft collection
  Token module->> Token module: Approve nft(s)
  Token module-->> Creator: Success
end
```

## Transfer nft/Batch transfer nfts

```mermaid
sequenceDiagram
Creator ->> Token module: Transfer nft/Batch transfer nfts
alt Restricted nft collection
  Token module->>Guard module: Is chain admin?
  Note left of Guard module: The guard module checks if the creator is the chain admin.
  alt Chain admin
    Guard module-->> Token module: Yes
  Token module->> Nft module: Transfer nft(s)
    Token module-->> Creator: Success
  else Not a chain admin
    Guard module-->> Token module: No
    Token module-->> Creator: Error
  end
else Not restricted nft collection
  Token module->> Nft module: Transfer nft(s)
  Token module-->> Creator: Success
end
```
