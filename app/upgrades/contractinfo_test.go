package upgrades_test

import (
	"encoding/hex"
	"testing"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/MANTRA-Chain/mantrachain/v8/app/upgrades"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

// ContractInfo of a mantra-1 contract written by a retracted wasmd.
const staleHex = "085212416d616e74726131303076666a79356b3538716437776478686379746c38666772726a71656d6a75666536347a6b3079787a66796b34797232356e736371347671711a2d6d616e7472613166353435366d3234707a6d66683234306e633637637377327830706470777264776d7735366e2215636c61696d64726f702d63616d706169676e2d31372a0508c898b0063a407761736d323130683072676a78646463733436656c326b7579616e7a70776635767339756563646e656478706a7374336a653264757333737a736833326c7466"

const creator = "mantra1f5456m24pzmfh240nc67csw2x0pdpwrdwmw56n"

func contractStore(t *testing.T) (sdk.Context, *storetypes.KVStoreKey, prefix.Store) {
	t.Helper()
	key := storetypes.NewKVStoreKey(wasmtypes.StoreKey)
	ctx := testutil.DefaultContext(key, storetypes.NewTransientStoreKey("transient_wasm"))
	return ctx, key, prefix.NewStore(ctx.KVStore(key), wasmtypes.ContractKeyPrefix)
}

func decode(t *testing.T, store prefix.Store, key string) wasmtypes.ContractInfo {
	t.Helper()
	var info wasmtypes.ContractInfo
	require.NoError(t, info.Unmarshal(store.Get([]byte(key))))
	return info
}

// oldLayout encodes ContractInfo as the retracted wasmd did: port at 7, extension at 8.
func oldLayout(port string, extension []byte) []byte {
	var bz []byte
	bz = protowire.AppendTag(bz, 1, protowire.VarintType)
	bz = protowire.AppendVarint(bz, 1)
	bz = protowire.AppendTag(bz, 2, protowire.BytesType)
	bz = protowire.AppendString(bz, creator)
	bz = protowire.AppendTag(bz, 4, protowire.BytesType)
	bz = protowire.AppendString(bz, "old")
	bz = protowire.AppendTag(bz, 5, protowire.BytesType)
	bz = protowire.AppendBytes(bz, []byte{0x08, 0x04})
	bz = protowire.AppendTag(bz, 7, protowire.BytesType)
	bz = protowire.AppendString(bz, port)
	if extension != nil {
		bz = protowire.AppendTag(bz, 8, protowire.BytesType)
		bz = protowire.AppendBytes(bz, extension)
	}
	return bz
}

func TestMigrateContractInfo(t *testing.T) {
	stale, err := hex.DecodeString(staleHex)
	require.NoError(t, err)
	require.ErrorContains(t, new(wasmtypes.ContractInfo).Unmarshal(stale), "illegal wireType 7")
	healthy, err := (&wasmtypes.ContractInfo{CodeID: 7, Creator: creator, Label: "healthy", IBC2PortID: "wasm2healthy"}).Marshal()
	require.NoError(t, err)

	ctx, key, store := contractStore(t)
	store.Set([]byte("stale"), stale)
	store.Set([]byte("healthy"), healthy)

	n, err := upgrades.MigrateContractInfo(ctx, key)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	want := wasmtypes.ContractInfo{
		CodeID:     82,
		Creator:    "mantra100vfjy5k58qd7wdxhcytl8fgrrjqemjufe64zk0yxzfyk4yr25nscq4vqq",
		Admin:      creator,
		Label:      "claimdrop-campaign-17",
		Created:    &wasmtypes.AbsoluteTxPosition{BlockHeight: 13372488},
		IBC2PortID: "wasm210h0rgjxddcs46el2kuyanzpwf5vs9uecdnedxpjst3je2dus3szsh32ltf",
	}
	require.Equal(t, want, decode(t, store, "stale"))
	// byte-identical to what this wasmd writes natively
	wantBz, err := want.Marshal()
	require.NoError(t, err)
	require.Equal(t, wantBz, store.Get([]byte("stale")))
	require.Equal(t, healthy, store.Get([]byte("healthy")))

	// idempotent
	n, err = upgrades.MigrateContractInfo(ctx, key)
	require.NoError(t, err)
	require.Zero(t, n)
}

func TestMigrateContractInfoMovesExtension(t *testing.T) {
	ext := &codectypes.Any{TypeUrl: "/test.Extension", Value: []byte("payload")}
	extBz, err := ext.Marshal()
	require.NoError(t, err)

	ctx, key, store := contractStore(t)
	store.Set([]byte("c"), oldLayout("wasm2ext", extBz))

	n, err := upgrades.MigrateContractInfo(ctx, key)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	got := decode(t, store, "c")
	require.Equal(t, "wasm2ext", got.IBC2PortID)
	require.Equal(t, ext, got.Extension)
}

func TestMigrateContractInfoAborts(t *testing.T) {
	for name, bad := range map[string][]byte{
		"truncated":  {0x0a},
		"not a port": oldLayout("not-a-port", nil),
	} {
		t.Run(name, func(t *testing.T) {
			repairable := oldLayout("wasm2fine", nil)
			ctx, key, store := contractStore(t)
			store.Set([]byte("a"), repairable)
			store.Set([]byte("b"), bad)

			n, err := upgrades.MigrateContractInfo(ctx, key)
			require.Error(t, err)
			require.Zero(t, n)
			// nothing written, not even the repairable sibling
			require.Equal(t, repairable, store.Get([]byte("a")))
			require.Equal(t, bad, store.Get([]byte("b")))
		})
	}
}
