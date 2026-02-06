package distrclaim

import (
	"context"

	cmn "github.com/cosmos/evm/precompiles/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegatorWithdrawAddrGetter is the minimal interface needed to read a delegator's withdraw address
// from the distribution module.
type DelegatorWithdrawAddrGetter interface {
	GetDelegatorWithdrawAddr(ctx context.Context, delAddr sdk.AccAddress) (sdk.AccAddress, error)
}

// DistributionKeeper is the distribution keeper surface expected by this precompile.
// It includes GetDelegatorWithdrawAddr so the precompile can call it directly.
type DistributionKeeper interface {
	cmn.DistributionKeeper
	DelegatorWithdrawAddrGetter
}
