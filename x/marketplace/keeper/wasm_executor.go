package keeper

import (
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

func (c *WasmExecutor) Transfer(contractAddress sdk.AccAddress, creator sdk.AccAddress, receiver sdk.AccAddress, amount string) (err error) {
	transferData := []byte("{\"transfer\": {\"recipient\": \"" + receiver.String() + "\", \"amount\": \"" + amount + "\"}}")
	_, err = c.wasmContractKeeper.Execute(c.ctx, contractAddress, creator, transferData, nil)
	return
}

func (c *WasmExecutor) Burn(contractAddress sdk.AccAddress, creator sdk.AccAddress, amount string) (err error) {
	burnData := []byte("{\"burn\": {\"amount\": \"" + amount + "\"}}")
	_, err = c.wasmContractKeeper.Execute(c.ctx, contractAddress, creator, burnData, nil)
	return
}
