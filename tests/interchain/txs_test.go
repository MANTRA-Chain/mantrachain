package interchain_test

import (
	"context"
	"fmt"
	"path"
	"testing"
	"time"

	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	"github.com/MANTRA-Chain/mantrachain/v5/tests/interchain/chainsuite"
)

const txAmount = 1_000_000_000

type TxSuite struct {
	*chainsuite.Suite
}

func txAmountUom() string {
	return fmt.Sprintf("%d%s", txAmount, chainsuite.Uom)
}

func (s *TxSuite) TestBankSend() {
	balanceBefore, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[1].Address, chainsuite.Uom)
	s.Require().NoError(err)

	_, err = s.Chain.Validators[0].ExecTx(
		s.GetContext(),
		s.Chain.ValidatorWallets[0].Moniker,
		"bank", "send",
		s.Chain.ValidatorWallets[0].Address, s.Chain.ValidatorWallets[1].Address, txAmountUom(),
	)
	s.Require().NoError(err)

	balanceAfter, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[1].Address, chainsuite.Uom)
	s.Require().NoError(err)
	s.Require().Equal(balanceBefore.Add(sdkmath.NewInt(txAmount)), balanceAfter)
}

func (s TxSuite) TestDelegateWithdrawUnbond() {
	// delegate tokens
	_, err := s.Chain.Validators[0].ExecTx(
		s.GetContext(),
		s.Chain.ValidatorWallets[0].Moniker,
		"staking", "delegate", s.Chain.ValidatorWallets[0].ValoperAddress, txAmountUom(),
	)
	s.Require().NoError(err)

	startingBalance, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[0].Address, chainsuite.Uom)
	s.Require().NoError(err)
	time.Sleep(20 * time.Second)
	// Withdraw rewards
	_, err = s.Chain.Validators[0].ExecTx(
		s.GetContext(),
		s.Chain.ValidatorWallets[0].Moniker,
		"distribution", "withdraw-rewards", s.Chain.ValidatorWallets[0].ValoperAddress,
	)
	s.Require().NoError(err)
	endingBalance, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[0].Address, chainsuite.Uom)
	s.Require().NoError(err)
	s.Require().Truef(endingBalance.GT(startingBalance), "endingBalance: %s, startingBalance: %s", endingBalance, startingBalance)

	// Unbond tokens
	_, err = s.Chain.Validators[0].ExecTx(
		s.GetContext(),
		s.Chain.ValidatorWallets[0].Moniker,
		"staking", "unbond", s.Chain.ValidatorWallets[0].ValoperAddress, txAmountUom(),
	)
	s.Require().NoError(err)
}

func (s TxSuite) TestAuthz() {
	s.Run("send", func() {
		balanceBefore, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[2].Address, chainsuite.Uom)
		s.Require().NoError(err)
		_, err = s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"authz", "grant", s.Chain.ValidatorWallets[1].Address, "send",
			"--spend-limit", fmt.Sprintf("%d%s", txAmount*2, chainsuite.Uom),
			"--allow-list", s.Chain.ValidatorWallets[2].Address,
		)
		s.Require().NoError(err)

		s.Require().Error(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "bank", "send", s.Chain.ValidatorWallets[0].Address, s.Chain.ValidatorWallets[3].Address, txAmountUom()))

		s.Require().NoError(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "bank", "send", s.Chain.ValidatorWallets[0].Address, s.Chain.ValidatorWallets[2].Address, txAmountUom()))
		balanceAfter, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[2].Address, chainsuite.Uom)
		s.Require().NoError(err)
		s.Require().Equal(balanceBefore.Add(sdkmath.NewInt(int64(txAmount))), balanceAfter)

		s.Require().Error(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "bank", "send", s.Chain.ValidatorWallets[0].Address, s.Chain.ValidatorWallets[2].Address, fmt.Sprintf("%d%s", txAmount+200, chainsuite.Uom)))

		_, err = s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"authz", "revoke", s.Chain.ValidatorWallets[1].Address, "/cosmos.bank.v1beta1.MsgSend",
		)
		s.Require().NoError(err)

		s.Require().Error(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "bank", "send", s.Chain.ValidatorWallets[0].Address, s.Chain.ValidatorWallets[2].Address, txAmountUom()))
	})

	s.Run("delegate", func() {
		_, err := s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"authz", "grant", s.Chain.ValidatorWallets[1].Address, "delegate",
			"--allowed-validators", s.Chain.ValidatorWallets[2].ValoperAddress,
		)
		s.Require().NoError(err)

		s.Require().NoError(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "staking", "delegate", s.Chain.ValidatorWallets[2].ValoperAddress, txAmountUom(), "--from", s.Chain.ValidatorWallets[0].Address))

		s.Require().Error(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "staking", "delegate", s.Chain.ValidatorWallets[0].ValoperAddress, txAmountUom(), "--from", s.Chain.ValidatorWallets[0].Address))

		_, err = s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"authz", "revoke", s.Chain.ValidatorWallets[1].Address, "/cosmos.staking.v1beta1.MsgDelegate",
		)
		s.Require().NoError(err)
		s.Require().Error(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "staking", "delegate", s.Chain.ValidatorWallets[2].ValoperAddress, txAmountUom(), "--from", s.Chain.ValidatorWallets[0].Address))
	})

	s.Run("unbond", func() {
		valHex, err := s.Chain.GetValidatorHex(s.GetContext(), 2)
		s.Require().NoError(err)
		powerBefore, err := s.Chain.GetValidatorPower(s.GetContext(), valHex)
		s.Require().NoError(err)
		_, err = s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"staking", "delegate", s.Chain.ValidatorWallets[2].ValoperAddress, txAmountUom(),
		)
		s.Require().NoError(err)
		s.Require().EventuallyWithT(func(c *assert.CollectT) {
			powerAfter, err := s.Chain.GetValidatorPower(s.GetContext(), valHex)
			s.Require().NoError(err)
			assert.NoError(c, err)
			assert.Greater(c, powerAfter, powerBefore)
		}, 15*chainsuite.CommitTimeout, chainsuite.CommitTimeout)

		_, err = s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"authz", "grant", s.Chain.ValidatorWallets[1].Address, "unbond",
			"--allowed-validators", s.Chain.ValidatorWallets[2].ValoperAddress,
		)
		s.Require().NoError(err)

		s.Require().NoError(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "staking", "unbond", s.Chain.ValidatorWallets[2].ValoperAddress, txAmountUom(), "--from", s.Chain.ValidatorWallets[0].Address))
		s.Require().Error(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "staking", "unbond", s.Chain.ValidatorWallets[0].ValoperAddress, txAmountUom(), "--from", s.Chain.ValidatorWallets[0].Address))

		s.Require().EventuallyWithT(func(c *assert.CollectT) {
			powerAfter, err := s.Chain.GetValidatorPower(s.GetContext(), valHex)
			s.Require().NoError(err)
			assert.NoError(c, err)
			assert.Equal(c, powerAfter, powerBefore)
		}, 15*chainsuite.CommitTimeout, chainsuite.CommitTimeout)

		_, err = s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"authz", "revoke", s.Chain.ValidatorWallets[1].Address, "/cosmos.staking.v1beta1.MsgUndelegate",
		)
		s.Require().NoError(err)
		s.Require().Error(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "staking", "unbond", s.Chain.ValidatorWallets[2].ValoperAddress, txAmountUom(), "--from", s.Chain.ValidatorWallets[0].Address))
	})

	s.Run("redelegate", func() {
		val0Hex, err := s.Chain.GetValidatorHex(s.GetContext(), 0)
		s.Require().NoError(err)
		val2Hex, err := s.Chain.GetValidatorHex(s.GetContext(), 2)
		s.Require().NoError(err)
		val0PowerBefore, err := s.Chain.GetValidatorPower(s.GetContext(), val0Hex)
		s.Require().NoError(err)
		_, err = s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"staking", "delegate", s.Chain.ValidatorWallets[0].ValoperAddress, txAmountUom(),
		)
		s.Require().NoError(err)
		s.Require().EventuallyWithT(func(c *assert.CollectT) {
			val0PowerAfter, err := s.Chain.GetValidatorPower(s.GetContext(), val0Hex)
			s.Require().NoError(err)
			s.Require().NoError(err)
			s.Require().Greater(val0PowerAfter, val0PowerBefore)
		}, 15*chainsuite.CommitTimeout, chainsuite.CommitTimeout)

		_, err = s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"authz", "grant", s.Chain.ValidatorWallets[1].Address, "redelegate",
			"--allowed-validators", s.Chain.ValidatorWallets[2].ValoperAddress,
		)
		s.Require().NoError(err)

		s.Require().Error(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "staking", "redelegate", s.Chain.ValidatorWallets[0].ValoperAddress, s.Chain.ValidatorWallets[1].ValoperAddress, txAmountUom(), "--from", s.Chain.ValidatorWallets[0].Address))

		val2PowerBefore, err := s.Chain.GetValidatorPower(s.GetContext(), val2Hex)
		s.Require().NoError(err)
		s.Require().NoError(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "staking", "redelegate", s.Chain.ValidatorWallets[0].ValoperAddress, s.Chain.ValidatorWallets[2].ValoperAddress, txAmountUom(), "--from", s.Chain.ValidatorWallets[0].Address))
		s.Require().EventuallyWithT(func(c *assert.CollectT) {
			val2PowerAfter, err := s.Chain.GetValidatorPower(s.GetContext(), val2Hex)
			s.Require().NoError(err)
			s.Require().Greater(val2PowerAfter, val2PowerBefore)
		}, 15*chainsuite.CommitTimeout, chainsuite.CommitTimeout)

		_, err = s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"authz", "revoke", s.Chain.ValidatorWallets[1].Address, "/cosmos.staking.v1beta1.MsgBeginRedelegate",
		)
		s.Require().NoError(err)

		s.Require().Error(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "staking", "redelegate", s.Chain.ValidatorWallets[0].ValoperAddress, s.Chain.ValidatorWallets[2].ValoperAddress, txAmountUom(), "--from", s.Chain.ValidatorWallets[0].Address))
	})

	s.Run("generic", func() {
		_, err := s.Chain.Validators[0].ExecTx(
			s.GetContext(),
			s.Chain.ValidatorWallets[0].Moniker,
			"authz", "grant", s.Chain.ValidatorWallets[1].Address, "generic",
			"--msg-type", "/cosmos.gov.v1.MsgVote",
		)
		s.Require().NoError(err)

		prop, err := s.Chain.BuildProposal(nil, "Test Proposal", "Test Proposal", "ipfs://CID", chainsuite.GovDepositAmount, s.Chain.ValidatorWallets[0].ValoperAddress, false)
		s.Require().NoError(err)
		result, err := s.Chain.SubmitProposal(s.GetContext(), s.Chain.ValidatorWallets[0].Moniker, prop)
		s.Require().NoError(err)
		s.Require().NoError(s.authzGenExec(s.GetContext(), s.Chain.ValidatorWallets[1], "gov", "vote", result.ProposalID, "yes", "--from", s.Chain.ValidatorWallets[0].Address))
	})
}

func (s TxSuite) TestFeegrant() {
	const (
		granter       = 5
		grantee       = 1
		fundsReceiver = 2
	)

	tests := []struct {
		name   string
		revoke func(expireTime time.Time)
	}{
		{
			name: "revoke",
			revoke: func(_ time.Time) {
				_, err := s.Chain.Validators[granter].ExecTx(
					s.GetContext(),
					s.Chain.ValidatorWallets[granter].Moniker,
					"feegrant", "revoke", s.Chain.ValidatorWallets[granter].Address, s.Chain.ValidatorWallets[grantee].Address,
				)
				s.Require().NoError(err)
			},
		},
		{
			name: "expire",
			revoke: func(expire time.Time) {
				<-time.After(time.Until(expire))
				err := testutil.WaitForBlocks(s.GetContext(), 1, s.Chain)
				s.Require().NoError(err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			expire := time.Now().Add(20 * chainsuite.CommitTimeout)
			_, err := s.Chain.Validators[granter].ExecTx(
				s.GetContext(),
				s.Chain.ValidatorWallets[granter].Moniker,
				"feegrant", "grant", s.Chain.ValidatorWallets[granter].Address, s.Chain.ValidatorWallets[grantee].Address,
				"--expiration", expire.Format(time.RFC3339),
			)
			s.Require().NoError(err)

			granterBalanceBefore, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[granter].Address, chainsuite.Uom)
			s.Require().NoError(err)
			granteeBalanceBefore, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[grantee].Address, chainsuite.Uom)
			s.Require().NoError(err)

			_, err = s.Chain.Validators[grantee].ExecTx(s.GetContext(), s.Chain.ValidatorWallets[grantee].Moniker,
				"bank", "send", s.Chain.ValidatorWallets[grantee].Address, s.Chain.ValidatorWallets[fundsReceiver].Address, txAmountUom(),
				"--fee-granter", s.Chain.ValidatorWallets[granter].Address,
			)
			s.Require().NoError(err)

			granteeBalanceAfter, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[grantee].Address, chainsuite.Uom)
			s.Require().NoError(err)
			granterBalanceAfter, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[granter].Address, chainsuite.Uom)
			s.Require().NoError(err)

			s.Require().True(granterBalanceAfter.LT(granterBalanceBefore), "granterBalanceBefore: %s, granterBalanceAfter: %s", granterBalanceBefore, granterBalanceAfter)
			s.Require().True(granteeBalanceAfter.Equal(granteeBalanceBefore.Sub(sdkmath.NewInt(txAmount))), "granteeBalanceBefore: %s, granteeBalanceAfter: %s", granteeBalanceBefore, granteeBalanceAfter)

			tt.revoke(expire)

			_, err = s.Chain.Validators[1].ExecTx(s.GetContext(), s.Chain.ValidatorWallets[grantee].Moniker,
				"bank", "send", s.Chain.ValidatorWallets[1].Address, s.Chain.ValidatorWallets[fundsReceiver].Address, txAmountUom(),
				"--fee-granter", s.Chain.ValidatorWallets[0].Address,
			)
			s.Require().Error(err)
		})
	}
}

func (s *TxSuite) TestMultisig() {
	pubkey1, _, err := s.Chain.Validators[1].ExecBin(s.GetContext(), "keys", "show", s.Chain.ValidatorWallets[1].Moniker, "--pubkey", "--keyring-backend", "test")
	s.Require().NoError(err)

	pubkey2, _, err := s.Chain.Validators[2].ExecBin(s.GetContext(), "keys", "show", s.Chain.ValidatorWallets[2].Moniker, "--pubkey", "--keyring-backend", "test")
	s.Require().NoError(err)

	_, _, err = s.Chain.Validators[0].ExecBin(s.GetContext(), "keys", "add", "val1", "--pubkey", string(pubkey1), "--keyring-backend", "test")
	s.Require().NoError(err)

	_, _, err = s.Chain.Validators[0].ExecBin(s.GetContext(), "keys", "add", "val2", "--pubkey", string(pubkey2), "--keyring-backend", "test")
	s.Require().NoError(err)

	multisigName := "multisig"
	_, _, err = s.Chain.Validators[0].ExecBin(s.GetContext(), "keys", "add", multisigName, "--multisig", fmt.Sprintf("%s,val1,val2", s.Chain.ValidatorWallets[0].Moniker), "--multisig-threshold", "2", "--keyring-backend", "test")
	s.Require().NoError(err)

	pubkeyMulti, _, err := s.Chain.Validators[0].ExecBin(s.GetContext(), "keys", "show", multisigName, "--pubkey", "--keyring-backend", "test")
	s.Require().NoError(err)

	_, _, err = s.Chain.Validators[1].ExecBin(s.GetContext(), "keys", "add", multisigName, "--pubkey", string(pubkeyMulti), "--keyring-backend", "test")
	s.Require().NoError(err)
	_, _, err = s.Chain.Validators[2].ExecBin(s.GetContext(), "keys", "add", multisigName, "--pubkey", string(pubkeyMulti), "--keyring-backend", "test")
	s.Require().NoError(err)
	// bogus validator, not in the multisig
	_, _, err = s.Chain.Validators[4].ExecBin(s.GetContext(), "keys", "add", multisigName, "--pubkey", string(pubkeyMulti), "--keyring-backend", "test")
	s.Require().NoError(err)

	defer func() {
		_, _, err = s.Chain.Validators[0].ExecBin(s.GetContext(), "keys", "delete", "val1", "--keyring-backend", "test", "-y")
		s.Require().NoError(err)
		_, _, err = s.Chain.Validators[0].ExecBin(s.GetContext(), "keys", "delete", "val2", "--keyring-backend", "test", "-y")
		s.Require().NoError(err)
		for i := 0; i < 3; i++ {
			_, _, err = s.Chain.Validators[i].ExecBin(s.GetContext(), "keys", "delete", multisigName, "--keyring-backend", "test", "-y")
			s.Require().NoError(err)
		}
		_, _, err = s.Chain.Validators[4].ExecBin(s.GetContext(), "keys", "delete", multisigName, "--keyring-backend", "test", "-y")
		s.Require().NoError(err)
	}()

	multisigAddr, err := s.Chain.Validators[0].KeyBech32(s.GetContext(), multisigName, "")
	s.Require().NoError(err)

	err = s.Chain.SendFunds(s.GetContext(), interchaintest.FaucetAccountKeyName, ibc.WalletAmount{
		Denom:   chainsuite.Uom,
		Amount:  sdkmath.NewInt(chainsuite.ValidatorFunds),
		Address: multisigAddr,
	})
	s.Require().NoError(err)

	balanceBefore, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[3].Address, chainsuite.Uom)
	s.Require().NoError(err)

	txjson, err := s.Chain.GenerateTx(
		s.GetContext(), 0, "bank", "send", multisigName, s.Chain.ValidatorWallets[3].Address, txAmountUom(),
		"--gas", "auto", "--gas-adjustment", fmt.Sprint(s.Chain.Config().GasAdjustment), "--gas-prices", s.Chain.Config().GasPrices,
	)
	s.Require().NoError(err)

	err = s.Chain.Validators[0].WriteFile(s.GetContext(), []byte(txjson), "tx.json")
	s.Require().NoError(err)

	signed0, _, err := s.Chain.Validators[0].Exec(s.GetContext(),
		s.Chain.Validators[0].TxCommand(s.Chain.ValidatorWallets[0].Moniker, "sign",
			path.Join(s.Chain.Validators[0].HomeDir(), "tx.json"),
			"--multisig", multisigAddr,
		), nil)
	s.Require().NoError(err)

	err = s.Chain.Validators[1].WriteFile(s.GetContext(), []byte(txjson), "tx.json")
	s.Require().NoError(err)

	signed1, _, err := s.Chain.Validators[1].Exec(s.GetContext(),
		s.Chain.Validators[1].TxCommand(s.Chain.ValidatorWallets[1].Moniker, "sign",
			path.Join(s.Chain.Validators[1].HomeDir(), "tx.json"),
			"--multisig", multisigAddr,
		), nil)
	s.Require().NoError(err)

	err = s.Chain.Validators[4].WriteFile(s.GetContext(), []byte(txjson), "tx.json")
	s.Require().NoError(err)
	_, _, err = s.Chain.Validators[4].Exec(s.GetContext(),
		s.Chain.Validators[4].TxCommand(s.Chain.ValidatorWallets[4].Moniker, "sign",
			path.Join(s.Chain.Validators[4].HomeDir(), "tx.json"),
			"--multisig", multisigAddr,
		), nil)
	s.Require().Error(err)

	err = s.Chain.Validators[0].WriteFile(s.GetContext(), signed0, "signed0.json")
	s.Require().NoError(err)

	_, _, err = s.Chain.Validators[0].Exec(s.GetContext(), s.Chain.Validators[0].TxCommand(
		multisigName,
		"multisign",
		path.Join(s.Chain.Validators[0].HomeDir(), "tx.json"),
		multisigName,
		path.Join(s.Chain.Validators[0].HomeDir(), "signed0.json"),
	), nil)
	s.Require().NoError(err)

	err = s.Chain.Validators[0].WriteFile(s.GetContext(), signed1, "signed1.json")
	s.Require().NoError(err)

	multisign, _, err := s.Chain.Validators[0].Exec(s.GetContext(), s.Chain.Validators[0].TxCommand(
		multisigName,
		"multisign",
		path.Join(s.Chain.Validators[0].HomeDir(), "tx.json"),
		multisigName,
		path.Join(s.Chain.Validators[0].HomeDir(), "signed0.json"),
		path.Join(s.Chain.Validators[0].HomeDir(), "signed1.json"),
	), nil)
	s.Require().NoError(err)

	err = s.Chain.Validators[0].WriteFile(s.GetContext(), multisign, "multisign.json")
	s.Require().NoError(err)

	_, err = s.Chain.Validators[0].ExecTx(s.GetContext(), multisigName, "broadcast", path.Join(s.Chain.Validators[0].HomeDir(), "multisign.json"))
	s.Require().NoError(err)
	balanceAfter, err := s.Chain.GetBalance(s.GetContext(), s.Chain.ValidatorWallets[3].Address, chainsuite.Uom)
	s.Require().NoError(err)
	s.Require().Equal(balanceBefore.Add(sdkmath.NewInt(txAmount)), balanceAfter)
}

func TestTransactions(t *testing.T) {
	txSuite := TxSuite{chainsuite.NewSuite(chainsuite.SuiteConfig{UpgradeOnSetup: true})}
	suite.Run(t, &txSuite)
}

func (s TxSuite) authzGenExec(ctx context.Context, grantee chainsuite.ValidatorWallet, command ...string) error {
	txjson, err := s.Chain.GenerateTx(ctx, 1, command...)
	s.Require().NoError(err)

	err = s.Chain.Validators[1].WriteFile(ctx, []byte(txjson), "tx.json")
	s.Require().NoError(err)

	_, err = s.Chain.Validators[1].ExecTx(
		ctx,
		grantee.Moniker,
		"authz", "exec", path.Join(s.Chain.Validators[1].HomeDir(), "tx.json"),
	)
	return err
}
