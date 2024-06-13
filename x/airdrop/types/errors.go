package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/airdrop module sentinel errors
var (
	ErrCampaignNameTooLong           = sdkerrors.Register(ModuleName, 2, "campaign name too long")
	ErrCampaignDescTooLong           = sdkerrors.Register(ModuleName, 3, "campaign description too long")
	ErrCampaignMtRootInvalid         = sdkerrors.Register(ModuleName, 4, "merkle tree root hash must be 32 bytes")
	ErrCampaignInvalidAmount         = sdkerrors.Register(ModuleName, 5, "invalid campaign amount")
	ErrCampaignStartTimeInvalid      = sdkerrors.Register(ModuleName, 6, "invalid campaign start time")
	ErrCampaignIdInvalid             = sdkerrors.Register(ModuleName, 7, "invalid campaign id")
	ErrCampaignReserveAddressInvalid = sdkerrors.Register(ModuleName, 8, "invalid campaign reserve address")
	ErrCampaignInvalidId             = sdkerrors.Register(ModuleName, 9, "invalid campaign id")
	ErrInvalidMerklePath             = sdkerrors.Register(ModuleName, 10, "invalid merkle path")
	ErrCampaignPaused                = sdkerrors.Register(ModuleName, 11, "campaign is paused")
	ErrInvalidMerklePathIndex        = sdkerrors.Register(ModuleName, 12, "invalid merkle path index")
	ErrCampaignHasEnded              = sdkerrors.Register(ModuleName, 13, "campaign has ended")
	ErrCampaignNotStarted            = sdkerrors.Register(ModuleName, 14, "campaign has not started")
	ErrInvalidDenom                  = sdkerrors.Register(ModuleName, 15, "invalid denom")
	ErrAlreadyClaimed                = sdkerrors.Register(ModuleName, 16, "already claimed")
	ErrCampaignTerminated            = sdkerrors.Register(ModuleName, 17, "campaign has been terminated")
	ErrCampaignAlreadyTerminated     = sdkerrors.Register(ModuleName, 18, "campaign has already been terminated")
)
