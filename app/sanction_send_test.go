package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sanctiontypes "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func sendCoins(amt int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(FeeDenom, amt)) }

// sanctionSendCtx funds a sanctioned and a clean account.
func sanctionSendCtx(t *testing.T) (*App, sdk.Context, sdk.AccAddress, sdk.AccAddress) {
	t.Helper()
	app, ctx, _, addOnly := sanctionAppCtx(t)
	require.NoError(t, app.BankKeeper.SetParams(ctx, banktypes.DefaultParams()))
	app.AccountKeeper.GetModuleAccount(ctx, minttypes.ModuleName)

	sanctioned := sdk.MustAccAddressFromBech32(sanctionTestAddr(1))
	clean := sdk.MustAccAddressFromBech32(sanctionTestAddr(2))
	require.NoError(t, banktestutil.FundAccount(ctx, app.BankKeeper, sanctioned, sendCoins(1000)))
	require.NoError(t, banktestutil.FundAccount(ctx, app.BankKeeper, clean, sendCoins(1000)))

	msg := sanctiontypes.NewMsgSanctionAddresses(addOnly, []string{sanctioned.String()}, "rule drain-v1")
	_, err := app.MsgServiceRouter().Handler(msg)(ctx, msg)
	require.NoError(t, err)
	return app, ctx, sanctioned, clean
}

// TestSanctionBlocksBankSends checks coins can't leave a sanctioned account by
// any bank path, including messages routed without the ante handler, as ICA
// host txs and contract-dispatched authz are. It can still receive.
func TestSanctionBlocksBankSends(t *testing.T) {
	app, ctx, sanctioned, clean := sanctionSendCtx(t)

	err := app.BankKeeper.SendCoins(ctx, sanctioned, clean, sendCoins(1))
	require.ErrorIs(t, err, sanctiontypes.ErrAccountBlacklisted)

	// Module sends from the account, as fees, IBC escrow and burns use.
	err = app.BankKeeper.SendCoinsFromAccountToModule(ctx, sanctioned, distrtypes.ModuleName, sendCoins(1))
	require.ErrorIs(t, err, sanctiontypes.ErrAccountBlacklisted)

	// MsgMultiSend's input side.
	err = app.BankKeeper.InputOutputCoins(ctx,
		banktypes.NewInput(sanctioned, sendCoins(1)),
		[]banktypes.Output{banktypes.NewOutput(clean, sendCoins(1))})
	require.ErrorIs(t, err, sanctiontypes.ErrAccountBlacklisted)

	// A MsgSend executed straight through the router.
	send := banktypes.NewMsgSend(sanctioned, clean, sendCoins(1))
	_, err = app.MsgServiceRouter().Handler(send)(ctx, send)
	require.ErrorIs(t, err, sanctiontypes.ErrAccountBlacklisted)

	require.Equal(t, sdkmath.NewInt(1000), app.BankKeeper.GetBalance(ctx, sanctioned, FeeDenom).Amount)

	// Receiving is unaffected.
	require.NoError(t, app.BankKeeper.SendCoins(ctx, clean, sanctioned, sendCoins(5)))
	require.Equal(t, sdkmath.NewInt(1005), app.BankKeeper.GetBalance(ctx, sanctioned, FeeDenom).Amount)

	// Clean accounts send as normal.
	require.NoError(t, app.BankKeeper.SendCoins(ctx, clean, sdk.MustAccAddressFromBech32(sanctionTestAddr(3)), sendCoins(1)))
}

// TestUnsanctionRestoresSends checks an account can send again once it's
// taken off the list.
func TestUnsanctionRestoresSends(t *testing.T) {
	app, ctx, sanctioned, clean := sanctionSendCtx(t)
	require.ErrorIs(t, app.BankKeeper.SendCoins(ctx, sanctioned, clean, sendCoins(1)), sanctiontypes.ErrAccountBlacklisted)

	gov := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	msg := sanctiontypes.NewMsgUnsanctionAddresses(gov, []string{sanctioned.String()}, "lifted")
	_, err := app.MsgServiceRouter().Handler(msg)(ctx, msg)
	require.NoError(t, err)

	require.NoError(t, app.BankKeeper.SendCoins(ctx, sanctioned, clean, sendCoins(1)))
}

// TestSanctionBlocksEVMBalanceSpend checks the EVM can't lower a sanctioned
// account's balance either. Value transfers and precompile transfers out of
// an account all end in this burn when the state commits.
func TestSanctionBlocksEVMBalanceSpend(t *testing.T) {
	app, ctx, sanctioned, clean := sanctionSendCtx(t)

	err := app.EVMKeeper.SetBalance(ctx, common.BytesToAddress(sanctioned), uint256.NewInt(400))
	require.ErrorIs(t, err, sanctiontypes.ErrAccountBlacklisted)
	require.Equal(t, sdkmath.NewInt(1000), app.BankKeeper.GetBalance(ctx, sanctioned, FeeDenom).Amount)

	// Raising its balance, as an incoming transfer does, still works.
	require.NoError(t, app.EVMKeeper.SetBalance(ctx, common.BytesToAddress(sanctioned), uint256.NewInt(1500)))
	require.Equal(t, sdkmath.NewInt(1500), app.BankKeeper.GetBalance(ctx, sanctioned, FeeDenom).Amount)

	require.NoError(t, app.EVMKeeper.SetBalance(ctx, common.BytesToAddress(clean), uint256.NewInt(400)))
}

// TestSanctionSendRestrictionChargesNoGas checks the blacklist lookup adds no
// gas to a transfer, whether or not the sender is listed.
func TestSanctionSendRestrictionChargesNoGas(t *testing.T) {
	app, ctx, sanctioned, clean := sanctionSendCtx(t)
	to := sdk.MustAccAddressFromBech32(sanctionTestAddr(3))

	gasCtx := ctx.WithGasMeter(storetypes.NewGasMeter(1_000_000))
	_, err := app.SanctionKeeper.SendRestrictionFn(gasCtx, clean, to, sendCoins(1))
	require.NoError(t, err)
	_, err = app.SanctionKeeper.SendRestrictionFn(gasCtx, sanctioned, to, sendCoins(1))
	require.ErrorIs(t, err, sanctiontypes.ErrAccountBlacklisted)
	require.Zero(t, gasCtx.GasMeter().GasConsumed())
}
