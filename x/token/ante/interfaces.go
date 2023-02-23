package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TokenKeeper interface {
	CheckCanTransfer(ctx sdk.Context, tokenKeeper TokenKeeper, classId string) (bool, error)
}
