package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

const (
	TypeMsgCreateNftCollection = "create_nft_collection"
)

var (
	_ sdk.Msg            = &MsgCreateNftCollection{}
	_ legacytx.LegacyMsg = &MsgCreateNftCollection{}
)

func NewMsgCreateNftCollection(creator string, collection *MsgCreateNftCollectionMetadata) *MsgCreateNftCollection {
	return &MsgCreateNftCollection{
		Creator:    creator,
		Collection: collection,
	}
}

func (msg *MsgCreateNftCollection) Route() string {
	return RouterKey
}

func (msg *MsgCreateNftCollection) Type() string {
	return TypeMsgCreateNftCollection
}

func (msg *MsgCreateNftCollection) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateNftCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Collection == nil {
		return errors.Wrapf(errorstypes.ErrKeyNotFound, "nft collection is empty")
	}
	if strings.TrimSpace(msg.Collection.Id) == "" {
		return errors.Wrap(errorstypes.ErrKeyNotFound, "collection id should not be empty")
	}
	return nil
}
