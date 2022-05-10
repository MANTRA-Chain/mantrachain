package simulation

import (
	"math/rand"

	"github.com/LimeChain/mantrachain/x/mns/keeper"
	"github.com/LimeChain/mantrachain/x/mns/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgCreateDomainName(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgCreateDomainName{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the CreateDomainName simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "CreateDomainName simulation not implemented"), nil, nil
	}
}
