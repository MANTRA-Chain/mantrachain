package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/LimeChain/mantrachain/x/token/types"
)

func (k Keeper) CheckCanTransferNft(ctx sdk.Context, classId string) (bool, error) {
	return false, errors.Wrapf(types.ErrNftModuleTransferNftDisabled, "nft module transfer nft disabled")
}
