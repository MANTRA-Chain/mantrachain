package keeper

import (
	"bytes"
	"context"

	"github.com/MANTRA-Finance/mantrachain/x/nft/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = Keeper{}

// Send implements Send method of the types.MsgServer.
func (k Keeper) Send(goCtx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	if len(msg.ClassId) == 0 {
		return nil, types.ErrEmptyClassID
	}

	if len(msg.Id) == 0 {
		return nil, types.ErrEmptyNFTID
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", msg.Sender)
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid receiver address (%s)", msg.Receiver)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := k.GetOwner(ctx, msg.ClassId, msg.Id)
	if !bytes.Equal(owner, sender) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not the owner of nft %s", msg.Sender, msg.Id)
	}

	if err := k.Transfer(ctx, msg.ClassId, msg.Id, receiver); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvent(&types.EventSend{
		ClassId:  msg.ClassId,
		Id:       msg.Id,
		Sender:   msg.Sender,
		Receiver: msg.Receiver,
	})
	return &types.MsgSendResponse{}, nil
}
