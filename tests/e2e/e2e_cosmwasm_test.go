package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
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
		contractWasm, err := os.ReadFile("test_data/contract_1.wasm")
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
		s.execWasmStoreCode(s.chainA, 0, senderAddr,
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
	})
}
