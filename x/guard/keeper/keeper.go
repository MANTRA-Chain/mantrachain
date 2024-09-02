package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/MANTRA-Finance/mantrachain/x/guard/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
		router    *baseapp.MsgServiceRouter

		authzKeeper       types.AuthzKeeper
		nftKeeper         types.NFTKeeper
		tokenKeeper       types.TokenKeeper
		coinFactoryKeeper types.CoinFactoryKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	router *baseapp.MsgServiceRouter,

	authzKeeper types.AuthzKeeper,
	nftKeeper types.NFTKeeper,
) *Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return &Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
		router:       router,

		authzKeeper: authzKeeper,
		nftKeeper:   nftKeeper,
	}
}

func SetSendRestrictions(k *Keeper, bankKeeper types.BankKeeper) {
	bankKeeper.AppendSendRestriction(func(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		if err := k.ValidateCoinsTransfers(
			sdkCtx,
			[]banktypes.Input{{Address: fromAddr.String(), Coins: amt}},
			[]banktypes.Output{{Address: toAddr.String(), Coins: amt}},
		); err != nil {
			return nil, err
		}
		return toAddr, nil
	})
}

func SetCoinFactoryKeeper(k *Keeper, coinFactoryKeeper types.CoinFactoryKeeper) {
	k.coinFactoryKeeper = coinFactoryKeeper
}

func SetTokenKeeper(k *Keeper, tokenKeeper types.TokenKeeper) {
	k.tokenKeeper = tokenKeeper
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
