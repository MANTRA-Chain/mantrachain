package e2e

import (
	"context"
	"fmt"
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

				gasFeesBurntMax := standardFees.Sub(sdk.NewCoin(uomDenom, math.NewInt(1000)))
				gasFeesBurntMin := standardFees.Sub(sdk.NewCoin(uomDenom, math.NewInt(1500)))
				maxDecremented := beforeAliceUomBalance.Sub(denomCreationFee).Sub(gasFeesBurntMin).IsGTE(afterAliceUomBalance)
				minDecremented := beforeAliceUomBalance.Sub(denomCreationFee).Sub(gasFeesBurntMax).IsLTE(afterAliceUomBalance)

				return minDecremented && maxDecremented
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

		escrowAddress, err := queryIBCEscrowAddress(chainEndpoint, "channel-0")
		s.Require().NoError(err)
		s.mintDenom(c, valIdx, alice.String(), toMint.String(), escrowAddress, standardFees.String(), true)
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
