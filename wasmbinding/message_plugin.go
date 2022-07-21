package wasmbinding

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/LimeChain/mantrachain/wasmbinding/bindings"
	mdbkeeper "github.com/LimeChain/mantrachain/x/mdb/keeper"
	mdbtypes "github.com/LimeChain/mantrachain/x/mdb/types"
)

func CustomMessageDecorator(mdbkeeper *mdbkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:   old,
			mdbkeeper: mdbkeeper,
		}
	}
}

type CustomMessenger struct {
	wrapped   wasmkeeper.Messenger
	mdbkeeper *mdbkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		var contractMsg bindings.MantraMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, sdkerrors.Wrap(err, "mantra msg")
		}
		if contractMsg.Mdb != nil {
			if contractMsg.Mdb.CreateNftCollection != nil {
				return m.createNftCollection(ctx, contractAddr, contractMsg.Mdb.CreateNftCollection)
			}
		}
	}

	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) createNftCollection(ctx sdk.Context, contractAddr sdk.AccAddress, createNftCollection *bindings.CreateNftCollection) ([]sdk.Event, [][]byte, error) {
	err := PerformCreateNftCollection(m.mdbkeeper, ctx, contractAddr, createNftCollection)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform create nft collection")
	}
	return nil, nil, nil
}

func PerformCreateNftCollection(f *mdbkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, createNftCollection *bindings.CreateNftCollection) error {
	if createNftCollection.Collection == nil {
		return wasmvmtypes.InvalidRequest{Err: "create nft collection null collection"}
	}

	msgServer := mdbkeeper.NewMsgServerImpl(*f)

	msgCreateNftCollection := mdbtypes.NewMsgCreateNftCollection(contractAddr.String(), &mdbtypes.MsgCreateNftCollectionMetadata{
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
