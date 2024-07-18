package types

import (
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewCampaign(id uint64, creator sdk.AccAddress, name, desc string, startTime, endTime time.Time, mtRoot []byte, amounts sdk.Coins) Campaign {
	return Campaign{
		Index:           GetCampaignIndex(strconv.FormatUint(id, 10)),
		Id:              id,
		Name:            name,
		Description:     desc,
		Creator:         creator.String(),
		CampaignAddress: CampaignReserveAddress(id).String(),
		StartTime:       startTime,
		EndTime:         endTime,
		MtRoot:          mtRoot,
		IsPaused:        false,
		IsTerminated:    false,
		Amounts:         amounts,
	}
}

func (campaign Campaign) GetCampaignCreator() sdk.AccAddress {
	if campaign.Creator == "" {
		return nil
	}
	addr, err := sdk.AccAddressFromBech32(campaign.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (campaign Campaign) GetReserveAddress() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(campaign.CampaignAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (campaign Campaign) Validate() error {
	if campaign.Id == 0 {
		return errors.Wrap(ErrCampaignIdInvalid, "campaign id must not be 0")
	}
	if _, err := sdk.AccAddressFromBech32(campaign.CampaignAddress); err != nil {
		return errors.Wrapf(ErrCampaignReserveAddressInvalid, "invalid reserve address %s: %s", campaign.CampaignAddress, err)
	}
	if _, err := sdk.AccAddressFromBech32(campaign.Creator); err != nil {
		return errors.Wrapf(errorstypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(campaign.Name) > MaxCampaignNameLen {
		return errors.Wrapf(ErrCampaignNameTooLong, "too long campaign name, maximum %d", MaxCampaignNameLen)
	}
	if len(campaign.Description) > MaxCampaignDescriptionLen {
		return errors.Wrapf(ErrCampaignDescTooLong, "too long campaign description, maximum %d", MaxCampaignDescriptionLen)
	}
	if len(campaign.MtRoot) != 32 {
		return errors.Wrap(ErrCampaignMtRootInvalid, "merkle tree root hash must be 32 bytes")
	}
	if campaign.Amounts.IsAnyNegative() {
		return errors.Wrap(ErrCampaignInvalidAmount, "campaign amount must not be negative")
	}
	if campaign.Amounts.IsZero() {
		return errors.Wrap(ErrCampaignInvalidAmount, "campaign amount must not be zero")
	}
	if campaign.StartTime.After(campaign.EndTime) {
		return errors.Wrap(ErrCampaignStartTimeInvalid, "start time must not be after end time")
	}
	return nil
}

func CampaignReserveAddress(poolId uint64) sdk.AccAddress {
	return DeriveAddress(
		ModuleName,
		strings.Join([]string{CampaignReserveAddressPrefix, strconv.FormatUint(poolId, 10)}, ModuleAddressNameSplitter),
	)
}
