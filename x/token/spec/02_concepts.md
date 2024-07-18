<!-- order: 2 -->

# Concepts

## Nft Collection

A nft collection is a set of nfts. It is identified by a collection creator and a unique id among the creators collections. It has a set of metadata, which includes: `id`, `name`, `images`, `url`, `description`, `links`, `options`, `category` and `symbol`.

The `id` is the collection id specified by the creator. The `category` is the category of the nft collection. Currently, the available categories are: `general`, `art`, `collectibles`, `music`, `photography`, `sports`, `trading-cards`, `utility` and `other`. All the rest of the fields like: `images`, `url`, `links`, `options` are optional and can be used to store additional metadata about the nft collection.
On collection creation there are aditional flags: `opened` and `soul_bonded_nfts` and `restricted_nfts`.
The nft collection can be restricted or unrestricted. If it is restricted, only the chain admin can mint nfts in this collection. Otherwise, if the collection is `opened` anyone can mint nfts in this collection.
If the collection is neither `restricted` nor `opened`, then only the collection creator can mint nfts in this collection. The `soul_bonded_nfts` flag indicates whether the nfts in this collection are soul bonded or not. If the nfts are soul bonded, then the nfts cannot be transfered.

## Nft

A nft is a non-fungible token. It is identified by a collection creator, a collection id and a unique id among the collection nfts. It has a set of metadata, which includes: `id`, `title`, `images`, `url`, `description`, `links` and `attributes`.
There are single operations and a batch operations to mint, burn, transfer and approve nfts among a collection. Also, there are additional flag `did` which indicates whether the nft should have DID or not. The DIDs are automatically generated and stored in the DID module. The DIDs are used to identify the nfts in the chain. The DIDs can be generated only for soul bonded nfts collections.
