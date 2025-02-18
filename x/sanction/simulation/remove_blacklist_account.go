package simulation

import (
	"math/rand"

	"github.com/MANTRA-Chain/mantrachain/v3/x/sanction/keeper"
	"github.com/MANTRA-Chain/mantrachain/v3/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgRemoveBlacklistAccount(
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgRemoveBlacklistAccount{
			Authority: k.GetAuthority(),
		}

		// TODO: Handling the RemoveBlacklistAccount simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "RemoveBlacklistAccount simulation not implemented"), nil, nil
	}
}
