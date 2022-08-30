package keeper

import (
	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (ak Keeper) GetAddress(ctx sdk.Context) (address sdk.AccAddress) {
	return ak.ac.GetModuleAddress(types.ModuleName)
}
