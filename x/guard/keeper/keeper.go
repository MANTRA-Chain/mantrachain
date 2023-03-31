package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc                           codec.BinaryCodec
		storeKey                      storetypes.StoreKey
		memKey                        storetypes.StoreKey
		paramstore                    paramtypes.Subspace
		whlstTransfersSendersAccAddrs map[string]bool
		router                        *baseapp.MsgServiceRouter
		ak                            types.AccountKeeper
		azk                           types.AuthzKeeper
		tk                            types.TokenKeeper
		nk                            types.NFTKeeper
		ck                            types.CoinFactoryKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	modAccAddrs map[string]bool,
	router *baseapp.MsgServiceRouter,
	ak types.AccountKeeper,
	azk types.AuthzKeeper,
	tk types.TokenKeeper,
	nk types.NFTKeeper,
	ck types.CoinFactoryKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	modAccAddrsCopy := make(map[string]bool)
	for k, v := range modAccAddrs {
		modAccAddrsCopy[k] = v
	}

	return Keeper{
		cdc:                           cdc,
		storeKey:                      storeKey,
		memKey:                        memKey,
		paramstore:                    ps,
		whlstTransfersSendersAccAddrs: modAccAddrsCopy,
		router:                        router,
		ak:                            ak,
		azk:                           azk,
		tk:                            tk,
		nk:                            nk,
		ck:                            ck,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}