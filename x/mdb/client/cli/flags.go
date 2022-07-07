package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagStrict   = "strict"
	FlagReceiver = "receiver"
)

var (
	FsNFT = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsNFT.String(FlagStrict, "", "Should return error on empty or missing nft/nfts/collection")
	FsNFT.String(FlagReceiver, "", "NFT/NFTs Mint/Approve/Transfer Receiver")
}
