package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagMarketplaceCreator = "marketplace-creator"
	FlagMarketplaceId      = "marketplace-id"
	FlagCollectionCreator  = "collection-creator"
	FlagCollectionId       = "collection-id"
)

var (
	FsImportNftCollection = flag.NewFlagSet("", flag.ContinueOnError)
	FsBuyNft              = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsImportNftCollection.String(FlagMarketplaceCreator, "", "The marketplace creator address")
	FsImportNftCollection.String(FlagMarketplaceId, "", "The marketplace id")
	FsImportNftCollection.String(FlagCollectionCreator, "", "The collection creator address")
	FsImportNftCollection.String(FlagCollectionId, "", "The collection id")

	FsBuyNft.String(FlagMarketplaceCreator, "", "The marketplace creator address")
	FsBuyNft.String(FlagMarketplaceId, "", "The marketplace id")
	FsBuyNft.String(FlagCollectionCreator, "", "The collection creator address")
	FsBuyNft.String(FlagCollectionId, "", "The collection id")
}
