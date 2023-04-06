package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/MANTRA-Finance/mantrachain/x/token/types"
)

func (k Keeper) CheckCanTransferNft(ctx sdk.Context, classId string) (bool, error) {
	return false, errors.Wrapf(types.ErrNftModuleTransferNftDisabled, "nft module transfer nft disabled")
}
