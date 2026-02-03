package distrclaim

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	vmtypes "github.com/cosmos/evm/x/vm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EVMKeeper is the minimal interface for running an EVM call via the x/vm keeper.
// This sees SDK-side EVM state updates (e.g. from ConvertCoin) that may not be visible to a
// nested call using the in-flight vm.EVM instance.
type EVMKeeper interface {
	CallEVMWithData(ctx sdk.Context, from common.Address, contract *common.Address, data []byte, commit bool, gasCap *big.Int) (*vmtypes.MsgEthereumTxResponse, error)
}
