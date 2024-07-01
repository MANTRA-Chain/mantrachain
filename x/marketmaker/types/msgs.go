package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ sdk.Msg            = &MsgApplyMarketMaker{}
	_ sdk.Msg            = &MsgClaimIncentives{}
	_ legacytx.LegacyMsg = &MsgApplyMarketMaker{}
	_ legacytx.LegacyMsg = &MsgClaimIncentives{}
)

// Message types for the marketmaker module
const (
	TypeMsgApplyMarketMaker = "apply_market_maker"
	TypeMsgClaimIncentives  = "claim_incentives"
)

// NewMsgApplyMarketMaker creates a new MsgApplyMarketMaker.
func NewMsgApplyMarketMaker(
	marketMaker sdk.AccAddress,
	pairIds []uint64,
) *MsgApplyMarketMaker {
	return &MsgApplyMarketMaker{
		Address: marketMaker.String(),
		PairIds: pairIds,
	}
}

func (msg MsgApplyMarketMaker) Route() string { return RouterKey }

func (msg MsgApplyMarketMaker) Type() string { return TypeMsgApplyMarketMaker }

func (msg MsgApplyMarketMaker) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid address %q: %v", msg.Address, err)
	}
	if len(msg.PairIds) == 0 {
		return errors.Wrap(errorstypes.ErrInvalidRequest, "pair ids must not be empty")
	}
	pairMap := make(map[uint64]struct{})
	for _, pair := range msg.PairIds {
		if _, ok := pairMap[pair]; ok {
			return errors.Wrapf(ErrInvalidPairId, "duplicated pair id %d", pair)
		}
		pairMap[pair] = struct{}{}
	}
	return nil
}

func (msg *MsgApplyMarketMaker) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgApplyMarketMaker) GetAccAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgClaimIncentives creates a new MsgClaimIncentives.
func NewMsgClaimIncentives(
	marketMaker sdk.AccAddress,
) *MsgClaimIncentives {
	return &MsgClaimIncentives{
		Address: marketMaker.String(),
	}
}

func (msg MsgClaimIncentives) Route() string { return RouterKey }

func (msg MsgClaimIncentives) Type() string { return TypeMsgClaimIncentives }

func (msg MsgClaimIncentives) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid address %q: %v", msg.Address, err)
	}
	return nil
}

func (msg *MsgClaimIncentives) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgClaimIncentives) GetAccAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return addr
}
