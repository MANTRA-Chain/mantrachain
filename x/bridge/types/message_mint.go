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
	mint *MsgMintMetadata,
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
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bridge creator should not be empty with strict-bridge flag")
	}
	if strings.TrimSpace(msg.BridgeId) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "bridge id should not be empty with strict-bridge flag")
	}
	return nil
}
