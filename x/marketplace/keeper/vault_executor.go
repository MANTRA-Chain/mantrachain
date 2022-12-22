package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	vaulttypes "github.com/LimeChain/mantrachain/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VaultExecutor struct {
	ctx         sdk.Context
	vaultKeeper types.VaultKeeper
}

func NewVaultExecutor(ctx sdk.Context, vaultKeeper types.VaultKeeper) *VaultExecutor {
	return &VaultExecutor{
		ctx:         ctx,
		vaultKeeper: vaultKeeper,
	}
}

func (c *VaultExecutor) CreateAndDelegateNftStakeStaked(
	marketplaceCreator string,
	marketplaceId string,
	collectionCreator string,
	collectionId string,
	nftId string,
	marketplaceIndex []byte,
	collectionIndex []byte,
	nftIndex []byte,
	creator sdk.AccAddress,
	amount sdk.Coin,
	stakingChain string,
	stakingValidator string,
	cw20ContractAddress sdk.AccAddress,
	nftMarketplaceEarningsOnYieldReward []*types.MarketplaceEarning,
) error {
	nftVaultEarningsOnYieldReward := []*vaulttypes.VaultEarning{}

	for _, k := range nftMarketplaceEarningsOnYieldReward {
		nftVaultEarningsOnYieldReward = append(nftVaultEarningsOnYieldReward, &vaulttypes.VaultEarning{
			Type:       k.Type,
			Address:    k.Address,
			Percentage: k.Percentage,
		})
	}

	return c.vaultKeeper.CreateAndDelegateNftStakeStaked(
		c.ctx,
		marketplaceCreator,
		marketplaceId,
		collectionCreator,
		collectionId,
		nftId,
		marketplaceIndex,
		collectionIndex,
		nftIndex,
		creator,
		amount,
		stakingChain,
		stakingValidator,
		cw20ContractAddress,
		nftVaultEarningsOnYieldReward,
	)
}
