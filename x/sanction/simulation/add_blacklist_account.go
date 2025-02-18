package simulation

import (
	"math/rand"

	"github.com/MANTRA-Chain/mantrachain/v3/x/sanction/keeper"
	"github.com/MANTRA-Chain/mantrachain/v3/x/sanction/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgAddBlacklistAccount(
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := &types.MsgAddBlacklistAccount{
			Authority: k.GetAuthority(),
		}

		// TODO: Handling the AddBlacklistAccount simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "AddBlacklistAccount simulation not implemented"), nil, nil
	}
}
