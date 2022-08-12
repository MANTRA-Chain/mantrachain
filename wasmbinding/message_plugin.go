package wasmbinding

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/LimeChain/mantrachain/wasmbinding/bindings"
	tokenkeeper "github.com/LimeChain/mantrachain/x/token/keeper"
	tokentypes "github.com/LimeChain/mantrachain/x/token/types"
)

func CustomMessageDecorator(tokenkeeper *tokenkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:     old,
			tokenkeeper: tokenkeeper,
		}
	}
}

type CustomMessenger struct {
	wrapped     wasmkeeper.Messenger
	tokenkeeper *tokenkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		var contractMsg bindings.MantraMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, sdkerrors.Wrap(err, "mantra msg")
		}
		if contractMsg.Token != nil {
			if contractMsg.Token.CreateNftCollection != nil {
				return m.createNftCollection(ctx, contractAddr, contractMsg.Token.CreateNftCollection)
			}
		}
	}

	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) createNftCollection(ctx sdk.Context, contractAddr sdk.AccAddress, createNftCollection *bindings.CreateNftCollection) ([]sdk.Event, [][]byte, error) {
	err := PerformCreateNftCollection(m.tokenkeeper, ctx, contractAddr, createNftCollection)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform create nft collection")
	}
	return nil, nil, nil
}

func PerformCreateNftCollection(f *tokenkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, createNftCollection *bindings.CreateNftCollection) error {
	if createNftCollection.Collection == nil {
		return wasmvmtypes.InvalidRequest{Err: "create nft collection null collection"}
	}

	msgServer := tokenkeeper.NewMsgServerImpl(*f)

	msgCreateNftCollection := tokentypes.NewMsgCreateNftCollection(contractAddr.String(), &tokentypes.MsgCreateNftCollectionMetadata{
		Id: createNftCollection.Collection.Id,
	})

	if err := msgCreateNftCollection.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "failed validating MsgCreateDenom")
	}

	_, err := msgServer.CreateNftCollection(
		sdk.WrapSDKContext(ctx),
		msgCreateNftCollection,
	)
	if err != nil {
		return sdkerrors.Wrap(err, "creating nft collection")
	}
	return nil
}
