package keeper

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/LimeChain/mantrachain/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateCw20Contract(goCtx context.Context, msg *types.MsgCreateCw20Contract) (*types.MsgCreateCw20ContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	adminAccount, err := sdk.AccAddressFromBech32(conf.AdminAccount)

	if err != nil {
		return nil, err
	}

	if !creator.Equals(adminAccount) {
		return nil, sdkerrors.Wrapf(types.ErrAdminAccountParamMismatch, "admin account param %s does not match the creator %s", adminAccount.String(), creator.String())
	}

	// Check if the value already exists
	_, isFound := k.GetCw20Contract(ctx)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "already set")
	}

	var cw20Contract = types.Cw20Contract{
		Creator: msg.Creator,
		CodeId:  msg.CodeId,
		Ver:     msg.Ver,
		Path:    msg.Path,
	}

	if msg.CodeId == 0 {
		if strings.TrimSpace(msg.Path) == "" {
			return nil, sdkerrors.Wrap(types.ErrInvalidPath, "path should not be empty")
		}

		cw20Code, err := ioutil.ReadFile(msg.Path)

		if err != nil {
			return nil, err
		}

		wasmExecutor := NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)
		codeId, err := wasmExecutor.Create(creator, cw20Code)

		if err != nil {
			return nil, err
		}

		cw20Contract.CodeId = codeId
	}

	k.SetCw20Contract(
		ctx,
		cw20Contract,
	)

	// TODO: Trigger event

	return &types.MsgCreateCw20ContractResponse{
		Creator: msg.Creator,
		CodeId:  cw20Contract.CodeId,
		Ver:     cw20Contract.Ver,
		Path:    cw20Contract.Path,
	}, nil
}

func (k msgServer) UpdateCw20Contract(goCtx context.Context, msg *types.MsgUpdateCw20Contract) (*types.MsgUpdateCw20ContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	adminAccount, err := sdk.AccAddressFromBech32(conf.AdminAccount)

	if err != nil {
		return nil, err
	}

	if !creator.Equals(adminAccount) {
		return nil, sdkerrors.Wrapf(types.ErrAdminAccountParamMismatch, "admin account param %s does not match the creator %s", adminAccount.String(), creator.String())
	}

	// Check if the value exists
	valFound, isFound := k.GetCw20Contract(ctx)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var cw20Contract = types.Cw20Contract{
		Creator: valFound.Creator,
		CodeId:  msg.CodeId,
		Ver:     msg.Ver,
		Path:    msg.Path,
	}

	if msg.CodeId == 0 {
		if strings.TrimSpace(msg.Path) == "" {
			return nil, sdkerrors.Wrap(types.ErrInvalidPath, "path should not be empty")
		}

		cw20Code, err := ioutil.ReadFile(msg.Path)

		if err != nil {
			return nil, err
		}

		wasmExecutor := NewWasmExecutor(ctx, k.wasmViewKeeper, k.wasmContractKeeper)
		codeId, err := wasmExecutor.Create(creator, cw20Code)

		if err != nil {
			return nil, err
		}

		cw20Contract.CodeId = codeId
	}

	k.SetCw20Contract(ctx, cw20Contract)

	// TODO: Trigger event

	return &types.MsgUpdateCw20ContractResponse{
		Creator: msg.Creator,
		CodeId:  cw20Contract.CodeId,
		Ver:     cw20Contract.Ver,
		Path:    cw20Contract.Path,
	}, nil
}

func (k msgServer) DeleteCw20Contract(goCtx context.Context, msg *types.MsgDeleteCw20Contract) (*types.MsgDeleteCw20ContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	conf := k.GetParams(ctx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)

	if err != nil {
		return nil, err
	}

	adminAccount, err := sdk.AccAddressFromBech32(conf.AdminAccount)

	if err != nil {
		return nil, err
	}

	// TODO: Check if can use roles instead of genesis state admin account
	if !creator.Equals(adminAccount) {
		return nil, sdkerrors.Wrapf(types.ErrAdminAccountParamMismatch, "admin account param %s does not match the creator %s", adminAccount.String(), creator.String())
	}

	// Check if the value exists
	valFound, isFound := k.GetCw20Contract(ctx)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveCw20Contract(ctx)

	// TODO: Trigger event

	return &types.MsgDeleteCw20ContractResponse{
		Creator: msg.Creator,
		CodeId:  valFound.CodeId,
		Ver:     valFound.Ver,
		Path:    valFound.Path,
	}, nil
}
