package keeper_test

import (
	"testing"
	"time"

	gogoproto "github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
	protov2 "google.golang.org/protobuf/proto"

	txsigning "cosmossdk.io/x/tx/signing"

	keepertest "github.com/MANTRA-Chain/mantrachain/v8/testutil/keeper"
	sanctionkeeper "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktxsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// mockTx implements authsigning.Tx with stub methods for fields unused by the decorator.
type mockTx struct {
	signers    [][]byte
	feeGranter []byte
	msgs       []sdk.Msg
}

func (m mockTx) GetMsgs() []sdk.Msg                                   { return m.msgs }
func (m mockTx) GetMsgsV2() ([]protov2.Message, error)                { return nil, nil }
func (m mockTx) GetSigners() ([][]byte, error)                        { return m.signers, nil }
func (m mockTx) GetPubKeys() ([]cryptotypes.PubKey, error)            { return nil, nil }
func (m mockTx) GetSignaturesV2() ([]sdktxsigning.SignatureV2, error) { return nil, nil }
func (m mockTx) GetMemo() string                                      { return "" }
func (m mockTx) GetGas() uint64                                       { return 0 }
func (m mockTx) GetFee() sdk.Coins                                    { return nil }
func (m mockTx) FeePayer() []byte                                     { return nil }
func (m mockTx) FeeGranter() []byte                                   { return m.feeGranter }
func (m mockTx) GetTimeoutHeight() uint64                             { return 0 }
func (m mockTx) GetTimeoutTimeStamp() time.Time                       { return time.Time{} }
func (m mockTx) GetUnordered() bool                                   { return false }
func (m mockTx) ValidateBasic() error                                 { return nil }

// newTestCodec builds a ProtoCodec whose signing context can resolve MsgSend signers.
func newTestCodec() codec.Codec {
	addrCodec := addresscodec.NewBech32Codec("mantra")
	registry, err := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: gogoproto.HybridResolver,
		SigningOptions: txsigning.Options{
			AddressCodec:          addrCodec,
			ValidatorAddressCodec: addresscodec.NewBech32Codec("mantravaloper"),
		},
	})
	if err != nil {
		panic(err)
	}
	banktypes.RegisterInterfaces(registry)
	authz.RegisterInterfaces(registry)
	return codec.NewProtoCodec(registry)
}

func passThrough(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) { return ctx, nil }

func TestBlacklistCheckDecorator_AnteHandle(t *testing.T) {
	k, ctx, _ := keepertest.SanctionKeeper(t)
	cdc := newTestCodec()
	dec := sanctionkeeper.NewBlacklistCheckDecorator(k, cdc)

	clean := sdk.AccAddress([]byte("clean_______________"))
	bad := sdk.AccAddress([]byte("blacklisted_________"))
	other := sdk.AccAddress([]byte("other_______________"))

	require.NoError(t, k.BlacklistAccounts.Set(ctx, bad.String()))

	msgSendFrom := func(from sdk.AccAddress) *banktypes.MsgSend {
		return banktypes.NewMsgSend(from, other, sdk.NewCoins(sdk.NewInt64Coin("amantra", 1)))
	}

	tests := []struct {
		name    string
		tx      mockTx
		wantErr bool
	}{
		{
			name:    "clean signer passes",
			tx:      mockTx{signers: [][]byte{clean}},
			wantErr: false,
		},
		{
			name:    "blacklisted signer blocked",
			tx:      mockTx{signers: [][]byte{bad}},
			wantErr: true,
		},
		{
			name:    "blacklisted fee granter blocked",
			tx:      mockTx{signers: [][]byte{clean}, feeGranter: bad},
			wantErr: true,
		},
		{
			name:    "clean fee granter passes",
			tx:      mockTx{signers: [][]byte{clean}, feeGranter: clean},
			wantErr: false,
		},
		{
			name: "blacklisted authz granter blocked (single level)",
			tx: mockTx{
				signers: [][]byte{clean},
				msgs: []sdk.Msg{
					func() sdk.Msg {
						exec := authz.NewMsgExec(clean, []sdk.Msg{msgSendFrom(bad)})
						return &exec
					}(),
				},
			},
			wantErr: true,
		},
		{
			name: "clean authz granter passes (single level)",
			tx: mockTx{
				signers: [][]byte{clean},
				msgs: []sdk.Msg{
					func() sdk.Msg {
						exec := authz.NewMsgExec(clean, []sdk.Msg{msgSendFrom(clean)})
						return &exec
					}(),
				},
			},
			wantErr: false,
		},
		{
			// Nested MsgExec is rejected regardless of whether the inner granter
			// is blacklisted — it cannot be safely inspected with a flat check.
			name: "nested MsgExec rejected (blacklisted inner granter)",
			tx: mockTx{
				signers: [][]byte{clean},
				msgs: []sdk.Msg{
					func() sdk.Msg {
						innerExec := authz.NewMsgExec(other, []sdk.Msg{msgSendFrom(bad)})
						outerExec := authz.NewMsgExec(clean, []sdk.Msg{&innerExec})
						return &outerExec
					}(),
				},
			},
			wantErr: true,
		},
		{
			// Even all-clean accounts cannot use nested MsgExec.
			name: "nested MsgExec rejected (all-clean accounts)",
			tx: mockTx{
				signers: [][]byte{clean},
				msgs: []sdk.Msg{
					func() sdk.Msg {
						innerExec := authz.NewMsgExec(other, []sdk.Msg{msgSendFrom(clean)})
						outerExec := authz.NewMsgExec(clean, []sdk.Msg{&innerExec})
						return &outerExec
					}(),
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := dec.AnteHandle(ctx, tc.tx, false, passThrough)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
