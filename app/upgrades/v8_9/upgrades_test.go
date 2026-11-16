package v8_9_test

import (
	"testing"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades/v8_9"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/types"
	"github.com/stretchr/testify/require"
)

func TestICAConstantsMatchIBCGo(t *testing.T) {
	require.Equal(t, icacontrollertypes.StoreKey, v8_9.ICAControllerStoreKey)
	require.Equal(t, icacontrollertypes.SubModuleName, v8_9.ICAControllerStoreKey)
	require.Equal(t, icahosttypes.StoreKey, v8_9.ICAHostStoreKey)
	require.Equal(t, icahosttypes.SubModuleName, v8_9.ICAHostStoreKey)
	require.Equal(t, icatypes.ModuleName, v8_9.ICAModuleName)
	require.ElementsMatch(t, []string{icacontrollertypes.StoreKey, icahosttypes.StoreKey}, v8_9.Upgrade.StoreUpgrades.Deleted)
}

func TestRemoveICAState(t *testing.T) {
	storekeys := map[string]*storetypes.KVStoreKey{
		upgradetypes.StoreKey: storetypes.NewKVStoreKey(upgradetypes.StoreKey),
		paramstypes.StoreKey:  storetypes.NewKVStoreKey(paramstypes.StoreKey),
	}
	ctx := testutil.DefaultContextWithKeys(storekeys, nil, nil)

	versionStore := prefix.NewStore(ctx.KVStore(storekeys[upgradetypes.StoreKey]), []byte{upgradetypes.VersionMapByte})
	versionStore.Set([]byte(icatypes.ModuleName), []byte{0, 0, 0, 0, 0, 0, 0, 3})
	versionStore.Set([]byte("bank"), []byte{0, 0, 0, 0, 0, 0, 0, 4})

	paramsStore := ctx.KVStore(storekeys[paramstypes.StoreKey])
	paramsStore.Set([]byte("icahost/HostEnabled"), []byte("true"))
	paramsStore.Set([]byte("icahost/AllowMessages"), []byte(`["*"]`))
	paramsStore.Set([]byte("icacontroller/ControllerEnabled"), []byte("true"))
	// Subspaces of other modules, including one sharing the "icahost" prefix
	// without the separator, are kept.
	paramsStore.Set([]byte("transfer/SendEnabled"), []byte("true"))
	paramsStore.Set([]byte("icahostx/Key"), []byte("kept"))

	vm := module.VersionMap{icatypes.ModuleName: 3, "bank": 4}

	removed := v8_9.RemoveICAState(ctx, storekeys, vm)

	require.Equal(t, 3, removed)
	require.Equal(t, module.VersionMap{"bank": 4}, vm)
	require.False(t, versionStore.Has([]byte(icatypes.ModuleName)))
	require.True(t, versionStore.Has([]byte("bank")))
	require.False(t, paramsStore.Has([]byte("icahost/HostEnabled")))
	require.False(t, paramsStore.Has([]byte("icahost/AllowMessages")))
	require.False(t, paramsStore.Has([]byte("icacontroller/ControllerEnabled")))
	require.True(t, paramsStore.Has([]byte("transfer/SendEnabled")))
	require.True(t, paramsStore.Has([]byte("icahostx/Key")))

	// Running it again, as on a chain without ICA state, is a no-op.
	require.Equal(t, 0, v8_9.RemoveICAState(ctx, storekeys, vm))
}
