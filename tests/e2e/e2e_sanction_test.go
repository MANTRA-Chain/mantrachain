package e2e

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"
)

const (
	proposalAddBlacklistAccountsFilename    = "proposal_add_blacklist_accounts.json"
	proposalRemoveBlacklistAccountsFilename = "proposal_remove_blacklist_accounts.json"
)

func (s *IntegrationTestSuite) writeAddBlacklistAccountsProposal(c *chain, blockAccount string) {
	template := `
	{
		"messages": [
		 {
		  "@type": "/mantrachain.sanction.v1.MsgAddBlacklistAccounts",
		  "authority": "%s",
		  "blacklist_accounts": [
		  	"%s"
		  ]
		 }
		],
		"metadata": "ipfs://CID",
		"deposit": "100uom",
		"title": "Add %s to blacklist",
		"summary": "e2e-test adding to blacklist"
	   }`
	propMsgBody := fmt.Sprintf(template,
		govAuthority,
		blockAccount,
		blockAccount,
	)

	err := writeFile(filepath.Join(c.validators[0].configDir(), "config", proposalAddBlacklistAccountsFilename), []byte(propMsgBody))
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) testAddToBlacklist() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorAddr, _ := validatorA.keyInfo.GetAddress()

	valIdx := 0
	alice, _ := s.chainA.genesisAccounts[1].keyInfo.GetAddress()
	bob, _ := s.chainA.genesisAccounts[2].keyInfo.GetAddress()
	// able to send tokens from alice to bob before blacklist
	s.execBankSend(s.chainA, valIdx, alice.String(), bob.String(), tokenAmount.String(), standardFees.String(), false)

	s.writeAddBlacklistAccountsProposal(s.chainA, alice.String())
	proposalCounter++
	submitGovFlags := []string{configFile(proposalAddBlacklistAccountsFilename)}
	depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.T().Logf("Submitting, deposit and vote Gov Proposal: Add %s to blacklist", alice.String())
	s.submitGovProposal(chainEndpoint, validatorAddr.String(), proposalCounter, "sanctiontypes.MsgAddBlacklistAccounts", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

	s.Require().Eventually(
		func() bool {
			s.T().Logf("After AddBlacklistAccount proposal")

			blacklist, err := queryBlacklist(chainEndpoint)
			s.Require().NoError(err)
			s.Require().Len(blacklist, 1)
			s.Require().Equal(alice.String(), blacklist[0])

			return true
		},
		15*time.Second,
		5*time.Second,
	)

	// alice sends tokens to bob
	s.execBankSend(s.chainA, valIdx, alice.String(), bob.String(), tokenAmount.String(), standardFees.String(), true)
	s.T().Logf("Failed to send token from Alice to Bob as Alice is blacklisted")
}

func (s *IntegrationTestSuite) writeRemoveBlacklistAccountsProposal(c *chain, blockAccount string) {
	template := `
	{
		"messages": [
		 {
		  "@type": "/mantrachain.sanction.v1.MsgRemoveBlacklistAccounts",
		  "authority": "%s",
		  "blacklist_accounts": [
		  	"%s"
		  ]
		 }
		],
		"metadata": "ipfs://CID",
		"deposit": "100uom",
		"title": "Remove %s to blacklist",
		"summary": "e2e-test remove from blacklist"
	   }`
	propMsgBody := fmt.Sprintf(template,
		govAuthority,
		blockAccount,
		blockAccount,
	)

	err := writeFile(filepath.Join(c.validators[0].configDir(), "config", proposalRemoveBlacklistAccountsFilename), []byte(propMsgBody))
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) testRemoveFromBlacklist() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorAAddr, _ := validatorA.keyInfo.GetAddress()

	valIdx := 0
	alice, _ := s.chainA.genesisAccounts[1].keyInfo.GetAddress()
	bob, _ := s.chainA.genesisAccounts[2].keyInfo.GetAddress()
	// unable to send tokens from alice to bob before remove from blacklist
	s.execBankSend(s.chainA, valIdx, alice.String(), bob.String(), tokenAmount.String(), standardFees.String(), true)
	s.T().Logf("Failed to send token from Alice to Bob as Alice is blacklisted")

	s.writeRemoveBlacklistAccountsProposal(s.chainA, alice.String())
	proposalCounter++
	submitGovFlags := []string{configFile(proposalRemoveBlacklistAccountsFilename)}
	depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.T().Logf("Submitting, deposit and vote Gov Proposal: Remove %s from blacklist", alice.String())
	s.submitGovProposal(chainEndpoint, validatorAAddr.String(), proposalCounter, "sanctiontypes.MsgRemoveBlacklistAccounts", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

	s.Require().Eventually(
		func() bool {
			s.T().Logf("After RemoveBlacklistAccount proposal")

			blacklist, err := queryBlacklist(chainEndpoint)
			s.Require().NoError(err)
			s.Require().Len(blacklist, 0)

			return true
		},
		15*time.Second,
		5*time.Second,
	)

	// alice sends tokens to bob
	s.execBankSend(s.chainA, valIdx, alice.String(), bob.String(), tokenAmount.String(), standardFees.String(), false)
}
