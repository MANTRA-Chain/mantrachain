package keeper

import (
	"github.com/LimeChain/mantrachain/x/vault/types"
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

func (c *WasmExecutor) Mint(contractAddress sdk.AccAddress, creator sdk.AccAddress, receiver sdk.AccAddress, amount string) (err error) {
	mintData := []byte("{\"mint\": {\"recipient\": \"" + receiver.String() + "\", \"amount\": \"" + amount + "\"}}")
	_, err = c.wasmContractKeeper.Execute(c.ctx, contractAddress, creator, mintData, nil)
	return
}

func (c *WasmExecutor) TransferFrom(contractAddress sdk.AccAddress, creator sdk.AccAddress, owner sdk.AccAddress, receiver sdk.AccAddress, amount string) (err error) {
	transferData := []byte("{\"transfer_from\": {\"owner\": \"" + owner.String() + "\", \"recipient\": \"" + receiver.String() + "\", \"amount\": \"" + amount + "\"}}")
	_, err = c.wasmContractKeeper.Execute(c.ctx, contractAddress, creator, transferData, nil)
	return
}

func (c *WasmExecutor) IncreaseAllowance(contractAddress sdk.AccAddress, creator sdk.AccAddress, spender sdk.AccAddress, amount string) (err error) {
	approveData := []byte("{\"increase_allowance\": {\"spender\": \"" + spender.String() + "\", \"amount\": \"" + amount + "\"}}")
	_, err = c.wasmContractKeeper.Execute(c.ctx, contractAddress, creator, approveData, nil)
	return
}