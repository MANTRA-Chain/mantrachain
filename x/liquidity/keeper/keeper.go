package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/MANTRA-Finance/mantrachain/x/liquidity/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		guardKeeper   types.GuardKeeper

		hooks types.LiquidityHooks
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	guardKeeper types.GuardKeeper,
) *Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return &Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		guardKeeper:   guardKeeper,

		hooks: nil,
	}
}

func (k *Keeper) Hooks() types.LiquidityHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiLiquidityHooks{}
	}

	return k.hooks
}

func (k *Keeper) SetHooks(gh types.LiquidityHooks) {
	if k.hooks != nil {
		panic("cannot set liquidity hooks twice")
	}

	k.hooks = gh
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
