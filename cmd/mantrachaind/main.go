package main

import (
	"fmt"
	"os"

	clienthelpers "cosmossdk.io/client/v2/helpers"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/MANTRA-Chain/mantrachain/v5/app"
	"github.com/MANTRA-Chain/mantrachain/v5/cmd/mantrachaind/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmdcfg "github.com/cosmos/evm/cmd/evmd/config"
)

func main() {
	sdk.SetCoinDenomRegex(MantraCoinDenomRegex)
	setupConfig()
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, clienthelpers.EnvPrefix, app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}

const (
	HumanCoinUnit = "om"
	BaseCoinUnit  = "uom"
	OmExponent    = 6

	DefaultBondDenom = BaseCoinUnit
)

var (
	Bech32Prefix = "mantra"
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key.
	Bech32PrefixAccPub = Bech32Prefix + "pub"
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address.
	Bech32PrefixValAddr = Bech32Prefix + "valoper"
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key.
	Bech32PrefixValPub = Bech32Prefix + "valoperpub"
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address.
	Bech32PrefixConsAddr = Bech32Prefix + "valcons"
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key.
	Bech32PrefixConsPub = Bech32Prefix + "valconspub"
)

// MantraCoinDenomRegex returns the mantra regex string
// this is used to override the default sdk coin denom regex
func MantraCoinDenomRegex() string {
	return `[a-zA-Z][a-zA-Z0-9/:._-]{1,127}`
}

func setupConfig() {
	// set the address prefixes
	config := sdk.GetConfig()
	SetAddressPrefixes(config)
	evmdcfg.SetBip44CoinType(config)
	config.Seal()
}

// SetAddressPrefixes builds the Config with Bech32 addressPrefix and publKeyPrefix for accounts, validators, and consensus nodes and verifies that addreeses have correct format.
func SetAddressPrefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(Bech32Prefix, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.SetAddressVerifier(wasmtypes.VerifyAddressLen())
}
