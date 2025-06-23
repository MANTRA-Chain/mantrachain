package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

var (
	deployedWasmCodeId      uint64 = 0
	deployedContractAddress string
)

func (s *IntegrationTestSuite) testQueryWasmParams() {
	s.Run("query_wasm_params", func() {
		chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

		params, err := queryWasmParams(chainEndpoint)
		s.Require().NoError(err)
		s.Require().Equal(params.CodeUploadAccess.Permission, wasmTypes.AccessTypeEverybody)
		s.Require().Equal(params.InstantiateDefaultPermission, wasmTypes.AccessTypeEverybody)
	})
}

func (s *IntegrationTestSuite) testStoreCode() {
	s.Run("store_wasm_code", func() {
		chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

		// Read the contract wasm file
		contractWasm, err := os.ReadFile("test_data/rwa_oracle.wasm")
		s.Require().NoError(err)

		// Get the initial count of stored codes
		initialCodes, err := queryWasmCodes(chainEndpoint)
		s.Require().NoError(err)
		initialCodeCount := len(initialCodes.CodeInfos)

		// Get validator address for sending transaction
		valAddr, _ := s.chainA.validators[0].keyInfo.GetAddress()
		senderAddr := valAddr.String()

		// Write the contract file to the validator's file system
		contractFileName := "contract_store_test.wasm"
		err = writeFile(filepath.Join(s.chainA.validators[0].configDir(), contractFileName), contractWasm)
		s.Require().NoError(err)

		// Store the code using wasm store command
		txHash := s.execWasmStoreCode(s.chainA, 0, senderAddr,
			filepath.Join(mantraHomePath, contractFileName), mantraHomePath,
		)

		// Verify the code was stored by checking the count increased
		s.Require().Eventually(
			func() bool {
				updatedCodes, err := queryWasmCodes(chainEndpoint)
				s.Require().NoError(err)
				return len(updatedCodes.CodeInfos) == initialCodeCount+1
			},
			30*time.Second,
			2*time.Second,
		)

		// Get the latest stored code info
		finalCodes, err := queryWasmCodes(chainEndpoint)
		s.Require().NoError(err)
		s.Require().Greater(len(finalCodes.CodeInfos), initialCodeCount)

		// Find the newly stored code (should be the one with the highest code_id)
		var newestCode *wasmTypes.CodeInfoResponse
		maxCodeID := uint64(0)
		for _, codeInfo := range finalCodes.CodeInfos {
			if codeInfo.CodeID > maxCodeID {
				maxCodeID = codeInfo.CodeID
				newestCode = &codeInfo
			}
		}

		s.Require().NotNil(newestCode, "Should have found the newly stored code")
		s.Require().Equal(senderAddr, newestCode.Creator)
		s.Require().Greater(newestCode.CodeID, uint64(0))

		event, err := queryTxEvents(chainEndpoint, txHash)
		s.Require().NoError(err)

		codeID, err := findCodeIdFromEvents(event)
		s.Require().NoError(err)
		s.Require().Greater(codeID, uint64(0))

		// Store the code ID for potential use in other tests
		deployedWasmCodeId = codeID
		s.T().Logf("Successfully stored wasm code with ID: %d", codeID)
	})
}

func (s *IntegrationTestSuite) testInstantiateContract() {
	s.Run("instantiate_wasm_contract", func() {
		// Skip if no code has been deployed
		if deployedWasmCodeId == 0 {
			s.T().Skip("No wasm code uploaded, skipping instantiate test")
		}

		chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

		// Get validator address for sending transaction
		valAddr, _ := s.chainA.validators[0].keyInfo.GetAddress()
		senderAddr := valAddr.String()

		// Simple init message for most contracts
		initMsg := `{}`
		label := "rwa_oracle"
		var contractAddr string
		var txHash string

		s.T().Logf("Trying instantiation with init message: %s", initMsg)

		// Try to instantiate the contract using execWasmInstantiate
		func() {
			defer func() {
				if r := recover(); r != nil {
					s.T().Logf("Instantiation failed with panic: %v", r)
				}
			}()

			txHash = s.execWasmInstantiate(
				s.chainA,
				0,
				senderAddr,
				deployedWasmCodeId,
				initMsg,
				label,
				senderAddr,
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
			contractAddr = addr
			s.T().Logf("Successfully instantiated contract at address: %s with init message: %s", contractAddr, initMsg)
		}()

		// Update the global variable regardless of success
		s.Require().NotEmpty(contractAddr)
		deployedContractAddress = contractAddr
		s.T().Logf("Contract instantiation successful. Address: %s", deployedContractAddress)

		contractInfo, err := queryWasmContractInfo(chainEndpoint, contractAddr)
		s.Require().NoError(err)
		s.Require().Equal(contractInfo.CodeID, deployedWasmCodeId)

		s.T().Log("Instantiation test completed")
	})
}

func (s *IntegrationTestSuite) testExecuteContractWithSimplyMessage() {
	s.Run("execute_wasm_contract_with_simple_message", func() {
		// Skip if no contract has been instantiated
		if deployedContractAddress == "" {
			s.T().Skip("No contract deployed, skipping execute test")
		}

		//chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

		// Get validator address for sending transaction
		valAddr, _ := s.chainA.validators[0].keyInfo.GetAddress()
		senderAddr := valAddr.String()

		// Simple message to execute on the contract
		execMsg := `{ "add_publishers": { "publishers": [ "mantra1hze5xhd5d5mwysddrutmdt7f89lztrx2xm3nk8" ] } }`

		txHash := s.execWasmExecute(
			s.chainA,
			0,
			senderAddr,
			deployedContractAddress,
			execMsg,
			mantraHomePath,
		)

		s.T().Logf("Execution transaction submitted with hash: %s", txHash)

		s.Require().NotEmpty(txHash)
	})
}

func findCodeIdFromEvents(events map[string][]string) (uint64, error) {
	// Look for store_code event
	if storeCodeAttrs, exists := events["store_code"]; exists {
		for _, attr := range storeCodeAttrs {
			// Each attribute is in "key=value" format
			parts := strings.Split(attr, "=")
			if len(parts) == 2 && parts[0] == "code_id" {
				var codeID uint64
				_, err := fmt.Sscanf(parts[1], "%d", &codeID)
				if err != nil {
					return 0, fmt.Errorf("failed to parse code ID: %w", err)
				}
				return codeID, nil
			}
		}
	}
	return 0, fmt.Errorf("code ID not found in events")
}

func findContractAddressFromEvents(events map[string][]string) (string, error) {
	// Look for instantiate event
	if instantiateAttrs, exists := events["instantiate"]; exists {
		for _, attr := range instantiateAttrs {
			// Each attribute is in "key=value" format
			parts := strings.Split(attr, "=")
			if len(parts) == 2 && parts[0] == "_contract_address" {
				return parts[1], nil
			}
		}
	}
	return "", fmt.Errorf("contract address not found in events")
}
