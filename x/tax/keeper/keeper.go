package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/MANTRA-Chain/mantrachain/x/tax/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		addressCodec address.Codec
		storeService store.KVStoreService
		logger       log.Logger
		authKeeper   types.AccountKeeper
		bankKeeper   types.BankKeeper

		// the address capable of executing a MsgUpdateParams message.
		// Typically, this should be the x/gov module account.
		authority string

		Schema collections.Schema
		Params collections.Item[types.Params]

		feeCollectorName string // name of the FeeCollector ModuleAccount
		// this line is used by starport scaffolding # collection/type
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	feeCollectorName string,
) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:              cdc,
		addressCodec:     addressCodec,
		storeService:     storeService,
		authority:        authority,
		logger:           logger,
		authKeeper:       ak,
		bankKeeper:       bk,
		feeCollectorName: feeCollectorName,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		// this line is used by starport scaffolding # collection/instantiate
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) AllocateMcaTax(ctx context.Context, mcaTax math.LegacyDec, mcaAddress sdk.AccAddress) error {
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	mcaTaxAllocation, _ := feesCollected.MulDec(mcaTax).TruncateDecimal()

	// transfer allocated mca tax to the specified account
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, mcaAddress, mcaTaxAllocation)
	if err != nil {
		return err
	}
	// emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMcaTaxAmount,
			sdk.NewAttribute(sdk.AttributeKeyAmount, mcaTaxAllocation.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, mcaAddress.String()),
		),
	)
	return nil
}
