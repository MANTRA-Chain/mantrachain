package keeper

import (
	"context"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	types2 "cosmossdk.io/store/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v3/types"
	"github.com/MANTRA-Chain/mantrachain/v5/x/tokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) setBeforeSendHook(ctx sdk.Context, denom, contractAddr string) error {
	// verify that denom is an x/tokenfactory denom
	_, _, err := types.DeconstructDenom(denom)
	if err != nil {
		return err
	}

	store := k.GetDenomPrefixStore(ctx, denom)

	// delete the store for denom prefix store when cosmwasm address is nil
	if contractAddr == "" {
		store.Delete([]byte(types.BeforeSendHookAddressPrefixKey))
		return nil
	}

	_, err = sdk.AccAddressFromBech32(contractAddr)
	if err != nil {
		return err
	}

	store.Set([]byte(types.BeforeSendHookAddressPrefixKey), []byte(contractAddr))

	return nil
}

func (k Keeper) GetBeforeSendHook(ctx context.Context, denom string) string {
	store := k.GetDenomPrefixStore(ctx, denom)

	bz := store.Get([]byte(types.BeforeSendHookAddressPrefixKey))
	if bz == nil {
		return ""
	}

	return string(bz)
}

func CWCoinFromSDKCoin(in sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom:  in.GetDenom(),
		Amount: in.Amount.String(),
	}
}

// Hooks wrapper struct for bank keeper
type Hooks struct {
	k Keeper
}

var _ types.BankHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// TrackBeforeSend calls the before send listener contract suppresses any errors
func (h Hooks) TrackBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) {
	_ = h.k.callBeforeSendListener(ctx, from, to, amount, false)
}

// TrackBeforeSend calls the before send listener contract returns any errors
func (h Hooks) BlockBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	return h.k.callBeforeSendListener(ctx, from, to, amount, true)
}

// callBeforeSendListener iterates over each coin and sends corresponding sudo msg to the contract address stored in state.
// If blockBeforeSend is true, sudoMsg wraps BlockBeforeSendMsg, otherwise sudoMsg wraps TrackBeforeSendMsg.
// Note that we gas meter trackBeforeSend to prevent infinite contract calls.
// CONTRACT: this should not be called in beginBlock or endBlock since out of gas will cause this method to panic.
func (k Keeper) callBeforeSendListener(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins, blockBeforeSend bool) (err error) {
	c := sdk.UnwrapSDKContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			err = types.ErrTrackBeforeSendOutOfGas
		}
	}()

	for _, coin := range amount {
		contractAddr := k.GetBeforeSendHook(ctx, coin.Denom)
		if contractAddr != "" {
			cwAddr, err := sdk.AccAddressFromBech32(contractAddr)
			if err != nil {
				return err
			}

			var msgBz []byte

			// get msgBz, either BlockBeforeSend or TrackBeforeSend
			// Note that for trackBeforeSend, we need to gas meter computations to prevent infinite loop
			// specifically because module to module sends are not gas metered.
			// We don't need to do this for blockBeforeSend since blockBeforeSend is not called during module to module sends.
			if blockBeforeSend {
				msg := types.BlockBeforeSendSudoMsg{
					BlockBeforeSend: types.BlockBeforeSendMsg{
						From:   from.String(),
						To:     to.String(),
						Amount: CWCoinFromSDKCoin(coin),
					},
				}
				msgBz, err = json.Marshal(msg)
			} else {
				msg := types.TrackBeforeSendSudoMsg{
					TrackBeforeSend: types.TrackBeforeSendMsg{
						From:   from.String(),
						To:     to.String(),
						Amount: CWCoinFromSDKCoin(coin),
					},
				}
				msgBz, err = json.Marshal(msg)
			}
			if err != nil {
				return err
			}
			em := sdk.NewEventManager()

			childCtx := c.WithGasMeter(types2.NewGasMeter(types.BeforeSendHookGasLimit))
			_, err = k.contractKeeper.Sudo(childCtx.WithEventManager(em), cwAddr, msgBz)
			if err != nil {
				return errorsmod.Wrapf(err, "failed to call before send hook for denom %s", coin.Denom)
			}

			// consume gas used for calling contract to the parent ctx
			c.GasMeter().ConsumeGas(childCtx.GasMeter().GasConsumed(), "track before send gas")
		}
	}
	return nil
}
