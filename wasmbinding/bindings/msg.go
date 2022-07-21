package bindings

type MantraMsg struct {
	Mdb *Mdb `json:"mdb,omitempty"`
}

type Mdb struct {
	CreateNftCollection *CreateNftCollection `json:"create_nft_collection,omitempty"`
}

type CreateNftCollection struct {
	Collection *NftCollectionMetadata `json:"collection"`
}

type NftCollectionMetadata struct {
	Id string `json:"id"`
}
