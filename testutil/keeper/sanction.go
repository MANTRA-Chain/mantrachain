package keeper

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/keeper"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
)

// SanctionModuleAccounts are the module names the test account keeper
// reports as module accounts.
var SanctionModuleAccounts = []string{govtypes.ModuleName, distrtypes.ModuleName, authtypes.FeeCollectorName, types.ModuleName}

// mockAccountKeeper reports SanctionModuleAccounts as module accounts and
// knows no other accounts.
type mockAccountKeeper struct{}

func (mockAccountKeeper) GetAccount(context.Context, sdk.AccAddress) sdk.AccountI { return nil }

func (mockAccountKeeper) GetModulePermissions() map[string]authtypes.PermissionsForAddress {
	perms := make(map[string]authtypes.PermissionsForAddress, len(SanctionModuleAccounts))
	for _, name := range SanctionModuleAccounts {
		perms[name] = authtypes.NewPermissionsForAddress(name, nil)
	}
	return perms
}

// SanctionEscrowAccount is the address the test tokenfactory keeper reports
// as an IBC escrow account.
var SanctionEscrowAccount = sdk.AccAddress([]byte("escrow______________"))

type mockTokenFactoryKeeper struct{}

func (mockTokenFactoryKeeper) IsEscrowAddress(_ sdk.Context, addr sdk.AccAddress) bool {
	return addr.Equals(SanctionEscrowAccount)
}

func SanctionKeeper(tb testing.TB) (keeper.Keeper, sdk.Context, address.Codec) {
	tb.Helper()
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(tb, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		mockAccountKeeper{},
		nil,
		mockTokenFactoryKeeper{},
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		tb.Fatalf("failed to set params: %v", err)
	}

	return k, ctx, addressCodec
}
