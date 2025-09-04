package v5rc5_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades"
	"github.com/MANTRA-Chain/mantrachain/v5/app/upgrades/v5rc5"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/stretchr/testify/require"
)

func TestCreateUpgradeHandler(t *testing.T) {
	storeKey := storetypes.NewKVStoreKey(erc20types.StoreKey)
	codec := app.MakeEncodingConfig(t).Codec
	a := app.MakeTestApp(t)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	ac := app.AddressCodec{}
	erc20Keeper := erc20keeper.NewKeeper(storeKey, codec, authority, a.AccountKeeper, a.BankKeeper, a.EVMKeeper, a.StakingKeeper, &a.TransferKeeper, ac)
	ctx := testutil.DefaultContext(storeKey, storetypes.NewTransientStoreKey("transient"))
	store := ctx.KVStore(storeKey)
	dynamicPrecompile := "0x6eC942095eCD4948d9C094337ABd59Dc3c521005"
	nativePrecompile := "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
	store.Set([]byte("DynamicPrecompiles"), []byte(dynamicPrecompile))
	store.Set([]byte("NativePrecompiles"), []byte(nativePrecompile))
	configurator := module.NewConfigurator(nil, nil, nil)
	keepers := &upgrades.UpgradeKeepers{Erc20Keeper: erc20Keeper}
	storekeys := map[string]*storetypes.KVStoreKey{
		erc20types.StoreKey: storeKey,
	}
	handler := v5rc5.CreateUpgradeHandler(&module.Manager{}, configurator, keepers, storekeys)
	_, err := handler(ctx, upgradetypes.Plan{}, module.VersionMap{})
	require.NoError(t, err)
	dynamicPrecompiles := erc20Keeper.GetDynamicPrecompiles(ctx)
	nativePrecompiles := erc20Keeper.GetNativePrecompiles(ctx)
	require.Len(t, dynamicPrecompiles, 1)
	require.Len(t, nativePrecompiles, 1)
	require.Equal(t, dynamicPrecompile, dynamicPrecompiles[0])
	require.Equal(t, nativePrecompile, nativePrecompiles[0])
}
