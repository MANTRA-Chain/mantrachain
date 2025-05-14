package simulation

import (
	"math/rand"

	"github.com/MANTRA-Chain/mantrachain/v5/x/sanction/keeper"
	"github.com/MANTRA-Chain/mantrachain/v5/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgAddBlacklistAccounts(
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgAddBlacklistAccounts{
			Authority: k.GetAuthority(),
		}

		// TODO: Handling the AddBlacklistAccounts simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "AddBlacklistAccounts simulation not implemented"), nil, nil
	}
}
