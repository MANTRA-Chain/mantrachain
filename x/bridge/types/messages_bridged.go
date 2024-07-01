package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

const (
	TypeMsgCreateMultiBridged = "create_multi_bridged"
)

var (
	_ legacytx.LegacyMsg = &MsgCreateMultiBridged{}
	_ sdk.Msg            = &MsgCreateMultiBridged{}
)

func NewMsgCreateMultiBridged(
	input Input,
	outputs []Output,
	ethTxHashes []string,
) *MsgCreateMultiBridged {
	return &MsgCreateMultiBridged{
		Input:       input,
		Outputs:     outputs,
		EthTxHashes: ethTxHashes,
	}
}

func (msg *MsgCreateMultiBridged) Route() string {
	return RouterKey
}

func (msg *MsgCreateMultiBridged) Type() string {
	return TypeMsgCreateMultiBridged
}

func (msg *MsgCreateMultiBridged) GetSigners() []sdk.AccAddress {
	inAddr, err := sdk.AccAddressFromBech32(msg.Input.Address)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{inAddr}
}

func (msg *MsgCreateMultiBridged) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateMultiBridged) ValidateBasic() error {
	if strings.TrimSpace(msg.Input.Address) == "" {
		return ErrNoInput
	}

	if len(msg.Input.Coins) != 1 {
		return ErrNoInputCoin
	}

	if msg.Input.Coins.Empty() {
		return ErrNoInputCoins
	}

	if len(msg.EthTxHashes) == 0 {
		return ErrNoEthTxHashes
	}

	if len(msg.Outputs) == 0 {
		return ErrNoOutputs
	}

	if len(msg.Outputs) != len(msg.EthTxHashes) {
		return ErrOutputEthTxHashMismatch
	}

	for _, ethTxHash := range msg.EthTxHashes {
		if strings.TrimSpace(ethTxHash) == "" {
			return ErrEmptyEthTxHash
		}
	}

	return ValidateInputOutputs(msg.Input, msg.Outputs)
}

func ValidateInputOutputs(input Input, outputs []Output) error {
	var totalOut sdk.Coins

	if err := input.ValidateBasic(); err != nil {
		return err
	}

	for _, out := range outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}

		if len(out.Coins) != 1 {
			return ErrMultipleOutputCoins
		}

		totalOut = totalOut.Add(out.Coins...)
	}

	// make sure input and outputs coins are equal
	if !input.Coins.Equal(totalOut) {
		return ErrInputOutputsCoinsMismatch
	}

	return nil
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(in.Address); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid input address: %s", err)
	}

	if !in.Coins.IsValid() {
		return errors.Wrap(errorstypes.ErrInvalidCoins, in.Coins.String())
	}

	if !in.Coins.IsAllPositive() {
		return errors.Wrap(errorstypes.ErrInvalidCoins, in.Coins.String())
	}

	return nil
}

func (out Output) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(out.Address); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid output address: %s", err)
	}

	if !out.Coins.IsValid() {
		return errors.Wrap(errorstypes.ErrInvalidCoins, out.Coins.String())
	}

	if !out.Coins.IsAllPositive() {
		return errors.Wrap(errorstypes.ErrInvalidCoins, out.Coins.String())
	}

	return nil
}
