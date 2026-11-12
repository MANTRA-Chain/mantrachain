package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	sanctiontypes "github.com/MANTRA-Chain/mantrachain/v8/x/sanction/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
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
// from and returns the tx hash. With wantErr nil the tx must succeed;
// otherwise it must fail with exactly that error's codespace and code,
// whether in CheckTx or on execution.
func (s *IntegrationTestSuite) execSanctionTx(c *chain, valIdx int, subcommand, arg, reason, from string, wantErr *errorsmod.Error) string { //nolint:unparam
	opts := applyOptions(c.id, []flagOption{
		withKeyValue(flagFrom, from),
		withKeyValue(flagFees, standardFees.String()),
		withKeyValue(flagGas, sanctionerTxGas),
		withKeyValue("reason", reason),
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	want := "success"
	if wantErr != nil {
		want = wantErr.Error()
	}
	s.T().Logf("%s %s as %s on chain %s, want: %s", subcommand, arg, from, c.id, want)

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

	return s.executeTxCommand(ctx, c, mantraCommand, valIdx, s.committedTxValidation(c, valIdx, wantErr))
}

// committedTxValidation waits for the tx to be committed before checking its
// result. Sanctioner errors happen during execution, after CheckTx has used
// the account sequence, so returning early could make the signer's next tx
// fail with a sequence mismatch. A tx rejected in CheckTx never commits, so
// that case is checked straight away.
func (s *IntegrationTestSuite) committedTxValidation(c *chain, valIdx int, wantErr *errorsmod.Error) func([]byte, []byte) bool {
	return func(stdOut []byte, stdErr []byte) bool {
		var txResp sdk.TxResponse
		s.Require().NoError(cdc.UnmarshalJSON(stdOut, &txResp), "stdOut: %s, stdErr: %s", string(stdOut), string(stdErr))
		if txResp.Code != 0 {
			s.requireTxOutcome("CheckTx", txResp.Codespace, txResp.Code, txResp.RawLog, wantErr)
			return true
		}

		endpoint := fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		var res *txResult
		s.Require().Eventuallyf(
			func() bool {
				var found bool
				res, found = queryTxResult(endpoint, txResp.TxHash)
				return found
			},
			time.Minute,
			time.Second,
			"tx %s was not committed", txResp.TxHash,
		)
		s.requireTxOutcome("execution", res.Codespace, res.Code, res.RawLog, wantErr)
		return true
	}
}

func (s *IntegrationTestSuite) requireTxOutcome(stage, codespace string, code uint32, log string, wantErr *errorsmod.Error) {
	if wantErr == nil {
		s.Require().Zero(code, "tx failed in %s: %s", stage, log)
		return
	}
	s.Require().NotZero(code, "tx succeeded, want %s", wantErr)
	s.Require().Equal(wantErr.Codespace(), codespace, "%s: %s", stage, log)
	s.Require().Equal(wantErr.ABCICode(), code, "%s: %s", stage, log)
}

// waitSanctioned waits until every address has the wanted sanction status.
func (s *IntegrationTestSuite) waitSanctioned(endpoint string, want bool, addrs ...string) {
	s.Require().Eventually(
		func() bool {
			for _, a := range addrs {
				res, err := queryIsSanctioned(endpoint, a)
				s.Require().NoError(err)
				if res.Sanctioned != want {
					return false
				}
			}
			return true
		},
		15*time.Second,
		5*time.Second,
		"want sanctioned=%v for %v", want, addrs,
	)
}

// waitSanctioners waits until the sanctioner list is exactly want.
func (s *IntegrationTestSuite) waitSanctioners(endpoint string, want map[string]sanctiontypes.SanctionerRole) {
	s.Require().Eventually(
		func() bool {
			sanctioners, err := querySanctioners(endpoint)
			s.Require().NoError(err)
			if len(sanctioners) != len(want) {
				return false
			}
			for _, sn := range sanctioners {
				if want[sn.Address] != sn.Role {
					return false
				}
			}
			return true
		},
		15*time.Second,
		5*time.Second,
		"want sanctioners %v", want,
	)
}

// listChangedEvents returns the EventSanctionListChanged events of a
// committed tx.
func (s *IntegrationTestSuite) listChangedEvents(endpoint, txHash string) []*sanctiontypes.EventSanctionListChanged {
	res, found := queryTxResult(endpoint, txHash)
	s.Require().True(found, "tx %s not found", txHash)
	return parseTypedEvents[*sanctiontypes.EventSanctionListChanged](s, res.Events)
}

// testSanctionerFastPath walks the sanctioner flow end to end on a running
// chain. It leaves the blacklist and the sanctioner list empty.
//
//  1. Governance appoints a full sanctioner (manager) and an add-only
//     sanctioner (bot).
//  2. The bot sanctions two addresses in one tx. Both are blocked and the tx
//     emits one event naming both, the bot and its role. Sanctioning them
//     again succeeds with no event.
//  3. Every rejected case fails with its exact error code.
//  4. The manager unsanctions both in one tx, with one event.
//  5. Governance sanctions the bot. The bot keeps its role but its txs are
//     rejected by the ante handler until the manager lifts the entry.
//  6. The manager revokes the bot, which can no longer sanction.
//  7. Governance removes the bot's remaining entry and empties the list,
//     turning the fast path off.
func (s *IntegrationTestSuite) testSanctionerFastPath() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	validatorA := s.chainA.validators[0]
	validatorAddr, _ := validatorA.keyInfo.GetAddress()

	valIdx := 0
	alice, _ := s.chainA.genesisAccounts[1].keyInfo.GetAddress()   // recipient, and a non-sanctioner signer
	bob, _ := s.chainA.genesisAccounts[2].keyInfo.GetAddress()     // target
	manager, _ := s.chainA.genesisAccounts[3].keyInfo.GetAddress() // full sanctioner
	bot, _ := s.chainA.genesisAccounts[4].keyInfo.GetAddress()     // add-only sanctioner
	// Addresses with no key.
	carol := sdk.AccAddress([]byte("e2e-sanction-carol__")).String()
	botEntry := sdk.AccAddress([]byte("e2e-bot-sanctioned__")).String()
	reason := "rule drain-v1, tx 0xe2e"

	// 1. Governance appoints the sanctioners.
	s.writeUpdateSanctionersProposal(s.chainA, manager.String(), bot.String())
	s.T().Logf("Submitting Gov Proposal: appoint %s (full) and %s (add-only)", manager, bot)
	s.submitSanctionProposal(chainEndpoint, validatorAddr.String(), proposalUpdateSanctionersFilename, "sanctiontypes.MsgUpdateSanctioners")
	s.waitSanctioners(chainEndpoint, map[string]sanctiontypes.SanctionerRole{
		manager.String(): sanctiontypes.SanctionerRole_SANCTIONER_ROLE_FULL,
		bot.String():     sanctiontypes.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY,
	})

	// 2. The bot sanctions bob and carol in one tx.
	s.execBankSend(s.chainA, valIdx, bob.String(), alice.String(), tokenAmount.String(), standardFees.String(), false)
	hash := s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", bob.String()+","+carol, reason, bot.String(), nil)
	s.waitSanctioned(chainEndpoint, true, bob.String(), carol)
	blacklist, err := queryBlacklist(chainEndpoint)
	s.Require().NoError(err)
	s.Require().ElementsMatch([]string{bob.String(), carol}, blacklist)
	s.Require().Equal([]*sanctiontypes.EventSanctionListChanged{{
		Addresses:  []string{bob.String(), carol},
		Sanctioned: true,
		Actor:      bot.String(),
		ActorRole:  sanctiontypes.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY,
		Reason:     reason,
	}}, s.listChangedEvents(chainEndpoint, hash))

	s.T().Logf("Bob sends while sanctioned, expect failure")
	s.execBankSend(s.chainA, valIdx, bob.String(), alice.String(), tokenAmount.String(), standardFees.String(), true)

	hash = s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", bob.String(), reason, bot.String(), nil)
	s.Require().Empty(s.listChangedEvents(chainEndpoint, hash), "already listed is a no-op")

	// 3. Rejected cases, each with its exact error.
	tooMany := make([]string, sanctiontypes.DefaultMaxAddressesPerMsg+1)
	for i := range tooMany {
		tooMany[i] = sdk.AccAddress([]byte(fmt.Sprintf("e2e-cap-%012d", i))).String()
	}
	tooLong := strings.Repeat("r", int(sanctiontypes.DefaultMaxReasonLength)+1)
	stranger := sdk.AccAddress([]byte("e2e-sanction-strangr")).String()

	s.execSanctionTx(s.chainA, valIdx, "unsanction-addresses", bob.String(), "should fail", bot.String(), sanctiontypes.ErrRoleNotPermitted)
	s.execSanctionTx(s.chainA, valIdx, "revoke-sanctioner", bot.String(), "should fail", bot.String(), sanctiontypes.ErrRoleNotPermitted)
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", manager.String(), "should fail", bot.String(), sanctiontypes.ErrProtectedAddress)
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", stranger+","+bot.String(), "should fail", manager.String(), sanctiontypes.ErrProtectedAddress)
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", stranger, "should fail", alice.String(), sanctiontypes.ErrNotSanctioner)
	s.execSanctionTx(s.chainA, valIdx, "unsanction-addresses", bob.String(), "should fail", alice.String(), sanctiontypes.ErrNotSanctioner)
	s.execSanctionTx(s.chainA, valIdx, "revoke-sanctioner", bot.String(), "should fail", alice.String(), sanctiontypes.ErrNotSanctioner)
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", stranger, tooLong, bot.String(), sanctiontypes.ErrReasonTooLong)
	s.execSanctionTx(s.chainA, valIdx, "unsanction-addresses", bob.String(), tooLong, manager.String(), sanctiontypes.ErrReasonTooLong)
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", strings.Join(tooMany, ","), "should fail", manager.String(), sanctiontypes.ErrTooManyAddresses)
	s.execSanctionTx(s.chainA, valIdx, "revoke-sanctioner", manager.String(), "should fail", manager.String(), sanctiontypes.ErrRoleNotPermitted)
	s.execSanctionTx(s.chainA, valIdx, "revoke-sanctioner", alice.String(), "should fail", manager.String(), sanctiontypes.ErrNotSanctioner)
	s.waitSanctioned(chainEndpoint, false, stranger, tooMany[0])
	s.waitSanctioned(chainEndpoint, true, bob.String(), carol)

	// 4. The manager reverses both in one tx.
	hash = s.execSanctionTx(s.chainA, valIdx, "unsanction-addresses", bob.String()+","+carol, "false positive", manager.String(), nil)
	s.waitSanctioned(chainEndpoint, false, bob.String(), carol)
	s.Require().Equal([]*sanctiontypes.EventSanctionListChanged{{
		Addresses:  []string{bob.String(), carol},
		Sanctioned: false,
		Actor:      manager.String(),
		ActorRole:  sanctiontypes.SanctionerRole_SANCTIONER_ROLE_FULL,
		Reason:     "false positive",
	}}, s.listChangedEvents(chainEndpoint, hash))
	s.execBankSend(s.chainA, valIdx, bob.String(), alice.String(), tokenAmount.String(), standardFees.String(), false)

	// 5. Governance sanctions the bot, which the ante handler then stops.
	s.writeAddBlacklistAccountsProposal(s.chainA, bot.String())
	s.T().Logf("Submitting Gov Proposal: Add sanctioner %s to blacklist", bot)
	s.submitSanctionProposal(chainEndpoint, validatorAddr.String(), proposalAddBlacklistAccountsFilename, "sanctiontypes.MsgSanctionAddresses")
	s.waitSanctioned(chainEndpoint, true, bot.String())
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", stranger, reason, bot.String(), sanctiontypes.ErrAccountBlacklisted)
	s.waitSanctioners(chainEndpoint, map[string]sanctiontypes.SanctionerRole{
		manager.String(): sanctiontypes.SanctionerRole_SANCTIONER_ROLE_FULL,
		bot.String():     sanctiontypes.SanctionerRole_SANCTIONER_ROLE_ADD_ONLY,
	})

	// The manager lifts governance's entry, and the bot works again.
	s.execSanctionTx(s.chainA, valIdx, "unsanction-addresses", bot.String(), "reviewed", manager.String(), nil)
	s.waitSanctioned(chainEndpoint, false, bot.String())
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", botEntry, reason, bot.String(), nil)
	s.waitSanctioned(chainEndpoint, true, botEntry)

	// 6. The manager revokes the bot, which can no longer sanction.
	hash = s.execSanctionTx(s.chainA, valIdx, "revoke-sanctioner", bot.String(), "key rotated", manager.String(), nil)
	s.waitSanctioners(chainEndpoint, map[string]sanctiontypes.SanctionerRole{
		manager.String(): sanctiontypes.SanctionerRole_SANCTIONER_ROLE_FULL,
	})
	res, found := queryTxResult(chainEndpoint, hash)
	s.Require().True(found)
	revoked := parseTypedEvents[*sanctiontypes.EventSanctionerRevoked](s, res.Events)
	s.Require().Len(revoked, 1)
	s.Require().Equal(bot.String(), revoked[0].Address)
	s.Require().Equal(manager.String(), revoked[0].RevokedBy)
	s.Require().Equal("key rotated", revoked[0].Reason)
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", stranger, reason, bot.String(), sanctiontypes.ErrNotSanctioner)
	s.waitSanctioned(chainEndpoint, true, botEntry)

	// 7. Governance removes the bot's entry, then empties the list.
	s.writeRemoveBlacklistAccountsProposal(s.chainA, botEntry)
	s.T().Logf("Submitting Gov Proposal: Remove %s from blacklist", botEntry)
	s.submitSanctionProposal(chainEndpoint, validatorAddr.String(), proposalRemoveBlacklistAccountsFilename, "sanctiontypes.MsgUnsanctionAddresses")

	s.writeUpdateSanctionersProposal(s.chainA, "", "")
	s.T().Logf("Submitting Gov Proposal: empty the sanctioner list")
	s.submitSanctionProposal(chainEndpoint, validatorAddr.String(), proposalUpdateSanctionersFilename, "sanctiontypes.MsgUpdateSanctioners")

	s.waitSanctioners(chainEndpoint, map[string]sanctiontypes.SanctionerRole{})
	s.Require().Eventually(
		func() bool {
			blacklist, err := queryBlacklist(chainEndpoint)
			s.Require().NoError(err)
			return len(blacklist) == 0
		},
		15*time.Second,
		5*time.Second,
	)
	s.execSanctionTx(s.chainA, valIdx, "sanction-addresses", stranger, reason, manager.String(), sanctiontypes.ErrNotSanctioner)
}

// txResult is the part of a committed tx these tests check.
type txResult struct {
	Code      uint32
	Codespace string
	RawLog    string
	Events    []abci.Event
}

// queryTxResult returns a committed tx's result, and false if the tx isn't found
// yet.
func queryTxResult(endpoint, txHash string) (*txResult, bool) {
	body, err := httpGet(fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", endpoint, txHash))
	if err != nil {
		return nil, false
	}
	var resp struct {
		TxResponse *struct {
			Code      uint32 `json:"code"`
			Codespace string `json:"codespace"`
			RawLog    string `json:"raw_log"`
			Events    []struct {
				Type       string `json:"type"`
				Attributes []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"attributes"`
			} `json:"events"`
		} `json:"tx_response"`
	}
	if err := json.Unmarshal(body, &resp); err != nil || resp.TxResponse == nil {
		return nil, false
	}
	res := &txResult{Code: resp.TxResponse.Code, Codespace: resp.TxResponse.Codespace, RawLog: resp.TxResponse.RawLog}
	for _, e := range resp.TxResponse.Events {
		ev := abci.Event{Type: e.Type}
		for _, a := range e.Attributes {
			ev.Attributes = append(ev.Attributes, abci.EventAttribute{Key: a.Key, Value: a.Value})
		}
		res.Events = append(res.Events, ev)
	}
	return res, true
}

// parseTypedEvents returns the events of type T.
func parseTypedEvents[T proto.Message](s *IntegrationTestSuite, events []abci.Event) []T {
	var zero T
	name := proto.MessageName(zero)
	out := []T{}
	for _, e := range events {
		if e.Type != name {
			continue
		}
		msg, err := sdk.ParseTypedEvent(e)
		s.Require().NoError(err)
		out = append(out, msg.(T))
	}
	return out
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
