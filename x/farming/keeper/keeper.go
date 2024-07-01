package keeper

import (
	"fmt"
	"strconv"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MANTRA-Finance/mantrachain/x/farming/types"
)

var (
	enableAdvanceEpoch = "false" // Set this to "true" using build flags to enable AdvanceEpoch msg handling.
	enableRatioPlan    = "false" // Set this to "true" using build flags to enable creation of RatioPlans.

	// EnableAdvanceEpoch indicates whether msgServer accepts MsgAdvanceEpoch or not.
	// Never set this to true in production mode. Doing that will expose serious attack vector.
	EnableAdvanceEpoch = false
	// EnableRatioPlan indicates whether msgServer and proposal handler accept
	// creation of RatioPlans.
	// Default is false, which means that RatioPlans can't be created through a
	// MsgCreateRatioPlan msg and a PublicPlanProposal.
	EnableRatioPlan = false
)

// TODO: test if it is executed
func init() {
	var err error
	EnableAdvanceEpoch, err = strconv.ParseBool(enableAdvanceEpoch)
	if err != nil {
		panic(err)
	}
	EnableRatioPlan, err = strconv.ParseBool(enableRatioPlan)
	if err != nil {
		panic(err)
	}
}

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
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	// ensure farming module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	guardKeeper.AddTransferAccAddressesWhitelist([]string{
		types.DefaultFarmingFeeCollector.String(),
		types.RewardsReserveAcc.String(),
		types.UnharvestedRewardsReserveAcc.String(),
	})

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		guardKeeper:   guardKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAccountKeeper() types.AccountKeeper {
	return k.accountKeeper
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}
