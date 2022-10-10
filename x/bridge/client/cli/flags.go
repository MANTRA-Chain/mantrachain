package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagBridgeCreator = "bridge-creator"
	FlagBridgeId      = "bridge-id"
)

var (
	FsMint = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsMint.String(FlagBridgeCreator, "", "The bridge creator address")
	FsMint.String(FlagBridgeId, "", "The bridge id")
}
