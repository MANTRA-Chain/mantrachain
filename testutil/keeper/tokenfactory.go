package keeper

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	"github.com/MANTRA-Chain/mantrachain/v8/x/tokenfactory/keeper"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// MockBankKeeper implements the tokenfactory types.BankKeeper interface for testing
type MockBankKeeper struct {
	HasDenomSupply bool
	DenomExists    bool
}

func (m MockBankKeeper) GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool) {
	return banktypes.Metadata{}, m.DenomExists
}

func (m MockBankKeeper) SetDenomMetaData(ctx context.Context, denomMetaData banktypes.Metadata) {}

func (m MockBankKeeper) HasSupply(ctx context.Context, denom string) bool {
	return m.HasDenomSupply
}

func (m MockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (m MockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return nil
}

func (m MockBankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	return nil
}

func (m MockBankKeeper) BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	return nil
}

func (m MockBankKeeper) SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (m MockBankKeeper) HasBalance(ctx context.Context, addr sdk.AccAddress, amt sdk.Coin) bool {
	return true
}

func (m MockBankKeeper) IsSendEnabledDenom(ctx context.Context, denom string) bool {
	return true
}

// MockERC20Keeper implements the tokenfactory types.ERC20Keeper interface for testing
type MockERC20Keeper struct {
	PrecompilesErr error
	IsRegistered   bool
}

func (m MockERC20Keeper) SetToken(ctx sdk.Context, pair erc20types.TokenPair) error {
	return nil
}

func (m MockERC20Keeper) EnableDynamicPrecompile(ctx sdk.Context, address ethcommon.Address) error {
	return m.PrecompilesErr
}

func (m MockERC20Keeper) IsERC20Registered(ctx sdk.Context, erc20 ethcommon.Address) bool {
	return m.IsRegistered
}

func TokenFactoryKeeper(
	tb testing.TB,
	bk MockBankKeeper,
	ek MockERC20Keeper,
) (keeper.Keeper, sdk.Context) {
	tb.Helper()
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(tb, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		[]string{},
		nil,
		bk,
		nil,
		ek,
		"authority",
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	return k, ctx
}
