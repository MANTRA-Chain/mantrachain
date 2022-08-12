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
	FsImportCollection = flag.NewFlagSet("", flag.ContinueOnError)
	FsBuyNft           = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsImportCollection.String(FlagMarketplaceCreator, "", "The marketplace creator address")
	FsImportCollection.String(FlagMarketplaceId, "", "The marketplace id")
	FsImportCollection.String(FlagCollectionCreator, "", "The collection creator address")
	FsImportCollection.String(FlagCollectionId, "", "The collection id")

	FsBuyNft.String(FlagMarketplaceCreator, "", "The marketplace creator address")
	FsBuyNft.String(FlagMarketplaceId, "", "The marketplace id")
	FsBuyNft.String(FlagCollectionCreator, "", "The collection creator address")
	FsBuyNft.String(FlagCollectionId, "", "The collection id")
}
