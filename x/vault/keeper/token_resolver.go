package keeper

import (
	tokentypes "github.com/LimeChain/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenResolver struct{}

func NewTokenResolver() *TokenResolver {
	return &TokenResolver{}
}

func (c *TokenResolver) GetCollectionIndex(creator sdk.AccAddress, id string) []byte {
	return tokentypes.GetNftCollectionIndex(creator, id)
}

func (c *TokenResolver) GetNftIndex(collectionCreator sdk.AccAddress, collectionId string, id string) []byte {
	collectionIndex := tokentypes.GetNftCollectionIndex(collectionCreator, collectionId)
	return tokentypes.GetNftIndex(collectionIndex, id)
}
