package keeper

import (
	"github.com/MANTRA-Finance/mantrachain/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Keeper of the nft store
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
	bk       types.BankKeeper
}

// NewKeeper creates a new nft Keeper instance
func NewKeeper(key sdk.StoreKey,
	cdc codec.BinaryCodec, ak types.AccountKeeper, bk types.BankKeeper) Keeper {
	// ensure nft module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the nft module account has not been set")
	}

	return Keeper{
		cdc:      cdc,
		storeKey: key,
		bk:       bk,
	}
}
