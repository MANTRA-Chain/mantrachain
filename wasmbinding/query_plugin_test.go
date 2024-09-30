package wasmbinding_test

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	proto "github.com/golang/protobuf/proto" //nolint:staticcheck // we're intentionally using this deprecated package to be compatible with cosmos protos
	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/osmomath"

	"github.com/MANTRA-Chain/mantrachain/app"

	"github.com/MANTRA-Chain/mantrachain/wasmbinding"
)

type StargateTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	app     *app.App
	HomeDir string
}

func (suite *StargateTestSuite) TearDownTestInternal() {
	os.RemoveAll(suite.HomeDir)
}

func TestStargateTestSuite(t *testing.T) {
	suite.Run(t, new(StargateTestSuite))
}

func (suite *StargateTestSuite) TestConvertProtoToJsonMarshal() {
	testCases := []struct {
		name                  string
		queryPath             string
		protoResponseStruct   proto.Message
		originalResponse      string
		expectedProtoResponse proto.Message
		expectedError         bool
	}{
		{
			name:                "successful conversion from proto response to json marshalled response",
			queryPath:           "/cosmos.bank.v1beta1.Query/AllBalances",
			originalResponse:    "0a090a036261721202333012050a03666f6f",
			protoResponseStruct: &banktypes.QueryAllBalancesResponse{},
			expectedProtoResponse: &banktypes.QueryAllBalancesResponse{
				Balances: sdk.NewCoins(sdk.NewCoin("bar", osmomath.NewInt(30))),
				Pagination: &query.PageResponse{
					NextKey: []byte("foo"),
				},
			},
		},
		{
			name:                "invalid proto response struct",
			queryPath:           "/cosmos.bank.v1beta1.Query/AllBalances",
			originalResponse:    "0a090a036261721202333012050a03666f6f",
			protoResponseStruct: &epochtypes.QueryCurrentEpochResponse{},
			expectedError:       true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTestInternal()
			defer suite.TearDownTestInternal()

			originalVersionBz, err := hex.DecodeString(tc.originalResponse)
			suite.Require().NoError(err)

			jsonMarshalledResponse, err := wasmbinding.ConvertProtoToJSONMarshal(tc.protoResponseStruct, originalVersionBz, suite.app.AppCodec())
			if tc.expectedError {
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)

			// check response by json marshalling proto response into json response manually
			jsonMarshalExpectedResponse, err := suite.app.AppCodec().MarshalJSON(tc.expectedProtoResponse)
			suite.Require().NoError(err)
			suite.Require().Equal(jsonMarshalledResponse, jsonMarshalExpectedResponse)
		})
	}
}

// TestDeterministicJsonMarshal tests that we get deterministic JSON marshalled response upon
// proto struct update in the state machine.
func (suite *StargateTestSuite) TestDeterministicJsonMarshal() {
	testCases := []struct {
		name                string
		testSetup           func()
		originalResponsebz  []byte
		updatedResponsebz   []byte
		queryPath           string
		responseProtoStruct interface{}
		expectedProto       func() proto.Message
	}{
		/**
		   * Origin Response
		   * balances:<denom:"bar" amount:"30" > pagination:<next_key:"foo" >
		   * New Version Response
		   * The binary built from the proto response with additional field address
		   * balances:<denom:"bar" amount:"30" > pagination:<next_key:"foo" > address:"cosmos1j6j5tsquq2jlw2af7l3xekyaq7zg4l8jsufu78"
		   // Origin proto
		   message QueryAllBalancesResponse {
		  	// balances is the balances of all the coins.
		  	repeated cosmos.base.v1beta1.Coin balances = 1
		  	[(gogoproto.nullable) = false, (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"];
		  	// pagination defines the pagination in the response.
		  	cosmos.base.query.v1beta1.PageResponse pagination = 2;
		  }
		  // Updated proto
		  message QueryAllBalancesResponse {
		  	// balances is the balances of all the coins.
		  	repeated cosmos.base.v1beta1.Coin balances = 1
		  	[(gogoproto.nullable) = false, (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"];
		  	// pagination defines the pagination in the response.
		  	cosmos.base.query.v1beta1.PageResponse pagination = 2;
		  	// address is the address to query all balances for.
		  	string address = 3;
		  }
		*/
		{
			"Query All Balances",
			func() {
				wasmbinding.SetWhitelistedQuery("/cosmos.bank.v1beta1.Query/AllBalances", &banktypes.QueryAllBalancesResponse{})
			},
			[]byte{10, 9, 10, 3, 98, 97, 114, 18, 2, 51, 48, 18, 5, 10, 3, 102, 111, 111},
			[]byte{
				10, 9, 10, 3, 98, 97, 114, 18, 2, 51, 48, 18, 5, 10, 3, 102, 111, 111, 26, 45, 99, 111, 115, 109, 111, 115, 49, 106,
				54, 106, 53, 116, 115, 113, 117, 113, 50, 106, 108, 119, 50, 97, 102, 55, 108, 51, 120, 101, 107, 121, 97, 113, 55, 122, 103,
				52, 108, 56, 106, 115, 117, 102, 117, 55, 56,
			},
			"/cosmos.bank.v1beta1.Query/AllBalances",
			&banktypes.QueryAllBalancesResponse{},
			func() proto.Message {
				return &banktypes.QueryAllBalancesResponse{
					Balances: sdk.NewCoins(sdk.NewCoin("bar", osmomath.NewInt(30))),
					Pagination: &query.PageResponse{
						NextKey: []byte("foo"),
					},
				}
			},
		},
		/**
		  // Origin proto
		  message QueryAccountResponse {
		    // account defines the account of the corresponding address.
		    google.protobuf.Any account = 1 [(cosmos_proto.accepts_interface) = "AccountI"];
		  }
		  // Updated proto
		  message QueryAccountResponse {
		    // account defines the account of the corresponding address.
		    google.protobuf.Any account = 1 [(cosmos_proto.accepts_interface) = "AccountI"];
		    // address is the address to query for.
		  	string address = 2;
		  }
		*/
		{
			"Query Account",
			nil,
			[]byte{
				10, 83, 10, 32, 47, 99, 111, 115, 109, 111, 115, 46, 97, 117, 116, 104, 46, 118, 49, 98, 101, 116, 97, 49, 46, 66, 97, 115,
				101, 65, 99, 99, 111, 117, 110, 116, 18, 47, 10, 45, 99, 111, 115, 109, 111, 115, 49, 102, 56, 117, 120, 117, 108, 116, 110, 56,
				115, 113, 122, 104, 122, 110, 114, 115, 122, 51, 113, 55, 55, 120, 119, 97, 113, 117, 104, 103, 114, 115, 103, 54, 106, 121, 118, 102, 121,
			},
			[]byte{
				10, 83, 10, 32, 47, 99, 111, 115, 109, 111, 115, 46, 97, 117, 116, 104, 46, 118, 49, 98, 101, 116, 97, 49, 46, 66, 97, 115,
				101, 65, 99, 99, 111, 117, 110, 116, 18, 47, 10, 45, 99, 111, 115, 109, 111, 115, 49, 102, 56, 117, 120, 117, 108, 116, 110, 56,
				115, 113, 122, 104, 122, 110, 114, 115, 122, 51, 113, 55, 55, 120, 119, 97, 113, 117, 104, 103, 114, 115, 103, 54, 106, 121, 118, 102, 121,
				18, 45, 99, 111, 115, 109, 111, 115, 49, 102, 56, 117, 120, 117, 108, 116, 110, 56, 115, 113, 122, 104, 122, 110, 114, 115, 122, 51, 113, 55,
				55, 120, 119, 97, 113, 117, 104, 103, 114, 115, 103, 54, 106, 121, 118, 102, 121,
			},
			"/cosmos.auth.v1beta1.Query/Account",
			&authtypes.QueryAccountResponse{},
			func() proto.Message {
				account := authtypes.BaseAccount{
					Address: "cosmos1f8uxultn8sqzhznrsz3q77xwaquhgrsg6jyvfy",
				}
				accountResponse, err := codectypes.NewAnyWithValue(&account)
				suite.Require().NoError(err)
				return &authtypes.QueryAccountResponse{
					Account: accountResponse,
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTestInternal()
			defer suite.TearDownTestInternal()

			if tc.testSetup != nil {
				tc.testSetup()
			}

			binding, err := wasmbinding.GetWhitelistedQuery(tc.queryPath)
			suite.Require().Nil(err)

			suite.Require().NoError(err)
			jsonMarshalledOriginalBz, err := wasmbinding.ConvertProtoToJSONMarshal(binding, tc.originalResponsebz, suite.app.AppCodec())
			suite.Require().NoError(err)

			jsonMarshalledUpdatedBz, err := wasmbinding.ConvertProtoToJSONMarshal(binding, tc.updatedResponsebz, suite.app.AppCodec())
			suite.Require().NoError(err)

			// json marshalled bytes should be the same since we use the same proto struct for unmarshalling
			suite.Require().Equal(jsonMarshalledOriginalBz, jsonMarshalledUpdatedBz)

			// raw build also make same result
			jsonMarshalExpectedResponse, err := suite.app.AppCodec().MarshalJSON(tc.expectedProto())
			suite.Require().NoError(err)
			suite.Require().Equal(jsonMarshalledUpdatedBz, jsonMarshalExpectedResponse)
		})
	}
}
