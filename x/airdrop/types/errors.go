package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/airdrop module sentinel errors
var (
	ErrCampaignNameTooLong           = sdkerrors.Register(ModuleName, 1901, "campaign name too long")
	ErrCampaignDescTooLong           = sdkerrors.Register(ModuleName, 1902, "campaign description too long")
	ErrCampaignMtRootInvalid         = sdkerrors.Register(ModuleName, 1903, "merkle tree root hash must be 32 bytes")
	ErrCampaignInvalidAmount         = sdkerrors.Register(ModuleName, 1904, "invalid campaign amount")
	ErrCampaignStartTimeInvalid      = sdkerrors.Register(ModuleName, 1905, "invalid campaign start time")
	ErrCampaignIdInvalid             = sdkerrors.Register(ModuleName, 1906, "invalid campaign id")
	ErrCampaignReserveAddressInvalid = sdkerrors.Register(ModuleName, 1907, "invalid campaign reserve address")
	ErrCampaignInvalidId             = sdkerrors.Register(ModuleName, 1908, "invalid campaign id")
	ErrInvalidMerklePath             = sdkerrors.Register(ModuleName, 1909, "invalid merkle path")
	ErrCampaignPaused                = sdkerrors.Register(ModuleName, 1910, "campaign is paused")
	ErrInvalidMerklePathIndex        = sdkerrors.Register(ModuleName, 1911, "invalid merkle path index")
	ErrCampaignHasEnded              = sdkerrors.Register(ModuleName, 1912, "campaign has ended")
	ErrCampaignNotStarted            = sdkerrors.Register(ModuleName, 1913, "campaign has not started")
	ErrInvalidDenom                  = sdkerrors.Register(ModuleName, 1914, "invalid denom")
	ErrAlreadyClaimed                = sdkerrors.Register(ModuleName, 1915, "already claimed")
	ErrCampaignTerminated            = sdkerrors.Register(ModuleName, 1916, "campaign has been terminated")
	ErrCampaignAlreadyTerminated     = sdkerrors.Register(ModuleName, 1917, "campaign has already been terminated")
)
