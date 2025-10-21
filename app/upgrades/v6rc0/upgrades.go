package v6rc0

import (
	"context"
	"strings"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/MANTRA-Chain/mantrachain/v7/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// ChainCoinInfo is a map of the chain id and its corresponding EvmCoinInfo
// that allows initializing the app with different coin info based on the
// chain id
var ChainCoinInfo = evmtypes.EvmCoinInfo{
	Denom:         "uom",
	ExtendedDenom: "aom",
	DisplayDenom:  "om",
	Decimals:      evmtypes.SixDecimals.Uint32(),
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
	storekeys map[string]*storetypes.KVStoreKey,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting module migrations...")

		keepers.BankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    ChainCoinInfo.Denom,
					Exponent: 0,
				},
				{
					Denom:    ChainCoinInfo.DisplayDenom,
					Exponent: ChainCoinInfo.Decimals,
				},
			},
			Base:    ChainCoinInfo.Denom,
			Display: ChainCoinInfo.DisplayDenom,
			Name:    ChainCoinInfo.DisplayDenom,
			Symbol:  "OM",
		})

		evmParams := keepers.EVMKeeper.GetParams(ctx)
		evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: ChainCoinInfo.ExtendedDenom}
		if err := keepers.EVMKeeper.SetParams(ctx, evmParams); err != nil {
			return vm, err
		}

		// update contract owner for all existing tokenfactory token_pairs
		pairs := keepers.Erc20Keeper.GetTokenPairs(ctx)
		for _, pair := range pairs {
			if strings.HasPrefix(pair.Denom, "factory/") {
				pair.ContractOwner = erc20types.OWNER_MODULE
				keepers.Erc20Keeper.SetTokenPair(ctx, pair)
			}
		}

		disableList := []string{
			"wasm/cosmos.evm.erc20.v1.MsgRegisterERC20",
			"wasm/cosmos.authz.v1beta1.MsgExec",
		}
		for _, msg := range disableList {
			if err := keepers.CircuitKeeper.DisableList.Set(ctx, msg); err != nil {
				return vm, err
			}
		}

		if err := keepers.EVMKeeper.InitEvmCoinInfo(ctx); err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v6.0.0-rc0 complete")
		return vm, nil
	}
}
