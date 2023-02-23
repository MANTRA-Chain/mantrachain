package keeper

import (
	"cosmossdk.io/errors"
	ante "github.com/LimeChain/mantrachain/x/token/ante"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/LimeChain/mantrachain/x/token/types"
)

func (k Keeper) CheckCanTransfer(ctx sdk.Context, _ ante.TokenKeeper, classId string) (bool, error) {
	// TODO: disable transfer only for soul bond nfts
	return false, errors.Wrapf(types.ErrNftModuleTransferNftDisabled, "nft module transfer nft disabled")
}
