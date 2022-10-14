package keeper

import (
	"strconv"

	"github.com/LimeChain/mantrachain/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type WasmExecutor struct {
	ctx                sdk.Context
	wasmViewKeeper     types.WasmViewKeeper
	wasmContractKeeper types.WasmContractOpsKeeper
}

func NewWasmExecutor(ctx sdk.Context, wasmViewKeeper types.WasmViewKeeper, wasmContractKeeper types.WasmContractOpsKeeper) *WasmExecutor {
	return &WasmExecutor{
		ctx:                ctx,
		wasmViewKeeper:     wasmViewKeeper,
		wasmContractKeeper: wasmContractKeeper,
	}
}

func (c *WasmExecutor) Transfer(contractAddress sdk.AccAddress, creator sdk.AccAddress, receiver sdk.AccAddress, amount uint64) (err error) {
	transferData := []byte("{\"transfer\": {\"recipient\": \"" + receiver.String() + "\", \"amount\": \"" + strconv.FormatUint(amount, 10) + "\"}}")
	_, err = c.wasmContractKeeper.Execute(c.ctx, contractAddress, creator, transferData, nil)
	return
}
