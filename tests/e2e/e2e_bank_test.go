package e2e

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testBankTokenTransfer() {
	s.Run("send_tokens_between_accounts", func() {
		var (
			err           error
			valIdx        = 0
			c             = s.chainA
			chainEndpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)

		// define one sender and two recipient accounts
		alice, _ := c.genesisAccounts[1].keyInfo.GetAddress()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()
		charlie, _ := c.genesisAccounts[3].keyInfo.GetAddress()

		var beforeAliceAmantraBalance,
			beforeBobAmantraBalance,
			beforeCharlieAmantraBalance,
			afterAliceAmantraBalance,
			afterBobAmantraBalance,
			afterCharlieAmantraBalance sdk.Coin

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				beforeAliceAmantraBalance, err = getSpecificBalance(chainEndpoint, alice.String(), amantraDenom)
				s.Require().NoError(err)

				beforeBobAmantraBalance, err = getSpecificBalance(chainEndpoint, bob.String(), amantraDenom)
				s.Require().NoError(err)

				beforeCharlieAmantraBalance, err = getSpecificBalance(chainEndpoint, charlie.String(), amantraDenom)
				s.Require().NoError(err)

				return beforeAliceAmantraBalance.IsValid() && beforeBobAmantraBalance.IsValid() && beforeCharlieAmantraBalance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		// alice sends tokens to bob
		s.execBankSend(s.chainA, valIdx, alice.String(), bob.String(), tokenAmount.String(), standardFees.String(), false)

		// check that the transfer was successful
		s.Require().Eventually(
			func() bool {
				afterAliceAmantraBalance, err = getSpecificBalance(chainEndpoint, alice.String(), amantraDenom)
				s.Require().NoError(err)

				afterBobAmantraBalance, err = getSpecificBalance(chainEndpoint, bob.String(), amantraDenom)
				s.Require().NoError(err)

				decremented := beforeAliceAmantraBalance.Sub(tokenAmount).Sub(standardFees).IsEqual(afterAliceAmantraBalance)
				incremented := beforeBobAmantraBalance.Add(tokenAmount).IsEqual(afterBobAmantraBalance)

				return decremented && incremented
			},
			10*time.Second,
			5*time.Second,
		)

		// save the updated account balances of alice and bob
		beforeAliceAmantraBalance, beforeBobAmantraBalance = afterAliceAmantraBalance, afterBobAmantraBalance

		// alice sends tokens to bob and charlie, at once
		s.execBankMultiSend(s.chainA, valIdx, alice.String(), []string{bob.String(), charlie.String()}, tokenAmount.String(), standardFees.String(), false)

		s.Require().Eventually(
			func() bool {
				afterAliceAmantraBalance, err = getSpecificBalance(chainEndpoint, alice.String(), amantraDenom)
				s.Require().NoError(err)

				afterBobAmantraBalance, err = getSpecificBalance(chainEndpoint, bob.String(), amantraDenom)
				s.Require().NoError(err)

				afterCharlieAmantraBalance, err = getSpecificBalance(chainEndpoint, charlie.String(), amantraDenom)
				s.Require().NoError(err)

				decremented := beforeAliceAmantraBalance.Sub(tokenAmount).Sub(tokenAmount).Sub(standardFees).IsEqual(afterAliceAmantraBalance)
				incremented := beforeBobAmantraBalance.Add(tokenAmount).IsEqual(afterBobAmantraBalance) &&
					beforeCharlieAmantraBalance.Add(tokenAmount).IsEqual(afterCharlieAmantraBalance)

				return decremented && incremented
			},
			10*time.Second,
			5*time.Second,
		)
	})
}
