package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace

		didKeeper types.DidKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,

	didKeeper types.DidKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		didKeeper:  didKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
