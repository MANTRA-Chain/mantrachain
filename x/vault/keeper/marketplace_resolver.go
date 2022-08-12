package keeper

import (
	marketplacetypes "github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MarketplaceResolver struct{}

func NewMarketplaceResolver() *MarketplaceResolver {
	return &MarketplaceResolver{}
}

func (c *MarketplaceResolver) GetMarketplaceIndex(creator sdk.AccAddress, id string) []byte {
	return marketplacetypes.GetMarketplaceIndex(creator, id)
}
