package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	subdenom = "test"
	mintAmt  = 1000
	burnAmt  = 800
)

func (s *IntegrationTestSuite) testTokenfactoryCreate() {
	s.Run("create_denom_tokenfactory", func() {
		var (
			err           error
			valIdx        = 0
			c             = s.chainA
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)

		// define one sender and two recipient accounts
		alice := s.getAlice()

		var beforeAliceUomBalance,
			afterAliceUomBalance sdk.Coin

		denomCreationFee, err := queryTokenfactoryDenomCreationFee(chainEndpoint)
		s.Require().NoError(err)
		s.Require().Equal(denomCreationFee.Denom, uomDenom)

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				beforeAliceUomBalance, err = getSpecificBalance(chainEndpoint, alice, uomDenom)
				s.Require().NoError(err)

				return beforeAliceUomBalance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		s.createDenom(c, valIdx, alice, subdenom, standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterAliceUomBalance, err = getSpecificBalance(chainEndpoint, alice, uomDenom)
				s.Require().NoError(err)

				beforeAlice := beforeAliceUomBalance.Sub(denomCreationFee).Sub(standardFees)

				return beforeAlice.Equal(afterAliceUomBalance)
			},
			10*time.Second,
			5*time.Second,
		)
	})
}

func (s *IntegrationTestSuite) testTokenfactoryAdmin() {
	s.Run("default_denom_admin_should_be_the_creator", func() {
		var (
			err    error
			valIdx = 0
			c      = s.chainA
		)

		chainEndpoint := fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))

		result, err := queryTokenfactoryDenomAuthorityMetadata(chainEndpoint, s.getAlice(), subdenom)

		s.Require().NoError(err)

		s.Require().Equal(result.Admin, s.getAlice(), "By default, the denom admin should be the creator")
	})
}

func (s *IntegrationTestSuite) testTokenfactorySetMetadata() {
	s.Run("set_denom_metadata_tokenfactory", func() {
		var (
			err    error
			valIdx = 0
			c      = s.chainA
		)

		// Create metadata JSON content using the global helper function
		metadataContent := s.BuildTokenMetadata()
		customDenom := metadataContent.Base
		metadataString := MetadataToString(metadataContent)

		// Write metadata to file in the validator's config directory
		metadataFileName := "metadata.json"
		metadataFile := filepath.Join(c.validators[valIdx].configDir(), metadataFileName)
		err = writeFile(metadataFile, []byte(metadataString))
		s.Require().NoError(err)
		s.T().Logf("Start setting metadata for denom %s", customDenom)

		// Set the metadata using the CLI command
		s.setDenomMetadata(c, valIdx, s.getAlice(), filepath.Join(mantraHomePath, metadataFileName), standardFees.String(), false)

		s.T().Logf("Successfully set metadata for denom %s", customDenom)

		// Query and verify the metadata was set correctly
		chainEndpoint := fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))

		s.Require().Eventually(
			func() bool {
				queriedMetadata, err := queryTokenfactoryDenomMetadata(chainEndpoint, customDenom)
				if err != nil {
					s.T().Logf("Error querying metadata: %v", err)
					return false
				}

				// Verify metadata fields
				if queriedMetadata.Description != metadataContent.Description {
					s.T().Logf("Description mismatch: expected %s, got %s", metadataContent.Description, queriedMetadata.Description)
					return false
				}

				if queriedMetadata.Base != metadataContent.Base {
					s.T().Logf("Base mismatch: expected %s, got %s", metadataContent.Base, queriedMetadata.Base)
					return false
				}

				if queriedMetadata.Display != metadataContent.Display {
					s.T().Logf("Display mismatch: expected %s, got %s", metadataContent.Display, queriedMetadata.Display)
					return false
				}

				if queriedMetadata.Name != metadataContent.Name {
					s.T().Logf("Name mismatch: expected %s, got %s", metadataContent.Name, queriedMetadata.Name)
					return false
				}

				if queriedMetadata.Symbol != metadataContent.Symbol {
					s.T().Logf("Symbol mismatch: expected %s, got %s", metadataContent.Symbol, queriedMetadata.Symbol)
					return false
				}

				if len(queriedMetadata.DenomUnits) != len(metadataContent.DenomUnits) {
					s.T().Logf("DenomUnits length mismatch: expected %d, got %d", len(metadataContent.DenomUnits), len(queriedMetadata.DenomUnits))
					return false
				}

				// Verify denom units
				for i, expectedUnit := range metadataContent.DenomUnits {
					if i >= len(queriedMetadata.DenomUnits) {
						s.T().Logf("Missing denom unit at index %d", i)
						return false
					}

					queriedUnit := queriedMetadata.DenomUnits[i]
					if queriedUnit.Denom != expectedUnit.Denom {
						s.T().Logf("DenomUnit[%d] Denom mismatch: expected %s, got %s", i, expectedUnit.Denom, queriedUnit.Denom)
						return false
					}

					if queriedUnit.Exponent != expectedUnit.Exponent {
						s.T().Logf("DenomUnit[%d] Exponent mismatch: expected %d, got %d", i, expectedUnit.Exponent, queriedUnit.Exponent)
						return false
					}

					if len(queriedUnit.Aliases) != len(expectedUnit.Aliases) {
						s.T().Logf("DenomUnit[%d] Aliases length mismatch: expected %d, got %d", i, len(expectedUnit.Aliases), len(queriedUnit.Aliases))
						return false
					}

					for j, expectedAlias := range expectedUnit.Aliases {
						if j >= len(queriedUnit.Aliases) || queriedUnit.Aliases[j] != expectedAlias {
							s.T().Logf("DenomUnit[%d] Alias[%d] mismatch: expected %s, got %s", i, j, expectedAlias, queriedUnit.Aliases[j])
							return false
						}
					}
				}

				s.T().Logf("Successfully verified metadata for denom %s", customDenom)
				return true
			},
			30*time.Second,
			2*time.Second,
		)
	})
}

func (s *IntegrationTestSuite) testTokenfactoryMint() {
	s.Run("mint_tokens_tokenfactory", func() {
		var (
			err           error
			valIdx        = 0
			c             = s.chainA
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)

		// define one admin and one recipient
		alice := s.getAlice()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()

		var beforeAliceCustomTokenBalance,
			afterAliceCustomTokenBalance,
			beforeBobCustomTokenBalance,
			afterBobCustomTokenBalance sdk.Coin

		customDenom := buildDenom(alice, subdenom)

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				beforeAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice, customDenom)
				s.Require().NoError(err)

				beforeBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				return beforeAliceCustomTokenBalance.IsValid() && beforeBobCustomTokenBalance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		toMint := sdk.NewCoin(customDenom, math.NewInt(mintAmt))
		s.mintDenom(c, valIdx, alice, toMint.String(), "", standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice, customDenom)
				s.Require().NoError(err)

				incrementedAlice := beforeAliceCustomTokenBalance.Add(toMint)

				return incrementedAlice.Equal(afterAliceCustomTokenBalance)
			},
			10*time.Second,
			5*time.Second,
		)

		s.mintDenom(c, valIdx, alice, toMint.String(), bob.String(), standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				incrementedBob := beforeBobCustomTokenBalance.Add(toMint)

				return incrementedBob.Equal(afterBobCustomTokenBalance)
			},
			10*time.Second,
			5*time.Second,
		)

		if !s.testOnSingleNode {
			// No need to test with IBC on single-node setup
			escrowAddress, err := queryIBCEscrowAddress(chainEndpoint, "channel-0")
			s.Require().NoError(err)
			s.mintDenom(c, valIdx, alice, toMint.String(), escrowAddress, standardFees.String(), true)
		}
	})
}

func (s *IntegrationTestSuite) testTokenfactoryBurn() {
	s.Run("burn_tokens_tokenfactory", func() {
		var (
			err           error
			valIdx        = 0
			c             = s.chainA
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)

		// define one admin and one recipient
		alice := s.getAlice()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()

		var beforeAliceCustomTokenBalance,
			afterAliceCustomTokenBalance,
			beforeBobCustomTokenBalance,
			afterBobCustomTokenBalance sdk.Coin

		customDenom := buildDenom(alice, subdenom)

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				beforeAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice, customDenom)
				s.Require().NoError(err)

				beforeBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				return beforeAliceCustomTokenBalance.IsValid() && beforeBobCustomTokenBalance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		toBurn := sdk.NewCoin(customDenom, math.NewInt(burnAmt))
		s.burnDenom(c, valIdx, alice, toBurn.String(), "", standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice, customDenom)
				s.Require().NoError(err)

				beforeAlice := beforeAliceCustomTokenBalance.Sub(toBurn)

				return beforeAlice.Equal(afterAliceCustomTokenBalance)
			},
			10*time.Second,
			5*time.Second,
		)

		s.burnDenom(c, valIdx, alice, toBurn.String(), bob.String(), standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				beforeBob := beforeBobCustomTokenBalance.Sub(toBurn)

				return beforeBob.Equal(afterBobCustomTokenBalance)
			},
			10*time.Second,
			5*time.Second,
		)

		if !s.testOnSingleNode {
			// No need to test with IBC on single-node setup
			escrowAddress, err := queryIBCEscrowAddress(chainEndpoint, "channel-0")
			s.Require().NoError(err)
			s.burnDenom(c, valIdx, alice, toBurn.String(), escrowAddress, standardFees.String(), true)
		}
	})
}

func buildDenom(sender, subDenom string) string {
	return fmt.Sprintf("factory/%s/%s", sender, subDenom)
}

func (s *IntegrationTestSuite) getAlice() string {
	alice, err := s.chainA.genesisAccounts[1].keyInfo.GetAddress()
	s.Require().NoError(err)
	return alice.String()
}

func (s *IntegrationTestSuite) BuildTokenMetadata() banktypes.Metadata {
	factoryDenom := buildDenom(s.getAlice(), subdenom)

	symbol := cases.Upper(language.English).String(subdenom)
	name := cases.Title(language.English).String(subdenom)
	metadata := banktypes.Metadata{
		Description: fmt.Sprintf("%s token for tokenfactory e2e tests", name),
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    factoryDenom,
				Exponent: 0,
				Aliases:  []string{"test_alias"},
			},
			{
				Denom:    symbol,
				Exponent: 6,
				Aliases:  []string{},
			},
		},
		Base:    factoryDenom,
		Display: symbol,
		Name:    fmt.Sprintf("%s Token", name),
		Symbol:  symbol,
	}

	return metadata
}

func MetadataToString(metadata banktypes.Metadata) string {
	metadataBytes, err := json.MarshalIndent(metadata, "", "\t")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal metadata: %v", err))
	}
	return string(metadataBytes)
}

func (s *IntegrationTestSuite) createDenom(c *chain, valIdx int, sender, subdenom, fees string, expErr bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ibcCmd := []string{
		mantrachaindBinary,
		txCommand,
		"tokenfactory",
		"create-denom",
		subdenom,
		fmt.Sprintf("--from=%s", sender),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--broadcast-mode=sync",
		"--output=json",
		"-y",
	}
	denom := buildDenom(sender, subdenom)
	s.T().Logf("%s is creating tokenfactory denom %s", sender, denom)
	if expErr {
		s.executeTxCommand(ctx, c, ibcCmd, valIdx, s.expectErrExecValidation(c, valIdx, true))
		s.T().Log("create tokenfactory denom unsuccessful")
	} else {
		s.executeTxCommand(ctx, c, ibcCmd, valIdx, s.defaultExecValidation(c, valIdx))
		s.T().Log("successfully created tokenfactory denom")
	}
}

func (s *IntegrationTestSuite) setDenomMetadata(c *chain, valIdx int, sender, metadataFile, fees string, expErr bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Sample tx: https://mantrascan.io/dukong/tx/4f40cc08aadb5ca005a4138353c707b1398858c577186458d5ce2a70bd3a67c8
	cmd := []string{
		mantrachaindBinary,
		txCommand,
		"tokenfactory",
		"set-denom-metadata",
		metadataFile,
		fmt.Sprintf("--from=%s", sender),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.5"),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--broadcast-mode=sync",
		"--output=json",
		"-y",
	}

	s.T().Logf("Address %s is setting denom metadata from file %s", sender, metadataFile)
	if expErr {
		s.executeTxCommand(ctx, c, cmd, valIdx, s.expectErrExecValidation(c, valIdx, true))
		s.T().Log("set denom metadata unsuccessful")
	} else {
		s.executeTxCommand(ctx, c, cmd, valIdx, s.defaultExecValidation(c, valIdx))
		s.T().Log("successfully set denom metadata")
	}
}

func (s *IntegrationTestSuite) mintDenom(c *chain, valIdx int, sender, mintCoin, mintTo, fees string, expErr bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var ibcCmd []string
	if mintTo == "" {
		ibcCmd = []string{
			mantrachaindBinary,
			txCommand,
			"tokenfactory",
			"mint",
			mintCoin,
			fmt.Sprintf("--from=%s", sender),
			fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
			fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.5"),
			fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
			"--keyring-backend=test",
			"--broadcast-mode=sync",
			"--output=json",
			"-y",
		}
	} else {
		ibcCmd = []string{
			mantrachaindBinary,
			txCommand,
			"tokenfactory",
			"mint",
			mintCoin,
			mintTo,
			fmt.Sprintf("--from=%s", sender),
			fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
			fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.5"),
			fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
			"--keyring-backend=test",
			"--broadcast-mode=sync",
			"--output=json",
			"-y",
		}
	}

	s.T().Logf("minting %s to %s", mintCoin, mintTo)
	if expErr {
		s.executeTxCommand(ctx, c, ibcCmd, valIdx, s.expectErrExecValidation(c, valIdx, true))
		s.T().Log("unsuccessful minting of tokenfactory denom")
	} else {
		s.executeTxCommand(ctx, c, ibcCmd, valIdx, s.defaultExecValidation(c, valIdx))
		s.T().Log("successfully minted tokenfactory denom")
	}
}

func (s *IntegrationTestSuite) burnDenom(c *chain, valIdx int, sender, burnCoin, burnFrom, fees string, expErr bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var ibcCmd []string
	if burnFrom == "" {
		ibcCmd = []string{
			mantrachaindBinary,
			txCommand,
			"tokenfactory",
			"burn",
			burnCoin,
			fmt.Sprintf("--from=%s", sender),
			fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
			fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.5"),
			fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
			"--keyring-backend=test",
			"--broadcast-mode=sync",
			"--output=json",
			"-y",
		}
	} else {
		ibcCmd = []string{
			mantrachaindBinary,
			txCommand,
			"tokenfactory",
			"burn",
			burnCoin,
			burnFrom,
			fmt.Sprintf("--from=%s", sender),
			fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
			fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.5"),
			fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
			"--keyring-backend=test",
			"--broadcast-mode=sync",
			"--output=json",
			"-y",
		}
	}

	s.T().Logf("burning %s from %s", burnCoin, burnFrom)
	if expErr {
		s.executeTxCommand(ctx, c, ibcCmd, valIdx, s.expectErrExecValidation(c, valIdx, true))
		s.T().Log("unsuccessful burning of tokenfactory denom")
	} else {
		s.executeTxCommand(ctx, c, ibcCmd, valIdx, s.defaultExecValidation(c, valIdx))
		s.T().Log("successfully burned tokenfactory denom")
	}
}
