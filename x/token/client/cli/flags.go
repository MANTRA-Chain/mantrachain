package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagStrict            = "strict"
	FlagCollectionCreator = "collection-creator"
	FlagCollectionId      = "collection-id"
	FlagReceiver          = "receiver"
	FlagOwner             = "owner"
	FlagDid               = "did"
)

var (
	FsMintNFT                     = flag.NewFlagSet("", flag.ContinueOnError)
	FsBurnNFT                     = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferNFT                 = flag.NewFlagSet("", flag.ContinueOnError)
	FsApproveNFT                  = flag.NewFlagSet("", flag.ContinueOnError)
	FsUpdateGuardSoulBondNFTImage = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsMintNFT.Bool(FlagStrict, false, "If true, throws an error when collection-creator/collection-id is empty or collection does not exist")
	FsMintNFT.Bool(FlagDid, false, "If true, throws an error when collection-creator/collection-id is not soul-bond collection")
	FsMintNFT.String(FlagCollectionCreator, "", "The collection creator address, if not filled, the default is the sender of the transaction")
	FsMintNFT.String(FlagCollectionId, "", "The collection id, if not filled, the default is the creator's 'default' collection")
	FsMintNFT.String(FlagReceiver, "", "NFT receiver's address on mint")

	FsBurnNFT.Bool(FlagStrict, false, "If true, throws an error when collection-creator/collection-id is empty or collection does not exist")
	FsBurnNFT.String(FlagCollectionCreator, "", "The collection creator address, if not filled, the default is the sender of the transaction")
	FsBurnNFT.String(FlagCollectionId, "", "The collection id, if not filled, the default is the creator's 'default' collection")

	FsTransferNFT.Bool(FlagStrict, false, "If true, throws an error when collection-creator/collection-id is empty or collection does not exist")
	FsTransferNFT.String(FlagCollectionCreator, "", "The collection creator address, if not filled, the default is the sender of the transaction")
	FsTransferNFT.String(FlagCollectionId, "", "The collection id, if not filled, the default is the creator's 'default' collection")

	FsApproveNFT.Bool(FlagStrict, false, "If true, throws an error when collection-creator/collection-id is empty or collection does not exist")
	FsApproveNFT.String(FlagCollectionCreator, "", "The collection creator address, if not filled, the default is the sender of the transaction")
	FsApproveNFT.String(FlagCollectionId, "", "The collection id, if not filled, the default is the creator's 'default' collection")

	FsUpdateGuardSoulBondNFTImage.String(FlagOwner, "", "The owner address of the NFT")
}
