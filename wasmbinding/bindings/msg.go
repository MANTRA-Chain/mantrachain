package bindings

type MantraMsg struct {
	Token *Token `json:"token,omitempty"`
}

type Token struct {
	CreateNftCollection *CreateNftCollection `json:"create_nft_collection,omitempty"`
}

type CreateNftCollection struct {
	Collection *NftCollectionMetadata `json:"collection"`
}

type NftCollectionMetadata struct {
	Id string `json:"id"`
}
