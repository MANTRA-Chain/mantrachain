package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/LimeChain/mantrachain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

		ac                 types.AccountKeeper
		bk                 types.BankKeeper
		wasmViewKeeper     types.WasmViewKeeper
		wasmContractKeeper types.WasmContractOpsKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	ac types.AccountKeeper,
	bk types.BankKeeper,
	wasmViewKeeper types.WasmViewKeeper,
	wasmContractKeeper types.WasmContractOpsKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		ac: ac, bk: bk,
		cdc:                cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		paramstore:         ps,
		wasmViewKeeper:     wasmViewKeeper,
		wasmContractKeeper: wasmContractKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
