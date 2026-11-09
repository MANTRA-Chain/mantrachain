package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	sanctiontypes "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	proposalUpdateSanctionersFilename = "proposal_update_sanctioners.json"
	sanctionerTxGas                   = "300000"
)

// writeUpdateSanctionersProposal writes a proposal that replaces the
// sanctioner list. Empty full and addOnly write an empty list. Proposals here
// carry one message each, because the harness checks proposal status through
// gov v1beta1, which can't represent multi-message proposals.
func (s *IntegrationTestSuite) writeUpdateSanctionersProposal(c *chain, full, addOnly string) {
	sanctioners := "[]"
	if full != "" || addOnly != "" {
		sanctioners = fmt.Sprintf(`[
			{"address": "%s", "role": "SANCTIONER_ROLE_FULL"},
			{"address": "%s", "role": "SANCTIONER_ROLE_ADD_ONLY"}
		]`, full, addOnly)
	}

	propMsgBody := fmt.Sprintf(`{
		"messages": [{
			"@type": "/mantrachain.sanction.v1.MsgUpdateSanctioners",
			"authority": "%s",
			"sanctioners": %s
		}],
		"metadata": "ipfs://CID",
		"deposit": "100000000000000amantra",
		"title": "Update sanctioners",
		"summary": "e2e-test updating the sanctioner list"
	}`, govAuthority, sanctioners)

	err := writeFile(filepath.Join(c.validators[0].configDir(), "config", proposalUpdateSanctionersFilename), []byte(propMsgBody))
	s.Require().NoError(err)
}

// submitSanctionProposal submits, deposits on and votes yes for the proposal
// in filename.
func (s *IntegrationTestSuite) submitSanctionProposal(chainEndpoint, sender, filename, proposalType string) {
	proposalCounter++
	submitGovFlags := []string{configFile(filename)}
	depositGovFlags := []string{strconv.Itoa(proposalCounter), depositAmount.String()}
	voteGovFlags := []string{strconv.Itoa(proposalCounter), "yes"}

	s.T().Logf("Proposal number: %d", proposalCounter)
	s.submitGovProposal(chainEndpoint, sender, proposalCounter, proposalType, submitGovFlags, depositGovFlags, voteGovFlags, "vote")
}

// execSanctionTx runs `tx sanction <subcommand> <arg> --reason <reason>` as
// from, and checks the tx succeeds or fails as expected.
func (s *IntegrationTestSuite) execSanctionTx(c *chain, valIdx int, subcommand, arg, reason, from string, expectErr bool) { //nolint:unparam
	opts := applyOptions(c.id, []flagOption{
		withKeyValue(flagFrom, from),
		withKeyValue(flagFees, standardFees.String()),
		withKeyValue(flagGas, sanctionerTxGas),
		withKeyValue("reason", reason),
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	s.T().Logf("%s %s as %s on chain %s, expect error: %v", subcommand, arg, from, c.id, expectErr)

	mantraCommand := []string{
		mantrachaindBinary,
		txCommand,
		sanctiontypes.ModuleName,
		subcommand,
		arg,
		"-y",
	}
	for flag, value := range opts {
		mantraCommand = append(mantraCommand, fmt.Sprintf("--%s=%v", flag, value))
	}

	s.executeTxCommand(ctx, c, mantraCommand, valIdx, s.committedTxValidation(c, valIdx, expectErr))
}

// committedTxValidation waits for the tx to be committed before checking its
// result. Sanctioner errors happen during execution, after CheckTx has used
// the account sequence, so returning early could make the signer's next tx
// fail with a sequence mismatch. A tx rejected in CheckTx never commits, so
// that case returns straight away.
func (s *IntegrationTestSuite) committedTxValidation(c *chain, valIdx int, expectErr bool) func([]byte, []byte) bool {
	return func(stdOut []byte, stdErr []byte) bool {
		var txResp sdk.TxResponse
		if err := cdc.UnmarshalJSON(stdOut, &txResp); err != nil {
			s.Require().True(expectErr, "stdOut: %s, stdErr: %s", string(stdOut), string(stdErr))
			return true
		}
		if txResp.Code != 0 {
			s.Require().True(expectErr, "rejected in CheckTx: %s", txResp.RawLog)
			return true
		}

		endpoint := fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		var code uint32
		s.Require().Eventuallyf(
			func() bool {
				var found bool
				code, found = queryTxCode(endpoint, txResp.TxHash)
				return found
			},
			time.Minute,
			time.Second,
			"tx %s was not committed", txResp.TxHash,
		)
		s.Require().Equal(expectErr, code != 0, "tx %s committed with code %d", txResp.TxHash, code)
		return true
	}
}

// testSanctionerFastPath walks the sanctioner flow end to end on a running
// chain: governance appoints a full and an add-only sanctioner and sanctions
// one address itself, the add-only
// sanctioner blocks an account, the full sanctioner reverses it and revokes
// the add-only sanctioner, and governance finally empties the list.
// It leaves the blacklist and sanctioner list empty.
func (s *IntegrationTestSuite) testSanctionerFastPath() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorAddr, _ := validatorA.keyInfo.GetAddress()

	valIdx := 0
	alice, _ := s.chainA.genesisAccounts[1].keyInfo.GetAddress()   // recipient
	bob, _ := s.chainA.genesisAccounts[2].keyInfo.GetAddress()     // target
	manager, _ := s.chainA.genesisAccounts[3].keyInfo.GetAddress() // full sanctioner
	bot, _ := s.chainA.genesisAccounts[4].keyInfo.GetAddress()     // add-only sanctioner
	// An address with no key, sanctioned by governance, to check sanctioners
	// can't remove governance entries.
	govEntry := sdk.AccAddress([]byte("e2e-gov-sanctioned__")).String()

	// Governance appoints the sanctioners.
	s.writeUpdateSanctionersProposal(s.chainA, manager.String(), bot.String())
	s.T().Logf("Submitting Gov Proposal: appoint %s (full) and %s (add-only)", manager, bot)
	s.submitSanctionProposal(chainEndpoint, validatorAddr.String(), proposalUpdateSanctionersFilename, "sanctiontypes.MsgUpdateSanctioners")

	s.Require().Eventually(
		func() bool {
			sanctioners, err := querySanctioners(chainEndpoint)
			s.Require().NoError(err)
			roles := map[string]sanctiontypes.SanctionerRole{}
			for _, sn := range sanctioners {
				roles[sn.Address] = sn.Role
			}
			return len(roles) == 2 &&
				roles[manager.String()] == sanctiontypes.SanctionerRole_SANCTIONER_ROLE_FULL &&
				roles[bot.String()] == sanctiontypes.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY
		},
		15*time.Second,
		5*time.Second,
	)

	// Governance sanctions govEntry.
	s.writeAddBlacklistAccountsProposal(s.chainA, govEntry)
	s.T().Logf("Submitting Gov Proposal: Add %s to blacklist", govEntry)
	s.submitSanctionProposal(chainEndpoint, validatorAddr.String(), proposalAddBlacklistAccountsFilename, "sanctiontypes.MsgSanctionAddresses")
	s.Require().Eventually(
		func() bool {
			res, err := queryIsSanctioned(chainEndpoint, govEntry)
			s.Require().NoError(err)
			return res.Sanctioned && res.Record != nil && res.Record.Actor == govAuthority
		},
		15*time.Second,
		5*time.Second,
	)

	// Bob can send before he's sanctioned.
	s.execBankSend(s.chainA, valIdx, bob.String(), alice.String(), tokenAmount.String(), standardFees.String(), false)

	// The add-only sanctioner blocks bob in one transaction.
	reason := "rule drain-v1, tx 0xe2e"
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", bob.String(), reason, bot.String(), false)
	s.Require().Eventually(
		func() bool {
			res, err := queryIsSanctioned(chainEndpoint, bob.String())
			s.Require().NoError(err)
			return res.Sanctioned &&
				res.Record != nil &&
				res.Record.Actor == bot.String() &&
				res.Record.ActorRole == sanctiontypes.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY &&
				res.Record.Reason == reason
		},
		15*time.Second,
		5*time.Second,
	)
	s.T().Logf("Bob sends while sanctioned, expect failure")
	s.execBankSend(s.chainA, valIdx, bob.String(), alice.String(), tokenAmount.String(), standardFees.String(), true)

	// The add-only sanctioner can't unblock anyone.
	s.execSanctionTx(s.chainA, valIdx, "unsanction-addresses", bob.String(), "should fail", bot.String(), true)

	// Sanctioners can't sanction another sanctioner or the gov authority.
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", manager.String(), "should fail", bot.String(), true)
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", govAuthority, "should fail", manager.String(), true)

	// Sanctioners can't remove a governance entry.
	s.execSanctionTx(s.chainA, valIdx, "unsanction-addresses", govEntry, "should fail", manager.String(), true)

	// The full sanctioner reverses the false positive.
	s.execSanctionTx(s.chainA, valIdx, "unsanction-addresses", bob.String(), "false positive", manager.String(), false)
	s.Require().Eventually(
		func() bool {
			res, err := queryIsSanctioned(chainEndpoint, bob.String())
			s.Require().NoError(err)
			return !res.Sanctioned
		},
		15*time.Second,
		5*time.Second,
	)
	s.execBankSend(s.chainA, valIdx, bob.String(), alice.String(), tokenAmount.String(), standardFees.String(), false)

	// The full sanctioner revokes the add-only sanctioner, which then can't
	// sanction anyone.
	s.execSanctionTx(s.chainA, valIdx, "revoke-sanctioner", bot.String(), "key rotated", manager.String(), false)
	s.Require().Eventually(
		func() bool {
			sanctioners, err := querySanctioners(chainEndpoint)
			s.Require().NoError(err)
			return len(sanctioners) == 1 && sanctioners[0].Address == manager.String()
		},
		15*time.Second,
		5*time.Second,
	)
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", bob.String(), "should fail", bot.String(), true)

	// Governance removes its own entry, then empties the list, turning the
	// fast path off.
	s.writeRemoveBlacklistAccountsProposal(s.chainA, govEntry)
	s.T().Logf("Submitting Gov Proposal: Remove %s from blacklist", govEntry)
	s.submitSanctionProposal(chainEndpoint, validatorAddr.String(), proposalRemoveBlacklistAccountsFilename, "sanctiontypes.MsgUnsanctionAddresses")

	s.writeUpdateSanctionersProposal(s.chainA, "", "")
	s.T().Logf("Submitting Gov Proposal: empty the sanctioner list")
	s.submitSanctionProposal(chainEndpoint, validatorAddr.String(), proposalUpdateSanctionersFilename, "sanctiontypes.MsgUpdateSanctioners")

	s.Require().Eventually(
		func() bool {
			sanctioners, err := querySanctioners(chainEndpoint)
			s.Require().NoError(err)
			blacklist, err := queryBlacklist(chainEndpoint)
			s.Require().NoError(err)
			return len(sanctioners) == 0 && len(blacklist) == 0
		},
		15*time.Second,
		5*time.Second,
	)
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", bob.String(), "should fail", manager.String(), true)
}

// queryTxCode returns a committed tx's result code, and false if the tx
// isn't found yet.
func queryTxCode(endpoint, txHash string) (uint32, bool) {
	body, err := httpGet(fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", endpoint, txHash))
	if err != nil {
		return 0, false
	}
	var resp struct {
		TxResponse *struct {
			Code uint32 `json:"code"`
		} `json:"tx_response"`
	}
	if err := json.Unmarshal(body, &resp); err != nil || resp.TxResponse == nil {
		return 0, false
	}
	return resp.TxResponse.Code, true
}

func querySanctioners(endpoint string) ([]sanctiontypes.Sanctioner, error) {
	body, err := httpGet(fmt.Sprintf("%s/MANTRA-Chain/mantrachain/sanction/v1/sanctioners", endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var resp sanctiontypes.QuerySanctionersResponse
	if err := cdc.UnmarshalJSON(body, &resp); err != nil {
		return nil, err
	}

	return resp.Sanctioners, nil
}

func queryIsSanctioned(endpoint, address string) (*sanctiontypes.QueryIsSanctionedResponse, error) {
	body, err := httpGet(fmt.Sprintf("%s/MANTRA-Chain/mantrachain/sanction/v1/is_sanctioned/%s", endpoint, address))
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	var resp sanctiontypes.QueryIsSanctionedResponse
	if err := cdc.UnmarshalJSON(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
