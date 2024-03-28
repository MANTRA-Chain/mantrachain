package keeper

import (
	"context"
	"strconv"

	"github.com/MANTRA-Finance/mantrachain/x/airdrop/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateCampaign(goCtx context.Context, msg *types.MsgCreateCampaign) (*types.MsgCreateCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.guardKeeper.CheckIsAdmin(ctx, msg.Creator); err != nil {
		return nil, sdkerrors.Wrap(err, "unauthorized")
	}

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	id, _ := k.GetLastCampaignId(ctx)
	nextId := id + 1

	campaign := types.NewCampaign(nextId, creator, msg.Name, msg.Desc, msg.StartTime, msg.EndTime, msg.MtRoot, sdk.NewCoins(msg.Amount))

	if err := campaign.Validate(); err != nil {
		return nil, err
	}

	whitelisted := k.guardKeeper.WhitelistTransferAccAddresses([]string{campaign.GetReserveAddress().String()}, true)
	if err := k.bankKeeper.SendCoins(ctx, creator, campaign.GetReserveAddress(), campaign.Amounts); err != nil {
		k.guardKeeper.WhitelistTransferAccAddresses(whitelisted, false)
		return nil, err
	}
	k.guardKeeper.WhitelistTransferAccAddresses(whitelisted, false)

	k.Keeper.SetCampaign(ctx, campaign)
	k.Keeper.SetLastCampaignId(ctx, nextId)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgCreateCampaign),
			sdk.NewAttribute(types.AttributeKeyCampaignId, strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute(types.AttributeKeyCampaignAddress, campaign.CampaignAddress),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
	)

	return &types.MsgCreateCampaignResponse{
		Id:              nextId,
		CampaignAddress: campaign.GetReserveAddress().String(),
	}, nil
}

func (k msgServer) DeleteCampaign(goCtx context.Context, msg *types.MsgDeleteCampaign) (*types.MsgDeleteCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	id, _ := k.GetLastCampaignId(ctx)
	if msg.Id > id {
		return nil, sdkerrors.Wrap(types.ErrCampaignInvalidId, "invalid campaign id")
	}

	campaignIndex := types.GetCampaignIndex(strconv.FormatUint(msg.Id, 10))

	campaign, found := k.GetCampaign(ctx, campaignIndex)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrCampaignInvalidId, "campaign not found")
	}

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	campaignCreator, err := sdk.AccAddressFromBech32(campaign.GetCreator())
	if err != nil {
		return nil, err
	}

	if !campaignCreator.Equals(creator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorized")
	}

	coins := k.bankKeeper.SpendableCoins(ctx, campaign.GetReserveAddress())

	whitelisted := k.guardKeeper.WhitelistTransferAccAddresses([]string{campaign.GetReserveAddress().String()}, true)
	if err := k.bankKeeper.SendCoins(ctx, campaign.GetReserveAddress(), creator, coins); err != nil {
		k.guardKeeper.WhitelistTransferAccAddresses(whitelisted, false)
		return nil, err
	}
	k.guardKeeper.WhitelistTransferAccAddresses(whitelisted, false)

	k.Keeper.RemoveCampaign(ctx, campaignIndex)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgDeleteCampaign),
			sdk.NewAttribute(types.AttributeKeyCampaignId, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
	)

	return &types.MsgDeleteCampaignResponse{}, nil
}

func (k msgServer) PauseCampaign(goCtx context.Context, msg *types.MsgPauseCampaign) (*types.MsgPauseCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	id, _ := k.GetLastCampaignId(ctx)
	if msg.Id > id {
		return nil, sdkerrors.Wrap(types.ErrCampaignInvalidId, "invalid campaign id")
	}

	campaign, found := k.GetCampaign(ctx, types.GetCampaignIndex(strconv.FormatUint(msg.Id, 10)))
	if !found {
		return nil, sdkerrors.Wrap(types.ErrCampaignInvalidId, "campaign not found")
	}

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	campaignCreator, err := sdk.AccAddressFromBech32(campaign.GetCreator())
	if err != nil {
		return nil, err
	}

	if !campaignCreator.Equals(creator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorized")
	}

	campaign.IsPaused = true

	k.Keeper.SetCampaign(ctx, campaign)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgPauseCampaign),
			sdk.NewAttribute(types.AttributeKeyCampaignId, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
	)

	return &types.MsgPauseCampaignResponse{}, nil
}

func (k msgServer) UnpauseCampaign(goCtx context.Context, msg *types.MsgUnpauseCampaign) (*types.MsgUnpauseCampaignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	id, _ := k.GetLastCampaignId(ctx)
	if msg.Id > id {
		return nil, sdkerrors.Wrap(types.ErrCampaignInvalidId, "invalid campaign id")
	}

	campaign, found := k.GetCampaign(ctx, types.GetCampaignIndex(strconv.FormatUint(msg.Id, 10)))
	if !found {
		return nil, sdkerrors.Wrap(types.ErrCampaignInvalidId, "campaign not found")
	}

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	campaignCreator, err := sdk.AccAddressFromBech32(campaign.GetCreator())
	if err != nil {
		return nil, err
	}

	if !campaignCreator.Equals(creator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorized")
	}

	campaign.IsPaused = false

	k.Keeper.SetCampaign(ctx, campaign)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgPauseCampaign),
			sdk.NewAttribute(types.AttributeKeyCampaignId, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
	)

	return &types.MsgUnpauseCampaignResponse{}, nil
}

func (k msgServer) CampaignClaim(goCtx context.Context, msg *types.MsgCampaignClaim) (*types.MsgCampaignClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	id, _ := k.GetLastCampaignId(ctx)
	if msg.Id > id {
		return nil, sdkerrors.Wrap(types.ErrCampaignInvalidId, "invalid campaign id")
	}

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	campaign, found := k.GetCampaign(ctx, types.GetCampaignIndex(strconv.FormatUint(msg.Id, 10)))
	if !found {
		return nil, sdkerrors.Wrap(types.ErrCampaignInvalidId, "campaign not found")
	}

	if campaign.IsTerminated {
		return nil, sdkerrors.Wrap(types.ErrCampaignTerminated, "campaign has been terminated")
	}

	if campaign.GetStartTime().After(ctx.BlockTime()) {
		return nil, sdkerrors.Wrap(types.ErrCampaignNotStarted, "campaign has not started")
	}

	if campaign.GetEndTime().Before(ctx.BlockTime()) {
		return nil, sdkerrors.Wrap(types.ErrCampaignHasEnded, "campaign has ended")
	}

	if campaign.IsPaused {
		return nil, sdkerrors.Wrap(types.ErrCampaignPaused, "campaign is paused")
	}

	claimedIndex := types.GetClaimedIndex(creator, campaign.Id)

	_, found = k.GetClaimed(ctx, claimedIndex)

	if found {
		return nil, sdkerrors.Wrap(types.ErrAlreadyClaimed, "already claimed")
	}

	leafHash, err := genLeafHash(creator.String(), msg.Amount.String())
	if err != nil {
		return nil, err
	}
	merklePath, err := pathToChunks(msg.Mip)
	if err != nil {
		return nil, err
	}
	rootHash := campaign.GetMtRoot()
	leafIndex := msg.Index

	// Verify Merkle path
	valid, err := verifyMerklePath(leafHash, merklePath, rootHash, leafIndex)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, sdkerrors.Wrap(types.ErrInvalidMerklePath, "invalid merkle path")
	}

	reward := sdk.NewCoin(msg.Amount.Denom, msg.Amount.Amount)

	k.SetClaimed(ctx, types.Claimed{
		Index:       claimedIndex,
		CampaignId:  campaign.Id,
		Address:     creator.String(),
		ClaimedTime: ctx.BlockTime(),
		Amounts:     sdk.NewCoins(reward),
	})

	whitelisted := k.guardKeeper.WhitelistTransferAccAddresses([]string{campaign.GetReserveAddress().String()}, true)
	if err := k.bankKeeper.SendCoins(ctx, campaign.GetReserveAddress(), creator, sdk.NewCoins(reward)); err != nil {
		k.guardKeeper.WhitelistTransferAccAddresses(whitelisted, false)
		return nil, err
	}
	k.guardKeeper.WhitelistTransferAccAddresses(whitelisted, false)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.TypeMsgCampaignClaim),
			sdk.NewAttribute(types.AttributeKeyCampaignId, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
	)

	return &types.MsgCampaignClaimResponse{}, nil
}

func (k Keeper) TerminateEndedCampaigns(ctx sdk.Context) (err error) {
	k.IterateAllCampaigns(ctx, func(campaign types.Campaign) (stop bool) {
		if campaign.IsTerminated {
			return false
		}
		if !ctx.BlockTime().Before(campaign.EndTime) {
			if err = k.TerminateCampaign(ctx, campaign); err != nil {
				return true
			}
		}
		return false
	})
	return err
}

func (k Keeper) TerminateCampaign(ctx sdk.Context, campaign types.Campaign) error {
	if campaign.IsTerminated {
		return types.ErrCampaignAlreadyTerminated
	}
	campaignAddr := campaign.GetReserveAddress()
	balances := k.bankKeeper.SpendableCoins(ctx, campaignAddr)
	if !balances.IsZero() {
		// Guard: whitelist account address
		whitelisted := k.guardKeeper.WhitelistTransferAccAddresses([]string{campaignAddr.String()}, true)
		if err := k.bankKeeper.SendCoins(
			ctx, campaignAddr, campaign.GetCampaignCreator(), balances); err != nil {
			k.guardKeeper.WhitelistTransferAccAddresses(whitelisted, false)
			return err
		}
		k.guardKeeper.WhitelistTransferAccAddresses(whitelisted, false)
	}
	campaign.IsTerminated = true
	k.SetCampaign(ctx, campaign)
	return nil
}
