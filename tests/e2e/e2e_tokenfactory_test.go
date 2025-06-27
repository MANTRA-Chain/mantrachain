package e2e

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	subdenom = "test"
	mintAmt  = 1000
	burnAmt  = 800
)

const (
	proposalDisableDenomSendFilename = "proposal_disable_denom_send.json"
	proposalEnableDenomSendFilename  = "proposal_enable_denom_send.json"
)

func (s *IntegrationTestSuite) writeDisableDenomSendProposal(c *chain, denom string) {
	template := `
	{
		"messages": [
			{
				"@type": "/cosmos.bank.v1beta1.MsgSetSendEnabled",
				"authority": "%s",
				"send_enabled": [
					{
						"denom": "%s"
					}
				],
				"use_default_for": []
        	}
		],
		"metadata": "ipfs://CID",
		"deposit": "100uom",
		"title": "Disable %s for sending",
		"summary": "e2e-test disable token send"
	   }`
	propMsgBody := fmt.Sprintf(template,
		govAuthority,
		denom,
		denom,
	)

	err := writeFile(filepath.Join(c.validators[0].configDir(), "config", proposalDisableDenomSendFilename), []byte(propMsgBody))
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) writeEnableDenomSendProposal(c *chain, denom string) {
	template := `
	{
		"messages": [
			{
				"@type": "/cosmos.bank.v1beta1.MsgSetSendEnabled",
				"authority": "%s",
				"send_enabled": [],
				"use_default_for": ["%s"]
        	}
		],
		"metadata": "ipfs://CID",
		"deposit": "100uom",
		"title": "Reenable %s for sending",
		"summary": "e2e-test reenable token send"
	   }`
	propMsgBody := fmt.Sprintf(template,
		govAuthority,
		denom,
		denom,
	)

	err := writeFile(filepath.Join(c.validators[0].configDir(), "config", proposalEnableDenomSendFilename), []byte(propMsgBody))
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) testTokenfactoryCreate() {
	s.Run("create_denom_tokenfactory", func() {
		var (
			err           error
			valIdx        = 0
			c             = s.chainA
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)

		// define one sender and two recipient accounts
		alice, _ := c.genesisAccounts[1].keyInfo.GetAddress()

		var beforeAliceUomBalance,
			afterAliceUomBalance sdk.Coin

		denomCreationFee, err := queryTokenfactoryDenomCreationFee(chainEndpoint)
		s.Require().NoError(err)
		s.Require().Equal(denomCreationFee.Denom, uomDenom)

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				beforeAliceUomBalance, err = getSpecificBalance(chainEndpoint, alice.String(), uomDenom)
				s.Require().NoError(err)

				return beforeAliceUomBalance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		s.createDenom(c, valIdx, alice.String(), subdenom, standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterAliceUomBalance, err = getSpecificBalance(chainEndpoint, alice.String(), uomDenom)
				s.Require().NoError(err)

				decremented := beforeAliceUomBalance.Sub(denomCreationFee).Sub(standardFees).IsEqual(afterAliceUomBalance)

				return decremented
			},
			10*time.Second,
			5*time.Second,
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
		alice, _ := c.genesisAccounts[1].keyInfo.GetAddress()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()

		var beforeAliceCustomTokenBalance,
			afterAliceCustomTokenBalance,
			beforeBobCustomTokenBalance,
			afterBobCustomTokenBalance sdk.Coin

		customDenom := fmt.Sprintf("factory/%s/%s", alice.String(), subdenom)

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				beforeAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice.String(), customDenom)
				s.Require().NoError(err)

				beforeBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				return beforeAliceCustomTokenBalance.IsValid() && beforeBobCustomTokenBalance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		toMint := sdk.NewCoin(customDenom, math.NewInt(mintAmt))
		s.mintDenom(c, valIdx, alice.String(), toMint.String(), "", standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice.String(), customDenom)
				s.Require().NoError(err)

				incremented := beforeAliceCustomTokenBalance.Add(toMint).IsEqual(afterAliceCustomTokenBalance)

				return incremented
			},
			10*time.Second,
			5*time.Second,
		)

		s.mintDenom(c, valIdx, alice.String(), toMint.String(), bob.String(), standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				incremented := beforeBobCustomTokenBalance.Add(toMint).IsEqual(afterBobCustomTokenBalance)

				return incremented
			},
			10*time.Second,
			5*time.Second,
		)

		// disable token send for tokenfactory token should prevent further minting
		s.writeDisableDenomSendProposal(s.chainA, customDenom)
		proposalCounter++
		submitGovFlags := []string{configFile(proposalDisableDenomSendFilename)}
		depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
		voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

		validatorA := s.chainA.validators[0]
		validatorAddr, _ := validatorA.keyInfo.GetAddress()

		s.T().Logf("Proposal number: %d", proposalCounter)
		s.T().Logf("Submitting, deposit and vote Gov Proposal: Disable %s for sending", customDenom)
		s.submitGovProposal(chainEndpoint, validatorAddr.String(), proposalCounter, "banktypes.MsgSetSendEnabled", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

		s.Require().Eventually(
			func() bool {
				s.T().Logf("After MsgSetSendEnabled proposal to disable denom %s", customDenom)

				sendEnabled, err := querySendEnabled(chainEndpoint)
				s.Require().NoError(err)
				s.Require().Len(sendEnabled, 1)
				s.Require().Equal(customDenom, sendEnabled[0].Denom)

				return true
			},
			15*time.Second,
			5*time.Second,
		)

		s.mintDenom(c, valIdx, alice.String(), toMint.String(), bob.String(), standardFees.String(), true)

		escrowAddress, err := queryIBCEscrowAddress(chainEndpoint, "channel-0")
		s.Require().NoError(err)
		s.mintDenom(c, valIdx, alice.String(), toMint.String(), escrowAddress, standardFees.String(), true)
	})
}

func (s *IntegrationTestSuite) testTokenfactoryBurn() {
	var (
		err           error
		valIdx        = 0
		c             = s.chainA
		chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
	)
	// define one admin and one recipient
	alice, _ := c.genesisAccounts[1].keyInfo.GetAddress()
	bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()

	var beforeAliceCustomTokenBalance,
		afterAliceCustomTokenBalance,
		beforeBobCustomTokenBalance,
		afterBobCustomTokenBalance sdk.Coin

	customDenom := fmt.Sprintf("factory/%s/%s", alice.String(), subdenom)

	// get balances of sender and recipient accounts
	s.Require().Eventually(
		func() bool {
			beforeAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice.String(), customDenom)
			s.Require().NoError(err)

			beforeBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
			s.Require().NoError(err)

			return beforeAliceCustomTokenBalance.IsValid() && beforeBobCustomTokenBalance.IsValid()
		},
		10*time.Second,
		5*time.Second,
	)

	toBurn := sdk.NewCoin(customDenom, math.NewInt(burnAmt))

	s.Run("reenable_token_send", func() {
		s.burnDenom(c, valIdx, alice.String(), toBurn.String(), "", standardFees.String(), true)

		s.writeEnableDenomSendProposal(s.chainA, customDenom)

		proposalCounter++
		submitGovFlags := []string{configFile(proposalEnableDenomSendFilename)}
		depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
		voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

		validatorA := s.chainA.validators[0]
		validatorAddr, _ := validatorA.keyInfo.GetAddress()

		s.T().Logf("Proposal number: %d", proposalCounter)
		s.T().Logf("Submitting, deposit and vote Gov Proposal: Reenable %s for sending", customDenom)
		s.submitGovProposal(chainEndpoint, validatorAddr.String(), proposalCounter, "banktypes.MsgSetSendEnabled", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

		s.Require().Eventually(
			func() bool {
				s.T().Logf("After MsgSetSendEnabled proposal to reenable denom %s", customDenom)
				sendEnabled, err := querySendEnabled(chainEndpoint)
				s.Require().NoError(err)
				s.Require().Len(sendEnabled, 0)
				return true
			},
			15*time.Second,
			5*time.Second,
		)
	})

	s.Run("burn_tokens_tokenfactory", func() {
		s.burnDenom(c, valIdx, alice.String(), toBurn.String(), "", standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice.String(), customDenom)
				s.Require().NoError(err)

				decremented := beforeAliceCustomTokenBalance.Sub(toBurn).IsEqual(afterAliceCustomTokenBalance)

				return decremented
			},
			10*time.Second,
			5*time.Second,
		)

		s.burnDenom(c, valIdx, alice.String(), toBurn.String(), bob.String(), standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				decremented := beforeBobCustomTokenBalance.Sub(toBurn).IsEqual(afterBobCustomTokenBalance)

				return decremented
			},
			10*time.Second,
			5*time.Second,
		)

		escrowAddress, err := queryIBCEscrowAddress(chainEndpoint, "channel-0")
		s.Require().NoError(err)
		s.burnDenom(c, valIdx, alice.String(), toBurn.String(), escrowAddress, standardFees.String(), true)
	})
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
	denom := fmt.Sprintf("factory/%s/%s", sender, subdenom)
	s.T().Logf("creating tokenfactory denom %s", denom)
	if expErr {
		s.executeTxCommand(ctx, c, ibcCmd, valIdx, s.expectErrExecValidation(c, valIdx, true))
		s.T().Log("create tokenfactory denom unsuccessful")
	} else {
		s.executeTxCommand(ctx, c, ibcCmd, valIdx, s.defaultExecValidation(c, valIdx))
		s.T().Log("successfully created tokenfactory denom")
	}
}

//nolint:unparam
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
		mintTo = sender
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
		burnFrom = sender
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
