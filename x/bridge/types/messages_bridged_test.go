package types

import (
	"testing"

	"github.com/MANTRA-Finance/mantrachain/testutil/sample"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateMultiBridged_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateMultiBridged
		err  error
	}{
		{
			name: "no inputs",
			msg: MsgCreateMultiBridged{
				Inputs: []Input{},
			},
			err: ErrNoInputs,
		}, {
			name: "valid address",
			msg: MsgCreateMultiBridged{
				Inputs: []Input{
					{
						Address: sample.AccAddress(),
						Coins:   sdk.NewCoins(sdk.NewCoin("uom", sdk.NewInt(100))),
					},
				},
				Outputs: []Output{
					{
						Address: sample.AccAddress(),
						Coins:   sdk.NewCoins(sdk.NewCoin("uom", sdk.NewInt(100))),
					},
				},
				EthTxHashes: []string{"0x1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
