package keeper

import (
	"encoding/json"
	"strconv"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/LimeChain/mantrachain/x/bridge/types"
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

func (c *WasmExecutor) GetMinter(
	contractAddress sdk.AccAddress,
) (sdk.AccAddress, error) {
	queryData := []byte("{\"minter\": {}}")

	bz, err := c.wasmViewKeeper.QuerySmart(c.ctx, contractAddress, queryData)

	if err != nil {
		return nil, err
	}

	var result map[string]string

	json.Unmarshal(bz, &result)

	minter, err := sdk.AccAddressFromBech32(result["minter"])

	if err != nil {
		return nil, err
	}

	return minter, nil
}

func (c *WasmExecutor) Create(creator sdk.AccAddress, wasmCode []byte) (codeId uint64, err error) {
	codeId, err = c.wasmContractKeeper.Create(c.ctx, creator, wasmCode, &wasmtypes.AccessConfig{
		Permission: wasmtypes.AccessTypeEverybody,
		Address:    creator.String(),
	})
	return
}

func (c *WasmExecutor) Instantiate(
	codeId uint64,
	creator,
	admin sdk.AccAddress,
	name string,
	symbol string,
	decimals uint32,
	initialBalances []*types.MsgCw20InitialBalances,
	mint *types.MsgCw20Mint,
) (contractAddress sdk.AccAddress, err error) {
	strInitialBalances, err := json.Marshal(initialBalances)

	if err != nil {
		return nil, err
	}

	strMint, err := json.Marshal(mint)

	if err != nil {
		return nil, err
	}

	instantiateData := []byte("{\"initial_balances\": " + string(strInitialBalances) + ", \"decimals\": " + strconv.FormatUint(uint64(decimals), 10) + ", \"name\": \"" + name + "\"" + ", \"symbol\": \"" + symbol + "\"" + ", \"mint\": " + string(strMint) + "}")

	contractAddress, _, err = c.wasmContractKeeper.Instantiate(c.ctx, codeId, creator, admin, instantiateData, name+" CW20", nil)
	return
}

func (c *WasmExecutor) Mint(contractAddress sdk.AccAddress, creator sdk.AccAddress, receiver sdk.AccAddress, amount uint64) (err error) {
	mintData := []byte("{\"mint\": {\"recipient\": \"" + receiver.String() + "\", \"amount\": \"" + strconv.FormatUint(amount, 10) + "\"}}")
	_, err = c.wasmContractKeeper.Execute(c.ctx, contractAddress, creator, mintData, nil)
	return
}
