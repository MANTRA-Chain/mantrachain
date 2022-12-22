package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMint = "mint"
)

var (
	_ sdk.Msg = &MsgMint{}
)

func NewMsgMint(creator string, bridgeCreator string, bridgeId string,
	mint *MsgMintListMetadata,
) *MsgMint {
	return &MsgMint{
		Creator:       creator,
		BridgeCreator: bridgeCreator,
		BridgeId:      bridgeId,
		Mint:          mint,
	}
}

func (msg *MsgMint) Route() string {
	return RouterKey
}

func (msg *MsgMint) Type() string {
	return TypeMsgMint
}

func (msg *MsgMint) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.BridgeCreator) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bridge creator should not be empty")
	}
	if strings.TrimSpace(msg.BridgeId) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bridge id should not be empty")
	}
	if msg.Mint == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "mint should not be empty")
	}
	if msg.Mint.MintList == nil || len(msg.Mint.MintList) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "mint list should not be empty")
	}
	for i, mint := range msg.Mint.MintList {
		_, err = sdk.AccAddressFromBech32(mint.Receiver)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid mint mint list receiver address (%s), index %d", err, i)
		}
		if strings.TrimSpace(mint.TxHash) == "" {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "mint mint list tx hash should not be empty, index %d", i)
		}
	}
	return nil
}
