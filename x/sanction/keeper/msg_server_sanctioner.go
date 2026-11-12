package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	if err := k.checkAuthority(msg.Authority); err != nil {
		return nil, err
	}
	if err := k.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) UpdateSanctioners(ctx context.Context, msg *types.MsgUpdateSanctioners) (*types.MsgUpdateSanctionersResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	if err := k.checkAuthority(msg.Authority); err != nil {
		return nil, err
	}
	if err := k.SetSanctioners(ctx, msg.Sanctioners); err != nil {
		return nil, err
	}
	sanctioners, err := k.GetSanctioners(ctx)
	if err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventSanctionersUpdated{
		Sanctioners: sanctioners,
		Height:      sdkCtx.BlockHeight(),
	}); err != nil {
		return nil, err
	}
	return &types.MsgUpdateSanctionersResponse{}, nil
}

// The handlers below check only that the signer may send the message. The
// keeper applies the limits and the rules on target addresses.

// SanctionAddresses accepts the gov authority or a sanctioner with any role.
func (k msgServer) SanctionAddresses(ctx context.Context, msg *types.MsgSanctionAddresses) (*types.MsgSanctionAddressesResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	if msg.Signer != k.GetAuthority() {
		if _, err := k.getSanctioner(ctx, msg.Signer); err != nil {
			return nil, err
		}
	}
	if err := k.Keeper.SanctionAddresses(ctx, msg.Addresses, msg.Signer, msg.Reason); err != nil {
		return nil, err
	}
	return &types.MsgSanctionAddressesResponse{}, nil
}

// UnsanctionAddresses accepts the gov authority or a full sanctioner.
func (k msgServer) UnsanctionAddresses(ctx context.Context, msg *types.MsgUnsanctionAddresses) (*types.MsgUnsanctionAddressesResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	if msg.Signer != k.GetAuthority() {
		if err := k.checkFullSanctioner(ctx, msg.Signer); err != nil {
			return nil, err
		}
	}
	if err := k.Keeper.UnsanctionAddresses(ctx, msg.Addresses, msg.Signer, msg.Reason); err != nil {
		return nil, err
	}
	return &types.MsgUnsanctionAddressesResponse{}, nil
}

// RevokeSanctioner accepts a full sanctioner.
func (k msgServer) RevokeSanctioner(ctx context.Context, msg *types.MsgRevokeSanctioner) (*types.MsgRevokeSanctionerResponse, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	if err := k.checkFullSanctioner(ctx, msg.Signer); err != nil {
		return nil, err
	}
	if err := k.Keeper.RevokeSanctioner(ctx, msg.Address, msg.Signer, msg.Reason); err != nil {
		return nil, err
	}
	return &types.MsgRevokeSanctionerResponse{}, nil
}

func (k msgServer) checkAuthority(authority string) error {
	if k.GetAuthority() != authority {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authority; expected %s, got %s", k.GetAuthority(), authority)
	}
	return nil
}

// checkFullSanctioner fails unless addr is a sanctioner with the full role.
func (k msgServer) checkFullSanctioner(ctx context.Context, addr string) error {
	s, err := k.getSanctioner(ctx, addr)
	if err != nil {
		return err
	}
	if s.Role != types.SanctionerRole_SANCTIONER_ROLE_FULL {
		return errorsmod.Wrapf(types.ErrRoleNotPermitted, "%s has role %s; this needs a full sanctioner", s.Address, s.Role)
	}
	return nil
}
