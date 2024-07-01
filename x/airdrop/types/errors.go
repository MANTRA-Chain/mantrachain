package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/airdrop module sentinel errors
var (
	ErrInvalidSigner                 = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrCampaignNameTooLong           = sdkerrors.Register(ModuleName, 1101, "campaign name too long")
	ErrCampaignDescTooLong           = sdkerrors.Register(ModuleName, 1102, "campaign description too long")
	ErrCampaignMtRootInvalid         = sdkerrors.Register(ModuleName, 1103, "merkle tree root hash must be 32 bytes")
	ErrCampaignInvalidAmount         = sdkerrors.Register(ModuleName, 1104, "invalid campaign amount")
	ErrCampaignStartTimeInvalid      = sdkerrors.Register(ModuleName, 1105, "invalid campaign start time")
	ErrCampaignIdInvalid             = sdkerrors.Register(ModuleName, 1106, "invalid campaign id")
	ErrCampaignReserveAddressInvalid = sdkerrors.Register(ModuleName, 1107, "invalid campaign reserve address")
	ErrCampaignInvalidId             = sdkerrors.Register(ModuleName, 1108, "invalid campaign id")
	ErrInvalidMerklePath             = sdkerrors.Register(ModuleName, 1109, "invalid merkle path")
	ErrCampaignPaused                = sdkerrors.Register(ModuleName, 1110, "campaign is paused")
	ErrInvalidMerklePathIndex        = sdkerrors.Register(ModuleName, 1111, "invalid merkle path index")
	ErrCampaignHasEnded              = sdkerrors.Register(ModuleName, 1112, "campaign has ended")
	ErrCampaignNotStarted            = sdkerrors.Register(ModuleName, 1113, "campaign has not started")
	ErrInvalidDenom                  = sdkerrors.Register(ModuleName, 1114, "invalid denom")
	ErrAlreadyClaimed                = sdkerrors.Register(ModuleName, 1115, "already claimed")
	ErrCampaignTerminated            = sdkerrors.Register(ModuleName, 1116, "campaign has been terminated")
	ErrCampaignAlreadyTerminated     = sdkerrors.Register(ModuleName, 1117, "campaign has already been terminated")
)
