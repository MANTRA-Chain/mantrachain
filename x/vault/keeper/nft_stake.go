package keeper

import (
	"strconv"
	"strings"

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

// TODO: move out the delegate mantrachain logic
func (k Keeper) CreateNftStakeStaked(
	ctx sdk.Context,
	marketplaceCreator string,
	marketplaceId string,
	collectionCreator string,
	collectionId string,
	nftId string,
	marketplaceIndex []byte,
	collectionIndex []byte,
	index []byte,
	creator sdk.AccAddress,
	amount sdk.Coin,
	stakingChain string,
	stakingValidator string,
	cw20ContractAddress sdk.AccAddress,
) error {
	var delegated bool = false
	nftStake, found := k.GetNftStake(ctx, marketplaceIndex, collectionIndex, index)
	delegate := strings.TrimSpace(stakingChain) == "" && strings.TrimSpace(stakingValidator) == ""

	if !found {
		nftStake = types.NftStake{
			Index:            index,
			MarketplaceIndex: marketplaceIndex,
			CollectionIndex:  collectionIndex,
			Staked:           []*types.NftStakeListItem{},
			Balances:         []*types.NftStakeBalance{},
			Creator:          creator,
		}
	}

	stakeAmount := sdk.NewDecFromInt(amount.Amount)

	staked := types.NftStakeListItem{
		Index:   uint32(len(nftStake.Staked)),
		Amount:  &stakeAmount,
		Denom:   amount.Denom,
		Creator: creator,
	}

	// Delegate stake on mantrachain validator
	if delegate {
		params := k.GetParams(ctx)

		se := NewStakingExecutor(ctx, k.ac, k.bk, k.sk)
		shares, err := se.Delegate(creator, amount, params.StakingValidatorAddress)

		if err != nil {
			return err
		}

		staked.StakedAt = ctx.BlockHeader().Time.Unix()
		staked.Chain = ctx.ChainID()
		staked.BlockHeight = ctx.BlockHeight()
		staked.Validator = params.StakingValidatorAddress
		staked.Shares = shares.String()

		// TODO: Add shares to chainValidatorBridge instead of setting them on SetEpochEnd

		lastEpochBlock, found := k.GetLastEpochBlock(ctx, ctx.ChainID(), params.StakingValidatorAddress)

		if !found {
			return sdkerrors.Wrap(types.ErrLastEpochBlockNotFound, "last epoch block not found")
		}

		staked.StakedEpoch = lastEpochBlock.BlockHeight

		delegated = true
	} else { // If the stake will be on a another chain
		staked.Chain = stakingChain
		staked.Validator = stakingValidator
		staked.Shares = "0"
		staked.Cw20ContractAddress = cw20ContractAddress.String()
	}

	nftStake.Staked = append(nftStake.Staked, &staked)

	k.SetNftStake(ctx, nftStake)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeNftStakeStakedCreated),
			sdk.NewAttribute(types.AttributeKeyMarketplaceCreator, marketplaceCreator),
			sdk.NewAttribute(types.AttributeKeyMarketplaceId, marketplaceId),
			sdk.NewAttribute(types.AttributeKeyCollectionCreator, collectionCreator),
			sdk.NewAttribute(types.AttributeKeyCollectionId, collectionId),
			sdk.NewAttribute(types.AttributeKeyNftId, nftId),
			sdk.NewAttribute(types.AttributeKeyChain, staked.Chain),
			sdk.NewAttribute(types.AttributeKeyValidator, staked.Validator),
			sdk.NewAttribute(types.AttributeKeyDelegated, strconv.FormatBool(delegated)),
			sdk.NewAttribute(types.AttributeNftStakeStakedIndex, strconv.Itoa(int(staked.Index))),
		),
	)

	return nil
}
