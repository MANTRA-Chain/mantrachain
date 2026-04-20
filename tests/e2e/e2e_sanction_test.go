package e2e

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ory/dockertest/v3/docker"
)

const (
	proposalAddBlacklistAccountsFilename    = "proposal_add_blacklist_accounts.json"
	proposalRemoveBlacklistAccountsFilename = "proposal_remove_blacklist_accounts.json"
	authzExecMsgFilename                    = "authz_exec_msg.json"
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
		"deposit": "100000000000000amantra",
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
		"deposit": "100000000000000amantra",
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
			s.Require().Empty(blacklist)

			return true
		},
		15*time.Second,
		5*time.Second,
	)

	// alice sends tokens to bob
	s.execBankSend(s.chainA, valIdx, alice.String(), bob.String(), tokenAmount.String(), standardFees.String(), false)
}

// generateAuthzExecTxFile runs bank send with --generate-only inside the validator container
// and writes the resulting unsigned tx JSON to the authz exec msg file so it can be used
// by `tx authz exec`, which calls ReadTxFromFile and expects a full tx (not a message array).
func (s *IntegrationTestSuite) generateAuthzExecTxFile(c *chain, valIdx int, from, to, amount string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	mantraCommand := []string{
		mantrachaindBinary,
		txCommand,
		"bank",
		"send",
		from,
		to,
		amount + amantraDenom,
		"--generate-only",
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", flags.FlagHome, mantraHomePath),
		"--keyring-backend=test",
		"--output=json",
	}

	var outBuf, errBuf bytes.Buffer
	exec, err := s.dkrPool.Client.CreateExec(docker.CreateExecOptions{
		Context:      ctx,
		AttachStdout: true,
		AttachStderr: true,
		Container:    s.valResources[c.id][valIdx].Container.ID,
		User:         "nonroot",
		Cmd:          mantraCommand,
	})
	s.Require().NoError(err)

	err = s.dkrPool.Client.StartExec(exec.ID, docker.StartExecOptions{
		Context:      ctx,
		Detach:       false,
		OutputStream: &outBuf,
		ErrorStream:  &errBuf,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(outBuf.Bytes(), "generate-only produced no output; stderr: %s", errBuf.String())

	err = writeFile(filepath.Join(c.validators[0].configDir(), "config", authzExecMsgFilename), outBuf.Bytes())
	s.Require().NoError(err)
}

// execAuthzGrantSend grants a send authorization from granter to grantee.
func (s *IntegrationTestSuite) execAuthzGrantSend(c *chain, valIdx int, granter, grantee, spendLimit string, opt ...flagOption) {
	opt = append(opt, withKeyValue(flagFrom, granter))
	opt = append(opt, withKeyValue(flagSpendLimit, spendLimit))
	opts := applyOptions(c.id, opt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("granting authz send from %s to %s on chain %s", granter, grantee, c.id)

	mantraCommand := []string{
		mantrachaindBinary,
		txCommand,
		"authz",
		"grant",
		grantee,
		"send",
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		fmt.Sprintf("--%s=%s", flags.FlagGas, "300000"),
		"--keyring-backend=test",
		"--output=json",
		"-y",
	}
	for flag, value := range opts {
		mantraCommand = append(mantraCommand, fmt.Sprintf("--%s=%s", flag, value))
	}

	s.executeTxCommand(ctx, c, mantraCommand, valIdx, s.defaultExecValidation(c, valIdx))
}

// execAuthzExec executes authz exec using the pre-written msg file.
func (s *IntegrationTestSuite) execAuthzExec(c *chain, valIdx int, grantee string, expectErr bool, opt ...flagOption) {
	opt = append(opt, withKeyValue(flagFrom, grantee))
	opts := applyOptions(c.id, opt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("executing authz exec as %s on chain %s", grantee, c.id)

	mantraCommand := []string{
		mantrachaindBinary,
		txCommand,
		"authz",
		"exec",
		configFile(authzExecMsgFilename),
		"-y",
	}
	for flag, value := range opts {
		mantraCommand = append(mantraCommand, fmt.Sprintf("--%s=%v", flag, value))
	}

	s.executeTxCommand(ctx, c, mantraCommand, valIdx, s.expectErrExecValidation(c, valIdx, expectErr))
}

// testAuthzGranterBlacklist verifies that a MsgExec transaction is blocked when the authz
// granter (inner message signer) is blacklisted, and succeeds again after removal.
func (s *IntegrationTestSuite) testAuthzGranterBlacklist() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorAddr, _ := validatorA.keyInfo.GetAddress()

	valIdx := 0
	alice, _ := s.chainA.genesisAccounts[1].keyInfo.GetAddress()   // granter (will be blacklisted)
	bob, _ := s.chainA.genesisAccounts[2].keyInfo.GetAddress()     // grantee (executes on alice's behalf)
	charlie, _ := s.chainA.genesisAccounts[3].keyInfo.GetAddress() // recipient

	// Grant a large spend limit so repeated execs never exhaust the allowance
	s.execAuthzGrantSend(s.chainA, valIdx, alice.String(), bob.String(), tokenAmount.String())

	// Write the authz exec message (bob sends a small amount from alice to charlie)
	s.generateAuthzExecTxFile(s.chainA, valIdx, alice.String(), charlie.String(), "1")

	// Bob executes the authz send — should succeed (alice is not yet blacklisted)
	s.T().Logf("Bob executes authz send on behalf of Alice — expect success")
	s.execAuthzExec(s.chainA, valIdx, bob.String(), false)

	// Blacklist alice via governance
	s.writeAddBlacklistAccountsProposal(s.chainA, alice.String())
	proposalCounter++
	submitGovFlags := []string{configFile(proposalAddBlacklistAccountsFilename)}
	depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.T().Logf("Submitting Gov Proposal: Add %s (authz granter) to blacklist", alice.String())
	s.submitGovProposal(chainEndpoint, validatorAddr.String(), proposalCounter, "sanctiontypes.MsgAddBlacklistAccounts", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

	s.Require().Eventually(
		func() bool {
			blacklist, err := queryBlacklist(chainEndpoint)
			s.Require().NoError(err)
			s.Require().Len(blacklist, 1)
			s.Require().Equal(alice.String(), blacklist[0])
			return true
		},
		15*time.Second,
		5*time.Second,
	)

	// Bob tries to execute the authz send — should fail (alice is blacklisted granter)
	s.T().Logf("Bob executes authz send on behalf of Alice — expect failure (alice blacklisted)")
	s.execAuthzExec(s.chainA, valIdx, bob.String(), true)

	// Remove alice from blacklist via governance
	s.writeRemoveBlacklistAccountsProposal(s.chainA, alice.String())
	proposalCounter++
	submitGovFlags = []string{configFile(proposalRemoveBlacklistAccountsFilename)}
	depositGovFlags = []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags = []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.T().Logf("Submitting Gov Proposal: Remove %s (authz granter) from blacklist", alice.String())
	s.submitGovProposal(chainEndpoint, validatorAddr.String(), proposalCounter, "sanctiontypes.MsgRemoveBlacklistAccounts", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

	s.Require().Eventually(
		func() bool {
			blacklist, err := queryBlacklist(chainEndpoint)
			s.Require().NoError(err)
			s.Require().Empty(blacklist)
			return true
		},
		15*time.Second,
		5*time.Second,
	)

	// Bob executes the authz send — should succeed again
	s.T().Logf("Bob executes authz send on behalf of Alice — expect success (alice removed from blacklist)")
	s.execAuthzExec(s.chainA, valIdx, bob.String(), false)
}

// testFeeGranterBlacklist verifies that a transaction is blocked when the fee granter
// is blacklisted, and succeeds again after removal.
func (s *IntegrationTestSuite) testFeeGranterBlacklist() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorAddr, _ := validatorA.keyInfo.GetAddress()

	valIdx := 0
	alice, _ := s.chainA.genesisAccounts[1].keyInfo.GetAddress()   // fee granter (will be blacklisted)
	bob, _ := s.chainA.genesisAccounts[2].keyInfo.GetAddress()     // sender (pays with alice's fee grant)
	charlie, _ := s.chainA.genesisAccounts[3].keyInfo.GetAddress() // recipient

	// Alice grants bob a fee allowance
	s.execFeeGrant(s.chainA, valIdx, alice.String(), bob.String(), tokenAmount.String())

	// Bob sends to charlie using alice's fee grant — should succeed
	s.T().Logf("Bob sends with Alice's fee grant — expect success")
	s.execBankSend(s.chainA, valIdx, bob.String(), charlie.String(), tokenAmount.String(), standardFees.String(), false,
		withKeyValue(flagFeeGranter, alice.String()),
	)

	// Blacklist alice via governance
	s.writeAddBlacklistAccountsProposal(s.chainA, alice.String())
	proposalCounter++
	submitGovFlags := []string{configFile(proposalAddBlacklistAccountsFilename)}
	depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.T().Logf("Submitting Gov Proposal: Add %s (fee granter) to blacklist", alice.String())
	s.submitGovProposal(chainEndpoint, validatorAddr.String(), proposalCounter, "sanctiontypes.MsgAddBlacklistAccounts", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

	s.Require().Eventually(
		func() bool {
			blacklist, err := queryBlacklist(chainEndpoint)
			s.Require().NoError(err)
			s.Require().Len(blacklist, 1)
			s.Require().Equal(alice.String(), blacklist[0])
			return true
		},
		15*time.Second,
		5*time.Second,
	)

	// Bob tries to send with alice's fee grant — should fail (alice is blacklisted fee granter)
	s.T().Logf("Bob sends with Alice's fee grant — expect failure (alice blacklisted)")
	s.execBankSend(s.chainA, valIdx, bob.String(), charlie.String(), tokenAmount.String(), standardFees.String(), true,
		withKeyValue(flagFeeGranter, alice.String()),
	)

	// Remove alice from blacklist via governance
	s.writeRemoveBlacklistAccountsProposal(s.chainA, alice.String())
	proposalCounter++
	submitGovFlags = []string{configFile(proposalRemoveBlacklistAccountsFilename)}
	depositGovFlags = []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags = []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.T().Logf("Submitting Gov Proposal: Remove %s (fee granter) from blacklist", alice.String())
	s.submitGovProposal(chainEndpoint, validatorAddr.String(), proposalCounter, "sanctiontypes.MsgRemoveBlacklistAccounts", submitGovFlags, depositGovFlags, voteGovFlags, "vote")

	s.Require().Eventually(
		func() bool {
			blacklist, err := queryBlacklist(chainEndpoint)
			s.Require().NoError(err)
			s.Require().Empty(blacklist)
			return true
		},
		15*time.Second,
		5*time.Second,
	)

	// Bob sends with alice's fee grant — should succeed again
	s.T().Logf("Bob sends with Alice's fee grant — expect success (alice removed from blacklist)")
	s.execBankSend(s.chainA, valIdx, bob.String(), charlie.String(), tokenAmount.String(), standardFees.String(), false,
		withKeyValue(flagFeeGranter, alice.String()),
	)
}
