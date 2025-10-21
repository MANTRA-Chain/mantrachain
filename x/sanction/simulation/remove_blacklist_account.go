package simulation

import (
	"math/rand"

	"github.com/MANTRA-Chain/mantrachain/v7/x/sanction/keeper"
	"github.com/MANTRA-Chain/mantrachain/v7/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgRemoveBlacklistAccounts(
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgRemoveBlacklistAccounts{
			Authority: k.GetAuthority(),
		}

		// TODO: Handling the RemoveBlacklistAccounts simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "RemoveBlacklistAccounts simulation not implemented"), nil, nil
	}
}
