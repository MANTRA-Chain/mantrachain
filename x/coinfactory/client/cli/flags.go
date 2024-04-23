package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagMintTo   = "mint-to"
	FlagBurnFrom = "burn-from"
)

var (
	FsMintTo   = flag.NewFlagSet("", flag.ContinueOnError)
	FsBurnFrom = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsMintTo.String(FlagMintTo, "", "If set, it will mint the token to the specified address")
	FsBurnFrom.String(FlagBurnFrom, "", "If set, it will burn the token from the specified address")
}
