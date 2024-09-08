package simulation

import (
	"math/rand"

	"github.com/MANTRA-Chain/mantrachain/x/xfeemarket/keeper"
	"github.com/MANTRA-Chain/mantrachain/x/xfeemarket/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgUpsertFeeDenom(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgUpsertFeeDenom{
			Authority: k.GetAuthority(),
		}

		// TODO: Handling the UpsertFeeDenom simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "UpsertFeeDenom simulation not implemented"), nil, nil
	}
}
