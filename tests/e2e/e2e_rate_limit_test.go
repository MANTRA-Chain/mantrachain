package e2e

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	sdkmath "cosmossdk.io/math"
)

const (
	proposalAddRateLimitAmantraFilename    = "proposal_add_rate_limit_amantra.json"
	proposalUpdateRateLimitAmantraFilename = "proposal_update_rate_limit_amantra.json"
	proposalResetRateLimitAmantraFilename  = "proposal_reset_rate_limit_amantra.json"
	proposalRemoveRateLimitAmantraFilename = "proposal_remove_rate_limit_amantra.json"
)

func (s *IntegrationTestSuite) writeAddRateLimitAmantraProposal(c *chain) {
	template := `
	{
		"messages": [
		 {
		  "@type": "/ratelimit.v1.MsgAddRateLimit",
		  "authority": "%s",
		  "denom": "%s",
		  "channel_or_client_id": "%s",
		  "max_percent_send": "%s",
		  "max_percent_recv": "%s",
		  "duration_hours": "%d"
		 }
		],
		"metadata": "ipfs://CID",
		"deposit": "100000000000000amantra",
		"title": "Add Rate Limit on (channel-0, amantra)",
		"summary": "e2e-test adding an IBC rate limit"
	   }`
	propMsgBody := fmt.Sprintf(template,
		govAuthority,
		amantraDenom,               // denom: amantra
		transferChannel,            // channel_id: channel-0
		sdkmath.NewInt(1).String(), // max_percent_send: 1%
		sdkmath.NewInt(1).String(), // max_percent_recv: 1%
		24,                         // duration_hours: 24
	)

	err := writeFile(filepath.Join(c.validators[0].configDir(), "config", proposalAddRateLimitAmantraFilename), []byte(propMsgBody))
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) writeUpdateRateLimitAmantraProposal(c *chain) {
	template := `
	{
		"messages": [
		 {
		  "@type": "/ratelimit.v1.MsgUpdateRateLimit",
		  "authority": "%s",
		  "denom": "%s",
		  "channel_or_client_id": "%s",
		  "max_percent_send": "%s",
		  "max_percent_recv": "%s",
		  "duration_hours": "%d"
		 }
		],
		"metadata": "ipfs://CID",
		"deposit": "100000000000000amantra",
		"title": "Update Rate Limit on (channel-0, amantra)",
		"summary": "e2e-test updating an IBC rate limit"
	   }`
	propMsgBody := fmt.Sprintf(template,
		govAuthority,
		amantraDenom,               // denom: amantra
		transferChannel,            // channel_id: channel-0
		sdkmath.NewInt(2).String(), // max_percent_send: 2%
		sdkmath.NewInt(1).String(), // max_percent_recv: 1%
		6,                          // duration_hours: 6
	)

	err := writeFile(filepath.Join(c.validators[0].configDir(), "config", proposalUpdateRateLimitAmantraFilename), []byte(propMsgBody))
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) writeResetRateLimitAmantraProposal(c *chain) {
	template := `
	{
		"messages": [
		 {
		  "@type": "/ratelimit.v1.MsgResetRateLimit",
		  "authority": "%s",
		  "denom": "%s",
		  "channel_or_client_id": "%s"
		 }
		],
		"metadata": "ipfs://CID",
		"deposit": "100000000000000amantra",
		"title": "Reset Rate Limit on (channel-0, amantra)",
		"summary": "e2e-test resetting an IBC rate limit"
	   }`
	propMsgBody := fmt.Sprintf(template,
		govAuthority,
		amantraDenom,    // denom: amantra
		transferChannel, // channel_id: channel-0
	)

	err := writeFile(filepath.Join(c.validators[0].configDir(), "config", proposalResetRateLimitAmantraFilename), []byte(propMsgBody))
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) writeRemoveRateLimitAmantraProposal(c *chain) {
	template := `
	{
		"messages": [
		 {
		  "@type": "/ratelimit.v1.MsgRemoveRateLimit",
		  "authority": "%s",
		  "denom": "%s",
		  "channel_or_client_id": "%s"
		 }
		],
		"metadata": "ipfs://CID",
		"deposit": "100000000000000amantra",
		"title": "Remove Rate Limit (channel-0, amantra)",
		"summary": "e2e-test removing an IBC rate limit"
	   }`
	propMsgBody := fmt.Sprintf(template,
		govAuthority,
		amantraDenom,    // denom: amantra
		transferChannel, // channel_id: channel-0
	)

	err := writeFile(filepath.Join(c.validators[0].configDir(), "config", proposalRemoveRateLimitAmantraFilename), []byte(propMsgBody))
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) testAddRateLimits() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorAAddr, _ := validatorA.keyInfo.GetAddress()

	s.writeAddRateLimitAmantraProposal(s.chainA)
	proposalCounter++
	submitGovFlags := []string{configFile(proposalAddRateLimitAmantraFilename)}
	depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.T().Logf("Submitting, deposit and vote Gov Proposal: Add IBC rate limit for (channel-0, amantra)")
	s.submitGovProposal(chainEndpoint, validatorAAddr.String(), proposalCounter, "ratelimittypes.MsgAddRateLimit", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

	s.Require().Eventually(
		func() bool {
			s.T().Logf("After AddRateLimit proposal (channel-0, amantra)")

			rateLimits, err := queryAllRateLimits(chainEndpoint)
			s.Require().NoError(err)
			s.Require().Len(rateLimits, 1)
			s.Require().Equal(transferChannel, rateLimits[0].Path.ChannelOrClientId)
			s.Require().Equal(amantraDenom, rateLimits[0].Path.Denom)
			s.Require().Equal(uint64(24), rateLimits[0].Quota.DurationHours)
			s.Require().Equal(sdkmath.NewInt(1), rateLimits[0].Quota.MaxPercentRecv)
			s.Require().Equal(sdkmath.NewInt(1), rateLimits[0].Quota.MaxPercentSend)

			res, err := queryRateLimit(chainEndpoint, transferChannel, amantraDenom)
			s.Require().NoError(err)
			s.Require().NotNil(res.RateLimit)
			s.Require().Equal(*rateLimits[0].Path, *res.RateLimit.Path)
			s.Require().Equal(*rateLimits[0].Quota, *res.RateLimit.Quota)

			rateLimitsByChainID, err := queryRateLimitsByChainID(chainEndpoint, s.chainB.id)
			s.Require().NoError(err)
			s.Require().Len(rateLimits, 1)
			s.Require().Equal(*rateLimits[0].Path, *rateLimitsByChainID[0].Path)
			s.Require().Equal(*rateLimits[0].Quota, *rateLimitsByChainID[0].Quota)

			return true
		},
		15*time.Second,
		5*time.Second,
	)
}

func (s *IntegrationTestSuite) testUpdateRateLimit() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorAAddr, _ := validatorA.keyInfo.GetAddress()

	s.writeUpdateRateLimitAmantraProposal(s.chainA)
	proposalCounter++
	submitGovFlags := []string{configFile(proposalUpdateRateLimitAmantraFilename)}
	depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.T().Logf("Submitting, deposit and vote Gov Proposal: Update IBC rate limit for (channel-0, amantra)")
	s.submitGovProposal(chainEndpoint, validatorAAddr.String(), proposalCounter, "ratelimittypes.MsgUpdateRateLimit", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

	s.Require().Eventually(
		func() bool {
			s.T().Logf("After UpdateRateLimit proposal")

			res, err := queryRateLimit(chainEndpoint, transferChannel, amantraDenom)
			s.Require().NoError(err)
			s.Require().NotNil(res.RateLimit)
			s.Require().Equal(sdkmath.NewInt(2), res.RateLimit.Quota.MaxPercentSend)
			s.Require().Equal(uint64(6), res.RateLimit.Quota.DurationHours)

			return true
		},
		15*time.Second,
		5*time.Second,
	)
}

func (s *IntegrationTestSuite) testResetRateLimit() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorAAddr, _ := validatorA.keyInfo.GetAddress()

	s.writeResetRateLimitAmantraProposal(s.chainA)
	proposalCounter++
	submitGovFlags := []string{configFile(proposalResetRateLimitAmantraFilename)}
	depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.T().Logf("Submitting, deposit and vote Gov Proposal: Reset IBC rate limit for (channel-0, amantra)")
	s.submitGovProposal(chainEndpoint, validatorAAddr.String(), proposalCounter, "ratelimittypes.MsgResetRateLimit", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

	s.Require().Eventually(
		func() bool {
			s.T().Logf("After ResetRateLimit proposal")

			res, err := queryRateLimit(chainEndpoint, transferChannel, amantraDenom)
			s.Require().NoError(err)
			s.Require().NotNil(res.RateLimit)
			s.Require().Equal(sdkmath.NewInt(0), res.RateLimit.Flow.Inflow)
			s.Require().Equal(sdkmath.NewInt(0), res.RateLimit.Flow.Outflow)

			return true
		},
		15*time.Second,
		5*time.Second,
	)
}

func (s *IntegrationTestSuite) testRemoveRateLimit() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorAAddr, _ := validatorA.keyInfo.GetAddress()

	s.writeRemoveRateLimitAmantraProposal(s.chainA)
	proposalCounter++
	submitGovFlags := []string{configFile(proposalRemoveRateLimitAmantraFilename)}
	depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.T().Logf("Submitting, deposit and vote Gov Proposal: Remove IBC rate limit for (channel-0, amantra)")
	s.submitGovProposal(chainEndpoint, validatorAAddr.String(), proposalCounter, "ratelimittypes.MsgRemoveRateLimit", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

	s.Require().Eventually(
		func() bool {
			s.T().Logf("After RemoveRateLimit proposal")

			rateLimits, err := queryAllRateLimits(chainEndpoint)
			s.Require().NoError(err)
			s.Require().Empty(rateLimits)

			res, err := queryRateLimit(chainEndpoint, transferChannel, amantraDenom)
			s.Require().NoError(err)
			s.Require().Nil(res.RateLimit)

			return true
		},
		15*time.Second,
		5*time.Second,
	)
}

func (s *IntegrationTestSuite) testIBCTransfer(expToFail bool) {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	address, _ := s.chainA.validators[0].keyInfo.GetAddress()
	sender := address.String()

	address, _ = s.chainB.validators[0].keyInfo.GetAddress()
	recipient := address.String()

	totalAmount, err := querySupplyOf(chainEndpoint, amantraDenom)
	s.Require().NoError(err)

	threshold := totalAmount.Amount.Mul(sdkmath.NewInt(1)).Quo(sdkmath.NewInt(100))
	tokenAmt := threshold.Add(sdkmath.NewInt(10)).String()
	s.sendIBC(s.chainA, 0, sender, recipient, tokenAmt+amantraDenom, standardFees.String(), "", expToFail)

	if !expToFail {
		s.T().Logf("After successful IBC transfer")

		res, err := queryRateLimit(chainEndpoint, transferChannel, amantraDenom)
		s.Require().NoError(err)
		s.Require().NotNil(res.RateLimit)
		s.Require().Equal(sdkmath.NewInt(0), res.RateLimit.Flow.Inflow)
		s.Require().NotEqual(sdkmath.NewInt(0), res.RateLimit.Flow.Outflow)
	}
}
