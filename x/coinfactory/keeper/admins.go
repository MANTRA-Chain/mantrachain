package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/LimeChain/mantrachain/x/coinfactory/types"
)

func (k Keeper) HasAdmin(ctx sdk.Context, denom string) bool {
	store := k.GetDenomPrefixStore(ctx, denom)
	return store.Has([]byte(types.DenomAuthorityMetadataKey))
}

// GetAuthorityMetadata returns the authority metadata for a specific denom
func (k Keeper) GetAdmin(ctx sdk.Context, denom string) (val sdk.AccAddress, found bool) {
	store := k.GetDenomPrefixStore(ctx, denom)

	if !k.HasAdmin(ctx, denom) {
		return []byte{}, false
	}

	b := store.Get([]byte(types.DenomAuthorityMetadataKey))

	if b == nil {
		return val, false
	}

	metadata := types.DenomAuthorityMetadata{}

	k.cdc.MustUnmarshal(b, &metadata)

	admin := sdk.MustAccAddressFromBech32(metadata.Admin)

	return admin, true
}

// GetAuthorityMetadata returns the authority metadata for a specific denom
func (k Keeper) GetAuthorityMetadata(ctx sdk.Context, denom string) (types.DenomAuthorityMetadata, error) {
	bz := k.GetDenomPrefixStore(ctx, denom).Get([]byte(types.DenomAuthorityMetadataKey))

	metadata := types.DenomAuthorityMetadata{}
	err := proto.Unmarshal(bz, &metadata)
	if err != nil {
		return types.DenomAuthorityMetadata{}, err
	}
	return metadata, nil
}

// setAuthorityMetadata stores authority metadata for a specific denom
func (k Keeper) setAuthorityMetadata(ctx sdk.Context, denom string, metadata types.DenomAuthorityMetadata) error {
	err := metadata.Validate()
	if err != nil {
		return err
	}

	store := k.GetDenomPrefixStore(ctx, denom)

	bz, err := proto.Marshal(&metadata)
	if err != nil {
		return err
	}

	store.Set([]byte(types.DenomAuthorityMetadataKey), bz)
	return nil
}

func (k Keeper) setAdmin(ctx sdk.Context, denom string, admin string) error {
	metadata, err := k.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return err
	}

	metadata.Admin = admin

	return k.setAuthorityMetadata(ctx, denom, metadata)
}
