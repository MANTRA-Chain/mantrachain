package cli_test

import (
	"fmt"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/status"

	"github.com/LimeChain/mantrachain/testutil/network"
	"github.com/LimeChain/mantrachain/testutil/nullify"
	"github.com/LimeChain/mantrachain/x/bridge/client/cli"
    "github.com/LimeChain/mantrachain/x/bridge/types"
)

func networkWithCw20ContractObjects(t *testing.T) (*network.Network, types.Cw20Contract) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
    require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	cw20Contract := &types.Cw20Contract{}
	nullify.Fill(&cw20Contract)
	state.Cw20Contract = cw20Contract
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), *state.Cw20Contract
}

func TestShowCw20Contract(t *testing.T) {
	net, obj := networkWithCw20ContractObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		args []string
		err  error
		obj  types.Cw20Contract
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowCw20Contract(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetCw20ContractResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Cw20Contract)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.Cw20Contract),
				)
			}
		})
	}
}

