package keeper

import (
	"github.com/LimeChain/mantrachain/x/vault/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) SetNftStake(ctx sdk.Context, nftStake types.NftStake) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStakeStoreKey(nftStake.MarketplaceIndex, nftStake.CollectionIndex))
	b := k.cdc.MustMarshal(&nftStake)
	store.Set(nftStake.Index, b)
}

func (k Keeper) GetNftStake(
	ctx sdk.Context,
	marketplaceIndex []byte,
	collectionIndex []byte,
	index []byte,
) (val types.NftStake, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStakeStoreKey(marketplaceIndex, collectionIndex))

	if !k.HasNftStake(ctx, marketplaceIndex, collectionIndex, index) {
		return types.NftStake{}, false
	}

	b := store.Get(index)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) HasNftStake(
	ctx sdk.Context,
	marketplaceIndex []byte,
	collectionIndex []byte,
	index []byte,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStakeStoreKey(marketplaceIndex, collectionIndex))
	return store.Has(index)
}

func (k Keeper) GetAllNftStake(ctx sdk.Context, marketplaceIndex []byte, collectionIndex []byte) (list []types.NftStake) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NftStakeStoreKey(marketplaceIndex, collectionIndex))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.NftStake
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) UpsertNftStake(
	ctx sdk.Context,
	marketplaceIndex []byte,
	collectionIndex []byte,
	index []byte,
	creator sdk.AccAddress,
	amount sdk.Coin,
	delegate bool,
) error {
	valAddress := sdk.AccAddress{}
	nftStake, found := k.GetNftStake(ctx, marketplaceIndex, collectionIndex, index)

	if !found {
		nftStake = types.NftStake{
			Index:            index,
			MarketplaceIndex: marketplaceIndex,
			CollectionIndex:  collectionIndex,
			Staked:           []types.Stake{},
			Balances:         []sdk.DecCoin{},
			Creator:          creator,
		}
	}

	staked := types.Stake{
		Amount:             amount.String(),
		Validator:          valAddress.String(),
		Chain:              ctx.ChainID(),
		StakedAt:           ctx.BlockHeader().Time.Unix(),
		Creator:            creator,
		BlockHeight:        ctx.BlockHeight(),
		LastEpochWithdrawn: types.UndefinedBlockHeight,
	}

	if delegate {
		params := k.GetParams(ctx)

		se := NewStakingExecutor(ctx, k.ac, k.bk, k.sk)
		shares, err := se.Delegate(creator, amount, params.StakingValidatorAddress)

		if err != nil {
			return err
		}

		staked.Shares = shares.String()
		staked.Validator = params.StakingValidatorAddress

		lastEpochBlock, found := k.GetLastEpochBlock(ctx, ctx.ChainID(), params.StakingValidatorAddress, params.StakingValidatorDenom)

		if !found {
			return sdkerrors.Wrap(types.ErrLastEpochBlockNotFound, "last epoch block not found")
		}

		staked.Epoch = lastEpochBlock.BlockHeight
	}

	nftStake.Staked = append(nftStake.Staked, staked)

	k.SetNftStake(ctx, nftStake)

	return nil
}
