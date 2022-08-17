package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagMarketplaceCreator = "marketplace-creator"
	FlagMarketplaceId      = "marketplace-id"
	FlagCollectionCreator  = "collection-creator"
	FlagCollectionId       = "collection-id"
	FlagReceiver           = "receiver"
)

var (
	FsWithdrawNftReward = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsWithdrawNftReward.String(FlagMarketplaceCreator, "", "The marketplace creator address")
	FsWithdrawNftReward.String(FlagMarketplaceId, "", "The marketplace id")
	FsWithdrawNftReward.String(FlagCollectionCreator, "", "The collection creator address")
	FsWithdrawNftReward.String(FlagCollectionId, "", "The collection id")
	FsWithdrawNftReward.String(FlagReceiver, "", "The withdraw reward receiver")
}
