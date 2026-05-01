package v8_1

import (
	"errors"
	"testing"
	"time"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtestutil "github.com/cosmos/cosmos-sdk/x/distribution/testutil"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var pks = simtestutil.CreateTestPubKeys(3)

func setup(t *testing.T) (distrkeeper.Keeper, *distrtestutil.MockStakingKeeper, sdk.Context) {
	t.Helper()
	ctrl := gomock.NewController(t)

	key := storetypes.NewKVStoreKey(disttypes.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(distribution.AppModuleBasic{})
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: time.Now()})

	bank := distrtestutil.NewMockBankKeeper(ctrl)
	staking := distrtestutil.NewMockStakingKeeper(ctrl)
	account := distrtestutil.NewMockAccountKeeper(ctrl)

	distrAcc := authtypes.NewEmptyModuleAccount(disttypes.ModuleName)
	account.EXPECT().GetModuleAddress(disttypes.ModuleName).Return(distrAcc.GetAddress())
	staking.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec(sdk.Bech32PrefixValAddr)).AnyTimes()
	account.EXPECT().AddressCodec().Return(address.NewBech32Codec(sdk.Bech32MainPrefix)).AnyTimes()

	k := distrkeeper.NewKeeper(
		encCfg.Codec,
		runtime.NewKVStoreService(key),
		account, bank, staking,
		"fee_collector",
		authtypes.NewModuleAddress("gov").String(),
	)
	require.NoError(t, k.FeePool.Set(ctx, disttypes.InitialFeePool()))
	require.NoError(t, k.Params.Set(ctx, disttypes.DefaultParams()))
	return k, staking, ctx
}

func makeDelegation(t *testing.T, idx int, tokens math.Int) (stakingtypes.Validator, sdk.AccAddress, stakingtypes.Delegation) {
	t.Helper()
	val, err := distrtestutil.CreateValidator(pks[idx], tokens)
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32(val.GetOperator())
	require.NoError(t, err)
	delAddr := sdk.AccAddress(valAddr)
	del := stakingtypes.NewDelegation(delAddr.String(), val.GetOperator(), val.DelegatorShares)
	return val, delAddr, del
}

func expectLookup(sk *distrtestutil.MockStakingKeeper, val stakingtypes.Validator, delAddr sdk.AccAddress, del stakingtypes.Delegation) {
	valAddr, _ := sdk.ValAddressFromBech32(val.GetOperator())
	sk.EXPECT().Validator(gomock.Any(), valAddr).Return(val, nil)
	sk.EXPECT().Delegation(gomock.Any(), delAddr, valAddr).Return(del, nil)
}

func setStartingStake(t *testing.T, k distrkeeper.Keeper, ctx sdk.Context, val stakingtypes.Validator, delAddr sdk.AccAddress, stake math.LegacyDec) {
	t.Helper()
	valAddr, err := sdk.ValAddressFromBech32(val.GetOperator())
	require.NoError(t, err)
	info := disttypes.NewDelegatorStartingInfo(0, stake, uint64(ctx.BlockHeight()))
	require.NoError(t, k.SetDelegatorStartingInfo(ctx, valAddr, delAddr, info))
}

func getStartingStake(t *testing.T, k distrkeeper.Keeper, ctx sdk.Context, val stakingtypes.Validator, delAddr sdk.AccAddress) math.LegacyDec {
	t.Helper()
	valAddr, err := sdk.ValAddressFromBech32(val.GetOperator())
	require.NoError(t, err)
	info, err := k.GetDelegatorStartingInfo(ctx, valAddr, delAddr)
	require.NoError(t, err)
	return info.Stake
}

// info.Stake > currentStake is clamped; <= currentStake is left alone.
func TestFixSilentlySkippedSlashes_ClampsOversizedStake(t *testing.T) {
	k, sk, ctx := setup(t)

	// pair 0: silent slash — Tokens dropped, Shares unchanged.
	val0, del0Addr, del0 := makeDelegation(t, 0, math.NewInt(1000))
	val0.Tokens = math.NewInt(900)
	cur0 := val0.TokensFromShares(del0.Shares)
	setStartingStake(t, k, ctx, val0, del0Addr, cur0.Add(math.LegacyNewDec(50)))

	// pair 1: info.Stake == currentStake.
	val1, del1Addr, del1 := makeDelegation(t, 1, math.NewInt(2000))
	cur1 := val1.TokensFromShares(del1.Shares)
	setStartingStake(t, k, ctx, val1, del1Addr, cur1)

	// pair 2: info.Stake < currentStake (typical post-rewards state).
	val2, del2Addr, del2 := makeDelegation(t, 2, math.NewInt(3000))
	below2 := val2.TokensFromShares(del2.Shares).Sub(math.LegacyNewDec(1))
	setStartingStake(t, k, ctx, val2, del2Addr, below2)

	expectLookup(sk, val0, del0Addr, del0)
	expectLookup(sk, val1, del1Addr, del1)
	expectLookup(sk, val2, del2Addr, del2)

	require.NoError(t, fixSilentlySkippedSlashes(ctx, sk, k))

	require.True(t, getStartingStake(t, k, ctx, val0, del0Addr).Equal(cur0))
	require.True(t, getStartingStake(t, k, ctx, val1, del1Addr).Equal(cur1))
	require.True(t, getStartingStake(t, k, ctx, val2, del2Addr).Equal(below2))
}

// empty store -> no-op, no error.
func TestFixSilentlySkippedSlashes_NoEntries(t *testing.T) {
	k, sk, ctx := setup(t)
	require.NoError(t, fixSilentlySkippedSlashes(ctx, sk, k))
}

func setSlashEvent(t *testing.T, k distrkeeper.Keeper, ctx sdk.Context, val stakingtypes.Validator, fraction string) {
	t.Helper()
	valAddr, err := sdk.ValAddressFromBech32(val.GetOperator())
	require.NoError(t, err)
	require.NoError(t, k.SetValidatorSlashEvent(ctx, valAddr,
		uint64(ctx.BlockHeight()), 1,
		disttypes.NewValidatorSlashEvent(1, math.LegacyMustNewDecFromStr(fraction))))
}

// Recorded slash fully explains info.Stake > currentStake → clamp is a no-op.
func TestFixSilentlySkippedSlashes_RecordedSlashLeftAlone(t *testing.T) {
	k, sk, ctx := setup(t)

	val, delAddr, del := makeDelegation(t, 0, math.NewInt(1000))
	preStake := val.TokensFromShares(del.Shares)
	val.Tokens = math.NewInt(950) // 5% recorded slash

	setStartingStake(t, k, ctx, val, delAddr, preStake)
	setSlashEvent(t, k, ctx, val, "0.05")
	expectLookup(sk, val, delAddr, del)

	require.NoError(t, fixSilentlySkippedSlashes(ctx, sk, k))
	require.True(t, getStartingStake(t, k, ctx, val, delAddr).Equal(preStake))
}

// Recorded 5% + silent 1%: clamp by silent residual only.
func TestFixSilentlySkippedSlashes_PartialSilentSkipResidue(t *testing.T) {
	k, sk, ctx := setup(t)

	val, delAddr, del := makeDelegation(t, 0, math.NewInt(1000))
	preStake := val.TokensFromShares(del.Shares)
	val.Tokens = math.NewInt(940) // 5% recorded + ~1% silent
	cur := val.TokensFromShares(del.Shares)

	setStartingStake(t, k, ctx, val, delAddr, preStake)
	setSlashEvent(t, k, ctx, val, "0.05")
	expectLookup(sk, val, delAddr, del)

	require.NoError(t, fixSilentlySkippedSlashes(ctx, sk, k))

	expected := preStake.MulTruncate(math.LegacyMustNewDecFromStr("0.95"))
	want := preStake.MulTruncate(cur.QuoTruncate(expected))
	require.True(t, getStartingStake(t, k, ctx, val, delAddr).Equal(want))
}

// staking read errors must propagate, not be swallowed.
func TestFixSilentlySkippedSlashes_PropagatesValidatorErr(t *testing.T) {
	k, sk, ctx := setup(t)

	val, delAddr, _ := makeDelegation(t, 0, math.NewInt(1000))
	seeded := math.LegacyNewDec(1000)
	setStartingStake(t, k, ctx, val, delAddr, seeded)

	valAddr, _ := sdk.ValAddressFromBech32(val.GetOperator())
	bang := errors.New("staking store corrupt")
	sk.EXPECT().Validator(gomock.Any(), valAddr).Return(nil, bang)

	require.ErrorIs(t, fixSilentlySkippedSlashes(ctx, sk, k), bang)
	require.True(t, getStartingStake(t, k, ctx, val, delAddr).Equal(seeded))
}
