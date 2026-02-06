package distrclaim

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	cmn "github.com/cosmos/evm/precompiles/common"
	vmtypes "github.com/cosmos/evm/x/vm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EVMKeeper is the minimal interface for running an EVM call via the x/vm keeper.
// This sees SDK-side EVM state updates (e.g. from ConvertCoin) that may not be visible to a
// nested call using the in-flight vm.EVM instance.
type EVMKeeper interface {
	CallEVMWithData(ctx sdk.Context, from common.Address, contract *common.Address, data []byte, commit bool, gasCap *big.Int) (*vmtypes.MsgEthereumTxResponse, error)
}

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
