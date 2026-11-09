package app

import (
	"fmt"
	"strings"
	"testing"

	storetypes "cosmossdk.io/store/types"
	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/keeper"
	sanctiontypes "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
)

func sanctionTestAddr(i int) string {
	return sdk.AccAddress([]byte(fmt.Sprintf("sanction-target-%04d", i))).String()
}

// sanctionAppCtx returns an app with one full and one add-only sanctioner.
func sanctionAppCtx(t *testing.T) (*App, sdk.Context, string, string) {
	t.Helper()
	// Setup's full genesis currently fails in x/vm (missing denom metadata),
	// so start from an empty store and initialise only what these tests use.
	app := SetupWithEmptyStore(t)
	ctx := app.BaseApp.NewUncachedContext(false, cmtproto.Header{Height: 1})
	require.NoError(t, app.SanctionKeeper.Params.Set(ctx, sanctiontypes.DefaultParams()))
	for _, name := range []string{govtypes.ModuleName, distrtypes.ModuleName, authtypes.FeeCollectorName} {
		app.AccountKeeper.GetModuleAccount(ctx, name) // creates it in x/auth
	}

	full := sdk.AccAddress([]byte("sanctioner-full_____")).String()
	addOnly := sdk.AccAddress([]byte("sanctioner-addonly__")).String()
	msg := sanctiontypes.NewMsgUpdateSanctioners(authtypes.NewModuleAddress(govtypes.ModuleName).String(), []sanctiontypes.Sanctioner{
		{Address: full, Role: sanctiontypes.SanctionerRole_SANCTIONER_ROLE_FULL},
		{Address: addOnly, Role: sanctiontypes.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY},
	})
	_, err := app.MsgServiceRouter().Handler(msg)(ctx, msg)
	require.NoError(t, err)
	return app, ctx, full, addOnly
}

// TestSanctionFastPathBlocksSigner routes a sanctioner message through the
// app's message router and checks the ante decorator then rejects the target.
func TestSanctionFastPathBlocksSigner(t *testing.T) {
	app, ctx, full, addOnly := sanctionAppCtx(t)
	attacker := sanctionTestAddr(1)

	buildSend := func(from string) authsigning.Tx {
		send := banktypes.NewMsgSend(sdk.MustAccAddressFromBech32(from), sdk.MustAccAddressFromBech32(full), sdk.NewCoins(sdk.NewInt64Coin("amantra", 1)))
		b := app.TxConfig().NewTxBuilder()
		require.NoError(t, b.SetMsgs(send))
		return b.GetTx()
	}
	dec := sanctionkeeper.NewBlacklistCheckDecorator(app.SanctionKeeper, app.AppCodec())
	next := func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) { return ctx, nil }

	_, err := dec.AnteHandle(ctx, buildSend(attacker), false, next)
	require.NoError(t, err)

	msg := sanctiontypes.NewMsgSanctionAddresses(addOnly, []string{attacker}, "rule drain-v1, tx 0xabc")
	_, err = app.MsgServiceRouter().Handler(msg)(ctx, msg)
	require.NoError(t, err)

	_, err = dec.AnteHandle(ctx, buildSend(attacker), false, next)
	require.ErrorIs(t, err, sanctiontypes.ErrAccountBlacklisted)

	unmsg := sanctiontypes.NewMsgUnsanctionAddresses(full, []string{attacker}, "reviewed")
	_, err = app.MsgServiceRouter().Handler(unmsg)(ctx, unmsg)
	require.NoError(t, err)
	_, err = dec.AnteHandle(ctx, buildSend(attacker), false, next)
	require.NoError(t, err)
}

// TestSanctionFastPathProtectsRealModuleAccounts uses the real account keeper
// and module list.
func TestSanctionFastPathProtectsRealModuleAccounts(t *testing.T) {
	app, ctx, full, _ := sanctionAppCtx(t)
	for _, name := range []string{govtypes.ModuleName, distrtypes.ModuleName, authtypes.FeeCollectorName, sanctiontypes.ModuleName} {
		target := authtypes.NewModuleAddress(name).String()
		msg := sanctiontypes.NewMsgSanctionAddresses(full, []string{target}, "should fail")
		_, err := app.MsgServiceRouter().Handler(msg)(ctx, msg)
		require.ErrorIs(t, err, sanctiontypes.ErrProtectedAddress, name)
	}
}

// TestSanctionGasAtCap measures execution gas for MsgSanctionAddresses and
// estimates tx-size gas, at mainnet's tx_size_cost_per_byte of 10.
func TestSanctionGasAtCap(t *testing.T) {
	app, ctx, _, addOnly := sanctionAppCtx(t)
	reason := strings.Repeat("r", int(sanctiontypes.DefaultMaxReasonLength))
	const txSizeCostPerByte = 10
	// A signed tx adds roughly this many bytes for the signature, pubkey,
	// fee and sequence beyond the unsigned body.
	const signedOverheadBytes = 150

	next := 0
	var atCap uint64
	for _, n := range []int{1, 10, int(sanctiontypes.DefaultMaxAddressesPerMsg)} {
		addrs := make([]string, n)
		for i := range addrs {
			addrs[i] = sanctionTestAddr(next)
			next++
		}
		msg := sanctiontypes.NewMsgSanctionAddresses(addOnly, addrs, reason)

		b := app.TxConfig().NewTxBuilder()
		require.NoError(t, b.SetMsgs(msg))
		txBytes, err := app.TxConfig().TxEncoder()(b.GetTx())
		require.NoError(t, err)
		sizeGas := uint64(len(txBytes)+signedOverheadBytes) * txSizeCostPerByte

		gctx := ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
		_, err = app.MsgServiceRouter().Handler(msg)(gctx, msg)
		require.NoError(t, err)
		execGas := gctx.GasMeter().GasConsumed()

		t.Logf("addresses=%d exec_gas=%d (%d/address) tx_bytes~%d size_gas~%d total~%d",
			n, execGas, execGas/uint64(n), len(txBytes)+signedOverheadBytes, sizeGas, execGas+sizeGas)
		if n == int(sanctiontypes.DefaultMaxAddressesPerMsg) {
			atCap = execGas + sizeGas
		}
	}
	// Mainnet block max_gas is 75,000,000. A batch at the cap must stay a
	// small fraction of a block.
	require.Less(t, atCap, uint64(2_000_000))
}
