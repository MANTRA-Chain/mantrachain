package token

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/MANTRA-Finance/mantrachain/x/token/keeper"
	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateNftCollection:
			res, err := msgServer.CreateNftCollection(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgMintNfts:
			res, err := msgServer.MintNfts(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgBurnNfts:
			res, err := msgServer.BurnNfts(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgTransferNfts:
			res, err := msgServer.TransferNfts(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgApproveNfts:
			res, err := msgServer.ApproveNfts(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgApproveAllNfts:
			res, err := msgServer.ApproveAllNfts(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgMintNft:
			res, err := msgServer.MintNft(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgBurnNft:
			res, err := msgServer.BurnNft(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgTransferNft:
			res, err := msgServer.TransferNft(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgApproveNft:
			res, err := msgServer.ApproveNft(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
			// this line is used by starport scaffolding # 1
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, errors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
