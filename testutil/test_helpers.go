package testutil

// DONTCOVER

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/ibc-go/v4/modules/core/exported"
	"github.com/cosmos/ibc-go/v4/testing/mock"
	ibctesting "github.com/cosmos/interchain-security/legacy_ibc_testing/testing"
	"github.com/cosmos/interchain-security/testutil/e2e"
	e2eutil "github.com/cosmos/interchain-security/testutil/e2e"
	icstestingutils "github.com/cosmos/interchain-security/testutil/ibc_testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/tendermint/spm/cosmoscmd"
	tmencoding "github.com/tendermint/tendermint/crypto/encoding"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	coinfactorytypes "github.com/MANTRA-Finance/mantrachain/x/coinfactory/types"
	guardtypes "github.com/MANTRA-Finance/mantrachain/x/guard/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	ibctmtypes "github.com/cosmos/ibc-go/v4/modules/light-clients/07-tendermint/types"
	consumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	ccv "github.com/cosmos/interchain-security/x/ccv/types"

	"github.com/MANTRA-Finance/mantrachain/app"
)

var (
	TestAdminAddress                          = "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw"
	TestAccountPrivilegesGuardNftCollectionId = "nft-guard-collection"
	SecondaryDenom                            = "ucoin"
	SecondaryAmount                           = sdk.NewInt(100000000)
)

func ConsumerAppIniter() (ibctesting.TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()
	encCdc := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)
	testApp := app.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})

	genesisState := app.NewDefaultGenesisState(encCdc.Marshaler)
	guardGenesis := guardtypes.NewGenesisState(guardtypes.NewParams(
		TestAdminAddress,
		TestAdminAddress,
		TestAccountPrivilegesGuardNftCollectionId,
		guardtypes.DefaultPrivileges,
	))
	genesisState[guardtypes.ModuleName] = testApp.AppCodec().MustMarshalJSON(guardGenesis)

	return testApp, genesisState
}

type IBCConnectionTestSuite struct {
	suite.Suite

	Coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	ChainProvider *ibctesting.TestChain
	Chain         *ibctesting.TestChain

	ProviderApp e2e.ProviderApp
	ChainApp    e2e.ConsumerApp

	CCVPath *ibctesting.Path

	TestAccs []sdk.AccAddress
}

func (suite *IBCConnectionTestSuite) GetApp(chain *ibctesting.TestChain) *app.App {
	testApp, ok := chain.App.(*app.App)
	if !ok {
		panic("not Consumer app")
	}

	return testApp
}

func (s *IBCConnectionTestSuite) SetupTest() {
	// Instantiate new coordinator and provider chain
	s.Coordinator = ibctesting.NewCoordinator(s.T(), 0)

	providerChain, providerApp := icstestingutils.AddProvider[e2eutil.ProviderApp](s.Coordinator, s.T(), icstestingutils.ProviderAppIniter)
	s.ChainProvider = providerChain
	s.ProviderApp = providerApp

	providerKeeper := providerApp.GetProviderKeeper()

	// re-assign all validator keys for the first consumer chain
	s.preProposalKeyAssignment()

	// start consumer chain
	bundle := icstestingutils.AddConsumer[e2eutil.ProviderApp, *app.App](s.Coordinator, &s.Suite, 0, ConsumerAppIniter)
	s.Chain = bundle.Chain
	s.ChainApp = bundle.Chain.App.(*app.App)

	genesisState, found := providerKeeper.GetConsumerGenesis(
		s.ChainProvider.GetContext(),
		s.Chain.ChainID,
	)
	s.Require().True(found, "consumer genesis not found")

	// initialize consumer chain with it's corresponding genesis state
	// stored on the provider.
	s.initConsumerChain(bundle, &genesisState)

	// try updating this consumer client on the provider chain
	err := bundle.Path.EndpointB.UpdateClient()
	s.Require().NoError(err)

	// try updating the provider client on this consumer chain
	err = bundle.Path.EndpointA.UpdateClient()
	s.Require().NoError(err)

	s.TestAccs = CreateRandomAccounts(3)
}

// preProposalKeyAssignment assigns keys to all provider validators for
// the consumer with chainID before the chain is registered, i.e.,
// before a client to the consumer is created
func (s *IBCConnectionTestSuite) preProposalKeyAssignment() {
	providerKeeper := s.ProviderApp.GetProviderKeeper()

	for _, val := range s.ChainProvider.Vals.Validators {
		// get SDK validator
		valAddr, err := sdk.ValAddressFromHex(val.Address.String())
		s.Require().NoError(err)
		validator, found := s.ProviderApp.GetE2eStakingKeeper().GetValidator(s.ChainProvider.GetContext(), valAddr)
		s.Require().True(found)

		// generate new PrivValidator
		privVal := mock.NewPV()
		tmPubKey, err := privVal.GetPubKey()
		s.Require().NoError(err)
		consumerKey, err := tmencoding.PubKeyToProto(tmPubKey)
		s.Require().NoError(err)

		// add Signer to the provider chain as there is no consumer chain to add it;
		// as a result, NewTestChainWithValSet in AddConsumer uses providerChain.Signers
		s.ChainProvider.Signers[tmPubKey.Address().String()] = privVal

		err = providerKeeper.AssignConsumerKey(s.ChainProvider.GetContext(), icstestingutils.FirstConsumerChainID, validator, consumerKey)
		s.Require().NoError(err)
	}
}

// GetConsumerEndpointClientAndConsState returns the client and consensus state
// for a particular consumer endpoint, as specified by the consumer's bundle.
func (s *IBCConnectionTestSuite) getConsumerEndpointClientAndConsState(consumerBundle icstestingutils.ConsumerBundle) (exported.ClientState, exported.ConsensusState) {

	ctx := consumerBundle.GetCtx()
	consumerKeeper := consumerBundle.GetKeeper()

	clientID, found := consumerKeeper.GetProviderClientID(ctx)
	s.Require().True(found)

	clientState, found := consumerBundle.App.GetIBCKeeper().ClientKeeper.GetClientState(ctx, clientID)
	s.Require().True(found)

	consState, found := consumerBundle.App.GetIBCKeeper().ClientKeeper.GetClientConsensusState(
		ctx, clientID, clientState.GetLatestHeight())
	s.Require().True(found)

	return clientState, consState
}

func (s *IBCConnectionTestSuite) validateEndpointsClientConfig(consumerBundle icstestingutils.ConsumerBundle) {
	consumerKeeper := consumerBundle.GetKeeper()
	providerStakingKeeper := s.ProviderApp.GetStakingKeeper()

	consumerUnbondingPeriod := consumerKeeper.GetUnbondingPeriod(consumerBundle.GetCtx())
	cs, ok := s.ProviderApp.GetIBCKeeper().ClientKeeper.GetClientState(s.ChainProvider.GetContext(),
		consumerBundle.Path.EndpointB.ClientID)
	s.Require().True(ok)
	s.Require().Equal(
		consumerUnbondingPeriod,
		cs.(*ibctmtypes.ClientState).UnbondingPeriod,
		"unexpected unbonding period in consumer client state",
	)

	providerUnbondingPeriod := providerStakingKeeper.UnbondingTime(s.ChainProvider.GetContext())
	cs, ok = consumerBundle.App.GetIBCKeeper().ClientKeeper.GetClientState(
		consumerBundle.GetCtx(), consumerBundle.Path.EndpointA.ClientID)
	s.Require().True(ok)
	s.Require().Equal(
		providerUnbondingPeriod,
		cs.(*ibctmtypes.ClientState).UnbondingPeriod,
		"unexpected unbonding period in provider client state",
	)
}

func (s *IBCConnectionTestSuite) initConsumerChain(
	bundle *icstestingutils.ConsumerBundle,
	genesisState *consumertypes.GenesisState,
) {
	providerKeeper := s.ProviderApp.GetProviderKeeper()

	// run CCV module init genesis
	s.NotPanics(func() {
		consumerKeeper := bundle.GetKeeper()
		consumerKeeper.InitGenesis(bundle.GetCtx(), genesisState)
	})

	// confirm client and cons state for consumer were set correctly in InitGenesis;
	// NOTE: on restart, both ProviderClientState and ProviderConsensusState are nil
	if genesisState.NewChain {
		consumerEndpointClientState,
			consumerEndpointConsState := s.getConsumerEndpointClientAndConsState(*bundle)
		s.Require().Equal(genesisState.ProviderClientState, consumerEndpointClientState)
		s.Require().Equal(genesisState.ProviderConsensusState, consumerEndpointConsState)
	}

	// create path for the CCV channel
	bundle.Path = ibctesting.NewPath(bundle.Chain, s.ChainProvider)

	// Set provider endpoint's clientID for each consumer
	providerEndpointClientID, found := providerKeeper.GetConsumerClientId(
		s.ChainProvider.GetContext(),
		bundle.Chain.ChainID,
	)
	s.Require().True(found, "provider endpoint clientID not found")
	bundle.Path.EndpointB.ClientID = providerEndpointClientID

	// Set consumer endpoint's clientID
	consumerKeeper := bundle.GetKeeper()
	consumerEndpointClientID, found := consumerKeeper.GetProviderClientID(bundle.GetCtx())
	s.Require().True(found, "consumer endpoint clientID not found")
	bundle.Path.EndpointA.ClientID = consumerEndpointClientID

	s.CCVPath = bundle.Path

	// Note: suite.path.EndpointA.ClientConfig and suite.path.EndpointB.ClientConfig are not populated,
	// since these IBC testing package fields are unused in our tests.

	// Confirm client config is now correct
	s.validateEndpointsClientConfig(*bundle)

	// - channel config
	bundle.Path.EndpointA.ChannelConfig.PortID = ccv.ConsumerPortID
	bundle.Path.EndpointB.ChannelConfig.PortID = ccv.ProviderPortID
	bundle.Path.EndpointA.ChannelConfig.Version = ccv.Version
	bundle.Path.EndpointB.ChannelConfig.Version = ccv.Version
	bundle.Path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	bundle.Path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED

	// create path for the transfer channel
	bundle.TransferPath = ibctesting.NewPath(bundle.Chain, s.ChainProvider)
	bundle.TransferPath.EndpointA.ChannelConfig.PortID = transfertypes.PortID
	bundle.TransferPath.EndpointB.ChannelConfig.PortID = transfertypes.PortID
	bundle.TransferPath.EndpointA.ChannelConfig.Version = transfertypes.Version
	bundle.TransferPath.EndpointB.ChannelConfig.Version = transfertypes.Version

	// commit state on this consumer chain
	s.Coordinator.CommitBlock(bundle.Chain)

	// try updating this consumer client on the provider chain
	err := bundle.Path.EndpointB.UpdateClient()
	s.Require().NoError(err)

	// try updating the provider client on this consumer chain
	err = bundle.Path.EndpointA.UpdateClient()
	s.Require().NoError(err)
}

func TestAddr(addr string, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	bechexpected := res.String()
	if bech != bechexpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bechres, res) {
		return nil, err
	}

	return res, nil
}

func CheckBalance(t *testing.T, app *app.App, addr sdk.AccAddress, balances sdk.Coins) {
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	require.True(t, balances.IsEqual(app.BankKeeper.GetAllBalances(ctxCheck, addr)))
}

func FundAccount(bankKeeper bankkeeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, coinfactorytypes.ModuleName, amounts); err != nil {
		return err
	}

	return bankKeeper.SendCoinsFromModuleToAccount(ctx, coinfactorytypes.ModuleName, addr, amounts)
}

func TestMessageAuthzSerialization(t *testing.T, msg sdk.Msg) {
	someDate := time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)
	const (
		mockGranter string = "cosmos1abc"
		mockGrantee string = "cosmos1xyz"
	)

	var (
		mockMsgGrant  authz.MsgGrant
		mockMsgRevoke authz.MsgRevoke
		mockMsgExec   authz.MsgExec
	)

	encodingConfig := app.MakeEncodingConfig()

	// Authz: Grant Msg
	typeURL := sdk.MsgTypeURL(msg)
	grant, err := authz.NewGrant(authz.NewGenericAuthorization(typeURL), someDate.Add(time.Hour))
	require.NoError(t, err)

	msgGrant := authz.MsgGrant{Granter: mockGranter, Grantee: mockGrantee, Grant: grant}
	msgGrantBytes := json.RawMessage(sdk.MustSortJSON(encodingConfig.Marshaler.MustMarshalJSON(&msgGrant)))
	err = encodingConfig.Marshaler.UnmarshalJSON(msgGrantBytes, &mockMsgGrant)
	require.NoError(t, err)

	// Authz: Revoke Msg
	msgRevoke := authz.MsgRevoke{Granter: mockGranter, Grantee: mockGrantee, MsgTypeUrl: typeURL}
	msgRevokeByte := json.RawMessage(sdk.MustSortJSON(encodingConfig.Marshaler.MustMarshalJSON(&msgRevoke)))
	err = encodingConfig.Marshaler.UnmarshalJSON(msgRevokeByte, &mockMsgRevoke)
	require.NoError(t, err)

	// Authz: Exec Msg
	msgAny, err := cdctypes.NewAnyWithValue(msg)
	require.NoError(t, err)
	msgExec := authz.MsgExec{Grantee: mockGrantee, Msgs: []*cdctypes.Any{msgAny}}
	execMsgByte := json.RawMessage(sdk.MustSortJSON(encodingConfig.Marshaler.MustMarshalJSON(&msgExec)))
	err = encodingConfig.Marshaler.UnmarshalJSON(execMsgByte, &mockMsgExec)
	require.NoError(t, err)
	require.Equal(t, msgExec.Msgs[0].Value, mockMsgExec.Msgs[0].Value)
}

// CreateRandomAccounts is a function return a list of randomly generated AccAddresses
func CreateRandomAccounts(numAccts int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, numAccts)
	for i := 0; i < numAccts; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}
