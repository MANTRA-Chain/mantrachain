package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagStrictCollection  = "strict-collection"
	FlagCollectionCreator = "collection-creator"
	FlagCollectionId      = "collection-id"
	FlagReceiver          = "receiver"
)

var (
	FsMintNFT     = flag.NewFlagSet("", flag.ContinueOnError)
	FsBurnNFT     = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferNFT = flag.NewFlagSet("", flag.ContinueOnError)
	FsApproveNFT  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsMintNFT.String(FlagStrictCollection, "", "If true will return error when collection-creator/collection-id flag is empty or collection does not exist")
	FsMintNFT.String(FlagCollectionCreator, "", "The collection creator address, if not filled, the default is the sender of the transaction")
	FsMintNFT.String(FlagCollectionId, "", "The collection id, if not filled, the default is the creator's 'default' collection")
	FsMintNFT.String(FlagReceiver, "", "NFT receiver address on mint")

	FsBurnNFT.String(FlagStrictCollection, "", "If true will return error when collection-creator/collection-id flag is empty or collection does not exist")
	FsBurnNFT.String(FlagCollectionCreator, "", "The collection creator address, if not filled, the default is the sender of the transaction")
	FsBurnNFT.String(FlagCollectionId, "", "The collection id, if not filled, the default is the creator's 'default' collection")

	FsTransferNFT.String(FlagStrictCollection, "", "If true will return error when collection-creator/collection-id flag is empty or collection does not exist")
	FsTransferNFT.String(FlagCollectionCreator, "", "The collection creator address, if not filled, the default is the sender of the transaction")
	FsTransferNFT.String(FlagCollectionId, "", "The collection id, if not filled, the default is the creator's 'default' collection")

	FsApproveNFT.String(FlagStrictCollection, "", "If true will return error when collection-creator/collection-id flag is empty or collection does not exist")
	FsApproveNFT.String(FlagCollectionCreator, "", "The collection creator address, if not filled, the default is the sender of the transaction")
	FsApproveNFT.String(FlagCollectionId, "", "The collection id, if not filled, the default is the creator's 'default' collection")

}
