package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	subdenom = "test"
	mintAmt  = 1000
	burnAmt  = 800
)

const (
	proposalDisableDenomSendFilename = "proposal_disable_denom_send.json"
	proposalEnableDenomSendFilename  = "proposal_enable_denom_send.json"

	transferCapContractFilename      = "transfer_cap.wasm"
	percentageCapContractFilename    = "percentage_cap.wasm"
	transferCapTrackContractFilename = "transfer_cap_track.wasm"
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

func (s *IntegrationTestSuite) writeWasmContracts(c *chain) {
	contractWasm, err := os.ReadFile(fmt.Sprint("test_data/", transferCapContractFilename))
	s.Require().NoError(err)
	err = writeFile(filepath.Join(c.validators[0].configDir(), transferCapContractFilename), contractWasm)
	s.Require().NoError(err)

	contractWasm, err = os.ReadFile(fmt.Sprint("test_data/", percentageCapContractFilename))
	s.Require().NoError(err)
	err = writeFile(filepath.Join(c.validators[0].configDir(), percentageCapContractFilename), contractWasm)
	s.Require().NoError(err)

	contractWasm, err = os.ReadFile(fmt.Sprint("test_data/", transferCapTrackContractFilename))
	s.Require().NoError(err)
	err = writeFile(filepath.Join(c.validators[0].configDir(), transferCapTrackContractFilename), contractWasm)
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
		alice := s.getAlice()

		var beforeAliceUomBalance,
			afterAliceUomBalance sdk.Coin

		denomCreationFee, err := queryTokenfactoryDenomCreationFee(chainEndpoint)
		s.Require().NoError(err)
		s.Require().Equal(denomCreationFee.Denom, uomDenom)

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				beforeAliceUomBalance, err = getSpecificBalance(chainEndpoint, alice, uomDenom)
				s.Require().NoError(err)

				return beforeAliceUomBalance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		s.createDenom(c, valIdx, alice, subdenom, standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterAliceUomBalance, err = getSpecificBalance(chainEndpoint, alice, uomDenom)
				s.Require().NoError(err)

				beforeAlice := beforeAliceUomBalance.Sub(denomCreationFee).Sub(standardFees)

				return beforeAlice.Equal(afterAliceUomBalance)
			},
			10*time.Second,
			5*time.Second,
		)
	})
}

func (s *IntegrationTestSuite) testTokenfactoryAdmin() {
	s.Run("default_denom_admin_should_be_the_creator", func() {
		var (
			err    error
			valIdx = 0
			c      = s.chainA
		)

		chainEndpoint := fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))

		result, err := queryTokenfactoryDenomAuthorityMetadata(chainEndpoint, s.getAlice(), subdenom)

		s.Require().NoError(err)

		s.Require().Equal(result.Admin, s.getAlice(), "By default, the denom admin should be the creator")
	})
}

func (s *IntegrationTestSuite) testTokenfactorySetMetadata() {
	s.Run("set_denom_metadata_tokenfactory", func() {
		var (
			err    error
			valIdx = 0
			c      = s.chainA
		)

		// Create metadata JSON content using the global helper function
		metadataContent := s.BuildTokenMetadata()
		customDenom := metadataContent.Base
		metadataString := MetadataToString(metadataContent)

		// Write metadata to file in the validator's config directory
		metadataFileName := "metadata.json"
		metadataFile := filepath.Join(c.validators[valIdx].configDir(), metadataFileName)
		err = writeFile(metadataFile, []byte(metadataString))
		s.Require().NoError(err)
		s.T().Logf("Start setting metadata for denom %s", customDenom)

		// Set the metadata using the CLI command
		s.setDenomMetadata(c, valIdx, s.getAlice(), filepath.Join(mantraHomePath, metadataFileName), standardFees.String(), false)

		s.T().Logf("Successfully set metadata for denom %s", customDenom)

		// Query and verify the metadata was set correctly
		chainEndpoint := fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))

		s.Require().Eventually(
			func() bool {
				queriedMetadata, err := queryTokenfactoryDenomMetadata(chainEndpoint, customDenom)
				if err != nil {
					s.T().Logf("Error querying metadata: %v", err)
					return false
				}

				// Verify metadata fields
				if queriedMetadata.Description != metadataContent.Description {
					s.T().Logf("Description mismatch: expected %s, got %s", metadataContent.Description, queriedMetadata.Description)
					return false
				}

				if queriedMetadata.Base != metadataContent.Base {
					s.T().Logf("Base mismatch: expected %s, got %s", metadataContent.Base, queriedMetadata.Base)
					return false
				}

				if queriedMetadata.Display != metadataContent.Display {
					s.T().Logf("Display mismatch: expected %s, got %s", metadataContent.Display, queriedMetadata.Display)
					return false
				}

				if queriedMetadata.Name != metadataContent.Name {
					s.T().Logf("Name mismatch: expected %s, got %s", metadataContent.Name, queriedMetadata.Name)
					return false
				}

				if queriedMetadata.Symbol != metadataContent.Symbol {
					s.T().Logf("Symbol mismatch: expected %s, got %s", metadataContent.Symbol, queriedMetadata.Symbol)
					return false
				}

				if len(queriedMetadata.DenomUnits) != len(metadataContent.DenomUnits) {
					s.T().Logf("DenomUnits length mismatch: expected %d, got %d", len(metadataContent.DenomUnits), len(queriedMetadata.DenomUnits))
					return false
				}

				// Verify denom units
				for i, expectedUnit := range metadataContent.DenomUnits {
					if i >= len(queriedMetadata.DenomUnits) {
						s.T().Logf("Missing denom unit at index %d", i)
						return false
					}

					queriedUnit := queriedMetadata.DenomUnits[i]
					if queriedUnit.Denom != expectedUnit.Denom {
						s.T().Logf("DenomUnit[%d] Denom mismatch: expected %s, got %s", i, expectedUnit.Denom, queriedUnit.Denom)
						return false
					}

					if queriedUnit.Exponent != expectedUnit.Exponent {
						s.T().Logf("DenomUnit[%d] Exponent mismatch: expected %d, got %d", i, expectedUnit.Exponent, queriedUnit.Exponent)
						return false
					}

					if len(queriedUnit.Aliases) != len(expectedUnit.Aliases) {
						s.T().Logf("DenomUnit[%d] Aliases length mismatch: expected %d, got %d", i, len(expectedUnit.Aliases), len(queriedUnit.Aliases))
						return false
					}

					for j, expectedAlias := range expectedUnit.Aliases {
						if j >= len(queriedUnit.Aliases) || queriedUnit.Aliases[j] != expectedAlias {
							s.T().Logf("DenomUnit[%d] Alias[%d] mismatch: expected %s, got %s", i, j, expectedAlias, queriedUnit.Aliases[j])
							return false
						}
					}
				}

				s.T().Logf("Successfully verified metadata for denom %s", customDenom)
				return true
			},
			30*time.Second,
			2*time.Second,
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
		alice := s.getAlice()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()

		var beforeAliceCustomTokenBalance,
			afterAliceCustomTokenBalance,
			beforeBobCustomTokenBalance,
			afterBobCustomTokenBalance sdk.Coin

		customDenom := buildDenom(alice, subdenom)

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				beforeAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice, customDenom)
				s.Require().NoError(err)

				beforeBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				return beforeAliceCustomTokenBalance.IsValid() && beforeBobCustomTokenBalance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		toMint := sdk.NewCoin(customDenom, math.NewInt(mintAmt))
		s.mintDenom(c, valIdx, alice, toMint.String(), "", standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice, customDenom)
				s.Require().NoError(err)

				incrementedAlice := beforeAliceCustomTokenBalance.Add(toMint)

				return incrementedAlice.Equal(afterAliceCustomTokenBalance)
			},
			10*time.Second,
			5*time.Second,
		)

		s.mintDenom(c, valIdx, alice, toMint.String(), bob.String(), standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				incrementedBob := beforeBobCustomTokenBalance.Add(toMint)

				return incrementedBob.Equal(afterBobCustomTokenBalance)
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

		s.mintDenom(c, valIdx, alice, toMint.String(), bob.String(), standardFees.String(), true)

		escrowAddress, err := queryIBCEscrowAddress(chainEndpoint, "channel-0")
		s.Require().NoError(err)
		s.mintDenom(c, valIdx, alice, toMint.String(), escrowAddress, standardFees.String(), true)
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
		alice := s.getAlice()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()

		var beforeAliceCustomTokenBalance,
			afterAliceCustomTokenBalance,
			beforeBobCustomTokenBalance,
			afterBobCustomTokenBalance sdk.Coin

		customDenom := buildDenom(alice, subdenom)

		// get balances of sender and recipient accounts
		s.Require().Eventually(
			func() bool {
				beforeAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice, customDenom)
				s.Require().NoError(err)

				beforeBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				return beforeAliceCustomTokenBalance.IsValid() && beforeBobCustomTokenBalance.IsValid()
			},
			10*time.Second,
			5*time.Second,
		)

		toBurn := sdk.NewCoin(customDenom, math.NewInt(burnAmt))

		s.T().Logf("Reenable token send %s", customDenom)
		s.burnDenom(c, valIdx, alice, toBurn.String(), "", standardFees.String(), true)

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

		s.burnDenom(c, valIdx, alice, toBurn.String(), "", standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice, customDenom)
				s.Require().NoError(err)

				beforeAlice := beforeAliceCustomTokenBalance.Sub(toBurn)

				return beforeAlice.Equal(afterAliceCustomTokenBalance)
			},
			10*time.Second,
			5*time.Second,
		)

		s.burnDenom(c, valIdx, alice, toBurn.String(), bob.String(), standardFees.String(), false)

		// check that the creation was successful
		s.Require().Eventually(
			func() bool {
				afterBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				beforeBob := beforeBobCustomTokenBalance.Sub(toBurn)

				return beforeBob.Equal(afterBobCustomTokenBalance)
			},
			10*time.Second,
			5*time.Second,
		)

		escrowAddress, err := queryIBCEscrowAddress(chainEndpoint, "channel-0")
		s.Require().NoError(err)
		s.burnDenom(c, valIdx, alice, toBurn.String(), escrowAddress, standardFees.String(), true)
	})
}

func (s *IntegrationTestSuite) testTokenfactoryHooks() {
	const (
		TRANSFER_CAP = int64(1000000)
	)
	var (
		err                 error
		valIdx              = 0
		c                   = s.chainA
		chainEndpoint       = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		testHooksMintAmount = int64(10000000)
	)

	// define one admin and one recipient
	alice, _ := c.genesisAccounts[1].keyInfo.GetAddress()
	bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()
	charlie, _ := c.genesisAccounts[3].keyInfo.GetAddress()

	var beforeAliceCustomTokenBalance,
		afterAliceCustomTokenBalance,
		beforeBobCustomTokenBalance,
		afterBobCustomTokenBalance sdk.Coin

	var transferCapContractCode,
		percentageCapContractCode,
		transferCapTrackContractCode int

	var transferCapContractAddr,
		percentageCapContractAddr,
		transferCapTrackContractAddr string

	customDenom := buildDenom(charlie.String(), subdenom)
	toMint := sdk.NewCoin(customDenom, math.NewInt(testHooksMintAmount))

	s.writeWasmContracts(c)

	s.Run("setup_hooks_denom", func() {
		s.createDenom(c, valIdx, charlie.String(), subdenom, standardFees.String(), false)

		s.mintDenom(c, valIdx, charlie.String(), toMint.String(), bob.String(), standardFees.String(), false)
	})

	s.Run("store_and_instantiate_hook_contracts", func() {
		initialCodes, err := queryWasmCodes(chainEndpoint)
		s.Require().NoError(err)
		initialCodeCount := len(initialCodes.CodeInfos)

		s.execWasmStoreCode(s.chainA, 0, charlie.String(),
			filepath.Join(mantraHomePath, transferCapContractFilename), mantraHomePath,
		)

		// Verify the code was stored by checking the count increased
		s.Require().Eventually(
			func() bool {
				updatedCodes, err := queryWasmCodes(chainEndpoint)
				s.Require().NoError(err)
				transferCapContractCode = len(updatedCodes.CodeInfos)
				return len(updatedCodes.CodeInfos) == initialCodeCount+1
			},
			30*time.Second,
			2*time.Second,
		)

		initialCodeCount = transferCapContractCode

		s.execWasmStoreCode(s.chainA, 0, charlie.String(),
			filepath.Join(mantraHomePath, percentageCapContractFilename), mantraHomePath,
		)

		// Verify the code was stored by checking the count increased
		s.Require().Eventually(
			func() bool {
				updatedCodes, err := queryWasmCodes(chainEndpoint)
				s.Require().NoError(err)
				percentageCapContractCode = len(updatedCodes.CodeInfos)
				return len(updatedCodes.CodeInfos) == initialCodeCount+1
			},
			30*time.Second,
			2*time.Second,
		)

		initialCodeCount = percentageCapContractCode

		s.execWasmStoreCode(s.chainA, 0, charlie.String(),
			filepath.Join(mantraHomePath, transferCapTrackContractFilename), mantraHomePath,
		)

		// Verify the code was stored by checking the count increased
		s.Require().Eventually(
			func() bool {
				updatedCodes, err := queryWasmCodes(chainEndpoint)
				s.Require().NoError(err)
				transferCapTrackContractCode = len(updatedCodes.CodeInfos)
				return len(updatedCodes.CodeInfos) == initialCodeCount+1
			},
			30*time.Second,
			2*time.Second,
		)

		initMsg := "{}"
		label := "test"

		txHash := s.execWasmInstantiate(
			s.chainA,
			0,
			charlie.String(),
			uint64(transferCapContractCode),
			initMsg,
			label,
			charlie.String(),
			"",
			mantraHomePath,
		)
		s.Require().NotEmpty(txHash)
		s.T().Logf("Instantiation transaction submitted with hash: %s", txHash)

		// Query transaction events to get contract address
		events, err := queryTxEvents(chainEndpoint, txHash)
		s.Require().NoError(err)
		addr, err := findContractAddressFromEvents(events)
		s.Require().NoError(err)
		s.NotEmpty(addr)
		transferCapContractAddr = addr
		s.T().Logf("Successfully instantiated contract at address: %s", transferCapContractAddr)

		txHash = s.execWasmInstantiate(
			s.chainA,
			0,
			charlie.String(),
			uint64(percentageCapContractCode),
			initMsg,
			label,
			charlie.String(),
			"",
			mantraHomePath,
		)
		s.Require().NotEmpty(txHash)
		s.T().Logf("Instantiation transaction submitted with hash: %s", txHash)

		// Query transaction events to get contract address
		events, err = queryTxEvents(chainEndpoint, txHash)
		s.Require().NoError(err)
		addr, err = findContractAddressFromEvents(events)
		s.Require().NoError(err)
		s.NotEmpty(addr)
		percentageCapContractAddr = addr
		s.T().Logf("Successfully instantiated contract at address: %s", percentageCapContractAddr)

		txHash = s.execWasmInstantiate(
			s.chainA,
			0,
			charlie.String(),
			uint64(transferCapTrackContractCode),
			initMsg,
			label,
			charlie.String(),
			"",
			mantraHomePath,
		)
		s.Require().NotEmpty(txHash)
		s.T().Logf("Instantiation transaction submitted with hash: %s", txHash)

		// Query transaction events to get contract address
		events, err = queryTxEvents(chainEndpoint, txHash)
		s.Require().NoError(err)
		addr, err = findContractAddressFromEvents(events)
		s.Require().NoError(err)
		s.NotEmpty(addr)
		transferCapTrackContractAddr = addr
		s.T().Logf("Successfully instantiated contract at address: %s", transferCapTrackContractAddr)
	})

	s.Run("transfer_cap_hook_test", func() {
		s.setBeforeSendHook(c, valIdx, charlie.String(), customDenom, transferCapContractAddr, standardFees.String(), false)

		toSendFail := sdk.NewCoin(customDenom, math.NewInt(TRANSFER_CAP+1))
		s.execBankSend(c, valIdx, bob.String(), alice.String(), toSendFail.String(), standardFees.String(), true)
		s.T().Log("Fail to send over transfer cap")

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

		toSendSucceed := sdk.NewCoin(customDenom, math.NewInt(TRANSFER_CAP))
		s.execBankSend(c, valIdx, bob.String(), alice.String(), toSendSucceed.String(), standardFees.String(), false)

		s.Require().Eventually(
			func() bool {
				afterAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice.String(), customDenom)
				s.Require().NoError(err)
				afterBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				incrementedAlice := beforeAliceCustomTokenBalance.Add(toSendSucceed)
				decrementedBob := beforeBobCustomTokenBalance.Sub(toSendSucceed)

				return incrementedAlice.Equal(afterAliceCustomTokenBalance) && decrementedBob.Equal(afterBobCustomTokenBalance)
			},
			10*time.Second,
			5*time.Second,
		)
		s.T().Log("Succeed send below transfer cap")
	})

	s.Run("percentage_cap_hook_test", func() {
		s.setBeforeSendHook(c, valIdx, charlie.String(), customDenom, percentageCapContractAddr, standardFees.String(), false)

		percentage_cap := toMint.Amount.Quo(math.NewInt(2))
		toSendFail := sdk.NewCoin(customDenom, percentage_cap.Add(math.OneInt()))
		s.execBankSend(c, valIdx, bob.String(), alice.String(), toSendFail.String(), standardFees.String(), true)
		s.T().Log("Fail to send over percentage cap")

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

		toSendSucceed := sdk.NewCoin(customDenom, percentage_cap)
		s.execBankSend(c, valIdx, bob.String(), alice.String(), toSendSucceed.String(), standardFees.String(), false)

		s.Require().Eventually(
			func() bool {
				afterAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice.String(), customDenom)
				s.Require().NoError(err)
				afterBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
				s.Require().NoError(err)

				incrementedAlice := beforeAliceCustomTokenBalance.Add(toSendSucceed)
				decrementedBob := beforeBobCustomTokenBalance.Sub(toSendSucceed)

				return incrementedAlice.Equal(afterAliceCustomTokenBalance) && decrementedBob.Equal(afterBobCustomTokenBalance)
			},
			10*time.Second,
			5*time.Second,
		)
		s.T().Log("Succeed send below percentage cap")
	})

	// s.Run("transfer_cap_hook_track_test", func() {
	// 	s.setBeforeSendHook(c, valIdx, charlie.String(), customDenom, transferCapTrackContractAddr, standardFees.String(), false)

	// 	// get balances of sender and recipient accounts
	// 	s.Require().Eventually(
	// 		func() bool {
	// 			beforeAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice.String(), customDenom)
	// 			s.Require().NoError(err)

	// 			beforeBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
	// 			s.Require().NoError(err)

	// 			return beforeAliceCustomTokenBalance.IsValid() && beforeBobCustomTokenBalance.IsValid()
	// 		},
	// 		10*time.Second,
	// 		5*time.Second,
	// 	)

	// 	toSend := sdk.NewCoin(customDenom, math.NewInt(TRANSFER_CAP+1))
	// 	s.execBankSend(c, valIdx, bob.String(), alice.String(), toSend.String(), standardFees.String(), false)
	// 	s.T().Log("Succeed to send although over transfer cap")

	// 	s.Require().Eventually(
	// 		func() bool {
	// 			afterAliceCustomTokenBalance, err = getSpecificBalance(chainEndpoint, alice.String(), customDenom)
	// 			s.Require().NoError(err)
	// 			afterBobCustomTokenBalance, err = getSpecificBalance(chainEndpoint, bob.String(), customDenom)
	// 			s.Require().NoError(err)

	// 			incrementedAlice := beforeAliceCustomTokenBalance.Add(toSend)
	// 			decrementedBob := beforeBobCustomTokenBalance.Sub(toSend)

	// 			return incrementedAlice.Equal(afterAliceCustomTokenBalance) && decrementedBob.Equal(afterBobCustomTokenBalance)
	// 		},
	// 		10*time.Second,
	// 		5*time.Second,
	// 	)
	// 	s.T().Log("Succeed send above transfer cap as it is only TrackBeforeSend")
	// })
}

func buildDenom(sender, subDenom string) string {
	return fmt.Sprintf("factory/%s/%s", sender, subDenom)
}

func (s *IntegrationTestSuite) getAlice() string {
	alice, err := s.chainA.genesisAccounts[1].keyInfo.GetAddress()
	s.Require().NoError(err)
	return alice.String()
}

func (s *IntegrationTestSuite) BuildTokenMetadata() banktypes.Metadata {
	factoryDenom := buildDenom(s.getAlice(), subdenom)

	symbol := cases.Upper(language.English).String(subdenom)
	name := cases.Title(language.English).String(subdenom)
	metadata := banktypes.Metadata{
		Description: fmt.Sprintf("%s token for tokenfactory e2e tests", name),
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    factoryDenom,
				Exponent: 0,
				Aliases:  []string{"test_alias"},
			},
			{
				Denom:    symbol,
				Exponent: 6,
				Aliases:  []string{},
			},
		},
		Base:    factoryDenom,
		Display: symbol,
		Name:    fmt.Sprintf("%s Token", name),
		Symbol:  symbol,
	}

	return metadata
}

func MetadataToString(metadata banktypes.Metadata) string {
	metadataBytes, err := json.MarshalIndent(metadata, "", "\t")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal metadata: %v", err))
	}
	return string(metadataBytes)
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
		fmt.Sprintf("--%s=%s", flags.FlagGas, "auto"),
		fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.5"),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, c.id),
		"--keyring-backend=test",
		"--broadcast-mode=sync",
		"--output=json",
		"-y",
	}
	denom := buildDenom(sender, subdenom)
	s.T().Logf("%s is creating tokenfactory denom %s", sender, denom)
	if expErr {
		s.executeTxCommand(ctx, c, ibcCmd, valIdx, s.expectErrExecValidation(c, valIdx, true))
		s.T().Log("create tokenfactory denom unsuccessful")
	} else {
		s.executeTxCommand(ctx, c, ibcCmd, valIdx, s.defaultExecValidation(c, valIdx))
		s.T().Log("successfully created tokenfactory denom")
	}
}

func (s *IntegrationTestSuite) setDenomMetadata(c *chain, valIdx int, sender, metadataFile, fees string, expErr bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Sample tx: https://mantrascan.io/dukong/tx/4f40cc08aadb5ca005a4138353c707b1398858c577186458d5ce2a70bd3a67c8
	cmd := []string{
		mantrachaindBinary,
		txCommand,
		"tokenfactory",
		"set-denom-metadata",
		metadataFile,
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

	s.T().Logf("Address %s is setting denom metadata from file %s", sender, metadataFile)
	if expErr {
		s.executeTxCommand(ctx, c, cmd, valIdx, s.expectErrExecValidation(c, valIdx, true))
		s.T().Log("set denom metadata unsuccessful")
	} else {
		s.executeTxCommand(ctx, c, cmd, valIdx, s.defaultExecValidation(c, valIdx))
		s.T().Log("successfully set denom metadata")
	}
}

func (s *IntegrationTestSuite) setBeforeSendHook(c *chain, valIdx int, sender, customDenom, contractAddr, fees string, expErr bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Sample tx: https://mantrascan.io/dukong/tx/4f40cc08aadb5ca005a4138353c707b1398858c577186458d5ce2a70bd3a67c8
	cmd := []string{
		mantrachaindBinary,
		txCommand,
		"tokenfactory",
		"set-before-send-hook",
		customDenom,
		contractAddr,
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

	s.T().Logf("Address %s is setting before send hook for denom %s to contract address %s", sender, customDenom, contractAddr)
	if expErr {
		s.executeTxCommand(ctx, c, cmd, valIdx, s.expectErrExecValidation(c, valIdx, true))
		s.T().Log("set before-send-hook unsuccessful")
	} else {
		s.executeTxCommand(ctx, c, cmd, valIdx, s.defaultExecValidation(c, valIdx))
		s.T().Log("successfully set before-send-hook")
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

func (s *IntegrationTestSuite) burnDenom(c *chain, valIdx int, sender, burnCoin, burnFrom, fees string, expErr bool) { //nolint:unparam
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
