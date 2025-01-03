package e2e

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testDistribution() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorB := s.chainA.validators[1]
	validatorBAddr, _ := validatorB.keyInfo.GetAddress()

	valOperAddressA := sdk.ValAddress(validatorBAddr).String()

	delegatorAddress, _ := s.chainA.genesisAccounts[2].keyInfo.GetAddress()

	newWithdrawalAddress, _ := s.chainA.genesisAccounts[3].keyInfo.GetAddress()

	beforeBalance, err := getSpecificBalance(chainEndpoint, newWithdrawalAddress.String(), uomDenom)
	s.Require().NoError(err)
	if beforeBalance.IsNil() {
		beforeBalance = sdk.NewCoin(uomDenom, math.NewInt(0))
	}

	s.execSetWithdrawAddress(s.chainA, 0, standardFees.String(), delegatorAddress.String(), newWithdrawalAddress.String(), mantraHomePath)

	// Verify
	s.Require().Eventually(
		func() bool {
			res, err := queryDelegatorWithdrawalAddress(chainEndpoint, delegatorAddress.String())
			s.Require().NoError(err)

			return res.WithdrawAddress == newWithdrawalAddress.String()
		},
		10*time.Second,
		5*time.Second,
	)

	s.execWithdrawReward(s.chainA, 0, delegatorAddress.String(), valOperAddressA, mantraHomePath)
	s.Require().Eventually(
		func() bool {
			afterBalance, err := getSpecificBalance(chainEndpoint, newWithdrawalAddress.String(), uomDenom)
			s.Require().NoError(err)

			return afterBalance.IsGTE(beforeBalance)
		},
		10*time.Second,
		5*time.Second,
	)
}

/*
fundCommunityPool tests the funding of the community pool on behalf of the distribution module.
Test Benchmarks:
1. Validation that balance of the distribution module account before funding
2. Execution funding the community pool
3. Verification that correct funds have been deposited to distribution module account
*/
func (s *IntegrationTestSuite) fundCommunityPool() {
	chainAAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
	sender, _ := s.chainA.validators[0].keyInfo.GetAddress()

	beforeDistUomBalance, _ := getSpecificBalance(chainAAPIEndpoint, distModuleAddress, tokenAmount.Denom)
	if beforeDistUomBalance.IsNil() {
		// Set balance to 0 if previous balance does not exist
		beforeDistUomBalance = sdk.NewInt64Coin(uomDenom, 0)
	}

	s.execDistributionFundCommunityPool(s.chainA, 0, sender.String(), tokenAmount.String(), standardFees.String())

	s.Require().Eventually(
		func() bool {
			afterDistUomBalance, err := getSpecificBalance(chainAAPIEndpoint, distModuleAddress, tokenAmount.Denom)
			s.Require().NoErrorf(err, "Error getting balance: %s", afterDistUomBalance)

			// check if the balance is increased by the tokenAmount
			return beforeDistUomBalance.Add(tokenAmount).IsEqual(afterDistUomBalance)
		},
		15*time.Second,
		5*time.Second,
	)
}
