package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/errors"
	"github.com/MANTRA-Finance/mantrachain/x/did/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestHandleMsgCreateDidDocument() {
	var (
		req    types.MsgCreateDidDocument
		errExp error
	)

	server := NewMsgServerImpl(suite.keeper)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"Pass: can create a an did",
			func() {
				req = *types.NewMsgCreateDidDocument("did:cosmos:mantrachain:subject", nil, nil, "subject")
				errExp = nil
			},
		},
		{
			"FAIL: did doc validation fails",
			func() {
				req = *types.NewMsgCreateDidDocument("invalid did", nil, nil, "subject")
				errExp = errors.Wrapf(types.ErrInvalidDIDFormat, "did %s", "invalid did")
			},
		},
		{
			"FAIL: did already exists",
			func() {
				did := "did:cosmos:mantrachain:subject"
				didDoc, _ := types.NewDidDocument(did)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				req = *types.NewMsgCreateDidDocument(did, nil, nil, "subject")
				errExp = errors.Wrapf(types.ErrDidDocumentFound, "a document with did %s already exists", did)
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			tc.malleate()
			_, err := server.CreateDidDocument(sdk.WrapSDKContext(suite.ctx), &req)
			if errExp == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(errExp.Error(), err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleMsgUpdateDidDocument() {
	var (
		req    types.MsgUpdateDidDocument
		errExp error
	)

	server := NewMsgServerImpl(suite.keeper)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"FAIL: not found",
			func() {
				req = *types.NewMsgUpdateDidDocument(&types.DidDocument{Id: "did:cosmos:mantrachain:subject"}, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrDidDocumentNotFound, "did document at %s not found", "did:cosmos:mantrachain:subject")
			},
		},
		{
			"FAIL: unauthorized",
			func() {

				did := "did:cosmos:mantrachain:subject"
				didDoc, _ := types.NewDidDocument(did)
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)

				req = *types.NewMsgUpdateDidDocument(&types.DidDocument{Id: didDoc.Id, Controller: []string{"did:cosmos:mantrachain:controller"}}, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrUnauthorized, "signer account %s not authorized to update the target did document at %s", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8", did)

			},
		},
		{
			"PASS: replace did document",
			func() {

				did := "did:cosmos:mantrachain:subject"
				didDoc, _ := types.NewDidDocument(did, types.WithVerifications(
					types.NewVerification(
						types.NewVerificationMethod(
							"did:cosmos:mantrachain:subject#key-1",
							"did:cosmos:mantrachain:subject",
							types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{types.Authentication},
						nil,
					),
				))
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))

				newDidDoc, err := types.NewDidDocument(did)
				suite.Require().Nil(err)

				req = *types.NewMsgUpdateDidDocument(&newDidDoc, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = nil
			},
		},
		{
			"FAIL: invalid controllers",
			func() {
				didDoc, _ := types.NewDidDocument("did:cosmos:mantrachain:subject", types.WithVerifications(
					types.NewVerification(
						types.NewVerificationMethod(
							"did:cosmos:mantrachain:subject#key-1",
							"did:cosmos:mantrachain:subject",
							types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{types.Authentication},
						nil,
					),
				))
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)

				controllers := []string{
					"did:cosmos:mantrachain:controller-1",
					"did:cosmos:mantrachain:controller-2",
					"invalid",
				}

				req = *types.NewMsgUpdateDidDocument(&types.DidDocument{Id: didDoc.Id, Controller: controllers}, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrInvalidDIDFormat, "invalid did document")
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			tc.malleate()

			_, err := server.UpdateDidDocument(sdk.WrapSDKContext(suite.ctx), &req)

			if errExp == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(errExp.Error(), err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleMsgAddVerification() {
	var (
		req    types.MsgAddVerification
		errExp error
	)

	server := NewMsgServerImpl(suite.keeper)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"FAIL: can not add verification, did does not exist",
			func() {
				req = *types.NewMsgAddVerification("did:cosmos:mantrachain:subject", nil, "subject")
				errExp = errors.Wrapf(types.ErrDidDocumentNotFound, "did document at %s not found", "did:cosmos:mantrachain:subject")
			},
		},
		{
			"FAIL: can not add verification, unauthorized",
			func() {
				// setup
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.CapabilityInvocation},
							nil,
						),
					),
				)
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				// actual test
				v := types.NewVerification(
					types.NewVerificationMethod(
						"did:cosmos:mantrachain:subject#key-2",
						"did:cosmos:mantrachain:subject",
						types.NewBlockchainAccountID("foochainid", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
					),
					[]string{types.Authentication},
					nil,
				)
				req = *types.NewMsgAddVerification(didDoc.Id, v, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrUnauthorized, "signer account %s not authorized to update the target did document at %s", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8", didDoc.Id)
			},
		},
		{
			"FAIL: can not add verification, unauthorized, key mismatch",
			func() {
				// setup
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				// actual test
				v := types.NewVerification(
					types.NewVerificationMethod(
						"did:cosmos:mantrachain:subject#key-2",
						"did:cosmos:mantrachain:subject",
						types.NewBlockchainAccountID("foochainid", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
					),
					[]string{types.Authentication},
					nil,
				)
				req = *types.NewMsgAddVerification(didDoc.Id, v, "cash1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2")
				errExp = errors.Wrapf(types.ErrUnauthorized, "signer account %s not authorized to update the target did document at %s", "cash1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2", didDoc.Id)
			},
		},
		{
			"FAIL: can not add verification, invalid verification",
			func() {
				// setup
				//signer := "subject"
				signer := "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				// actual test
				v := types.NewVerification(
					types.NewVerificationMethod(
						"",
						"did:cosmos:mantrachain:subject",
						types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
					),
					[]string{types.Authentication},
					nil,
				)
				req = *types.NewMsgAddVerification(didDoc.Id, v, signer)
				errExp = errors.Wrapf(types.ErrInvalidDIDURLFormat, "verification method id: %v", "")
			},
		},
		{
			"PASS: can add verification to did document",
			func() {
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)

				v := types.NewVerification(
					types.NewVerificationMethod(
						"did:cosmos:mantrachain:subject#key-2",
						"did:cosmos:mantrachain:subject",
						types.NewBlockchainAccountID("foochainid", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
					),
					[]string{types.Authentication},
					nil,
				)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))
				req = *types.NewMsgAddVerification(didDoc.Id, v, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = nil
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			tc.malleate()

			_, err := server.AddVerification(sdk.WrapSDKContext(suite.ctx), &req)

			if errExp == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(errExp.Error(), err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleMsgSetVerificationRelationships() {
	var (
		req    types.MsgSetVerificationRelationships
		errExp error
	)

	server := NewMsgServerImpl(suite.keeper)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"FAIL: can not add verification relationship, did does not exist",
			func() {
				req = *types.NewMsgSetVerificationRelationships(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:mantrachain:subject#key-1",
					[]string{types.Authentication},
					"subject",
				)
				errExp = errors.Wrapf(types.ErrDidDocumentNotFound, "did document at %s not found", "did:cosmos:mantrachain:subject")
			},
		},
		{
			"FAIL: can not add verification relationship, unauthorized",
			func() {
				// setup
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.CapabilityInvocation},
							nil,
						),
					),
				)
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				// actual test
				req = *types.NewMsgSetVerificationRelationships(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:mantrachain:subject#key-1",
					[]string{types.Authentication},
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
				)

				errExp = errors.Wrapf(types.ErrUnauthorized, "signer account %s not authorized to update the target did document at %s", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8", "did:cosmos:mantrachain:subject")
			},
		},
		{
			"FAIL: can not add verification relationship, invalid relationship provided",
			func() {
				// setup
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								//"did:cosmos:mantrachain:subject#cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				// actual test
				req = *types.NewMsgSetVerificationRelationships(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:mantrachain:subject#key-1",
					nil,
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
				)
				errExp = errors.Wrap(types.ErrEmptyRelationships, "at least a verification relationship is required")
			},
		},
		{
			"FAIL: verification method does not exist ",
			func() {
				// setup
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				// actual test
				req = *types.NewMsgSetVerificationRelationships(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:mantrachain:subject#key-does-not-exists",
					[]string{types.Authentication, types.CapabilityInvocation},
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
				)
				errExp = errors.Wrapf(types.ErrVerificationMethodNotFound, "verification method %v not found", "did:cosmos:mantrachain:subject#key-does-not-exists")
			},
		},
		{
			"PASS: add a new relationship",
			func() {
				// setup
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))
				// actual test
				req = *types.NewMsgSetVerificationRelationships(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:mantrachain:subject#key-1",
					[]string{types.Authentication, types.CapabilityInvocation},
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
				)
				errExp = nil
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			tc.malleate()

			_, err := server.SetVerificationRelationships(sdk.WrapSDKContext(suite.ctx), &req)

			if errExp == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(errExp.Error(), err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleMsgRevokeVerification() {
	var (
		req    types.MsgRevokeVerification
		errExp error
	)

	server := NewMsgServerImpl(suite.keeper)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"FAIL: can not revoke verification, did does not exist",
			func() {
				req = *types.NewMsgRevokeVerification("did:cosmos:mantrachain:2222", "service-id", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrDidDocumentNotFound, "did document at %s not found", "did:cosmos:mantrachain:2222")
			},
		},
		{
			"FAIL: can not revoke verification, not found",
			func() {
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				req = *types.NewMsgRevokeVerification(didDoc.Id, "did:cosmos:mantrachain:subject#not-existent", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrVerificationMethodNotFound, "verification method id: %v", "did:cosmos:mantrachain:subject#not-existent")
			},
		},
		{
			"FAIL: can not revoke verification, unauthorized",
			func() {
				signer := "controller-1"
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.CapabilityDelegation},
							nil,
						),
					),
				)

				vmID := "did:cosmos:mantrachain:subject#key-1"

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				// controller-1 does not exists
				req = *types.NewMsgRevokeVerification(didDoc.Id, vmID, signer)

				errExp = errors.Wrapf(types.ErrUnauthorized, "signer account %s not authorized to update the target did document at %s", signer, didDoc.Id)
			},
		},
		{
			"PASS: can revoke verification",
			func() {
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))
				req = *types.NewMsgRevokeVerification(didDoc.Id,
					"did:cosmos:mantrachain:subject#key-1",
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
				)
				errExp = nil
			},
		},
	}
	for i, tc := range testCases {
		suite.Run(fmt.Sprintf("TestHandleMsgRevokeVerification#%v", i), func() {
			tc.malleate()

			_, err := server.RevokeVerification(sdk.WrapSDKContext(suite.ctx), &req)

			if errExp == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(errExp.Error(), err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleMsgAddService() {
	var (
		req    types.MsgAddService
		errExp error
	)

	server := NewMsgServerImpl(suite.keeper)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"FAIL: can not add service, did does not exist",
			func() {
				service := types.NewService(
					"service-id",
					"NonUserCredential",
					"cash/multihash",
				)
				req = *types.NewMsgAddService("did:cosmos:mantrachain:subject", service, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrDidDocumentNotFound, "did document at %s not found", "did:cosmos:mantrachain:subject")
			},
		},
		{
			"FAIL: can not add service, service not defined",
			func() {
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)
				// create the did
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				// try adding a service
				req = *types.NewMsgAddService("did:cosmos:mantrachain:subject", nil, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrap(types.ErrInvalidInput, "service is not defined")
			},
		},
		{
			"FAIL: cannot add service to did document (unauthorized, wrong relationship)",
			func() {
				signer := "subject"
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.CapabilityInvocation, types.CapabilityDelegation},
							nil,
						),
					),
				)

				service := types.NewService(
					"service-id",
					"UserCredential",
					"cash/multihash",
				)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				req = *types.NewMsgAddService(didDoc.Id, service, signer)

				errExp = errors.Wrapf(types.ErrUnauthorized, "signer account %s not authorized to update the target did document at %s", signer, didDoc.Id)
			},
		},
		{
			"FAIL: cannot add service to did document with an empty type",
			func() {
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)

				service := types.NewService(
					"service-id",
					"",
					"cash/multihash",
				)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				req = *types.NewMsgAddService(didDoc.Id, service, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrap(types.ErrInvalidInput, "service type cannot be empty")
			},
		},
		{
			"FAIL: duplicated service",
			func() {
				//signer := "subject"
				signer := "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
					types.WithServices(
						types.NewService(
							"service-id",
							"UserCredential",
							"cash/multihash",
						),
					),
				)

				service := types.NewService(
					"service-id",
					"UserCredential",
					"cash/multihash",
				)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				req = *types.NewMsgAddService(didDoc.Id, service, signer)
				errExp = errors.Wrapf(types.ErrInvalidInput, "duplicated verification method id %s", "service-id")
			},
		},
		{
			"PASS: can add service to did document",
			func() {
				signer := "subject"
				didDoc, err := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewBlockchainAccountID("foochainid", signer),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)

				if err != nil {
					suite.FailNow("test setup failed: ", err)
				}

				service := types.NewService(
					"service-id",
					"UserCredential",
					"cash/multihash",
				)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))

				req = *types.NewMsgAddService(didDoc.Id, service, signer)
				errExp = nil
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()

			_, err := server.AddService(sdk.WrapSDKContext(suite.ctx), &req)

			if errExp == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(errExp.Error(), err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleMsgDeleteService() {
	var (
		req    types.MsgDeleteService
		errExp error
	)

	server := NewMsgServerImpl(suite.keeper)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"FAIL: can not delete service, did does not exist",
			func() {
				req = *types.NewMsgDeleteService("did:cosmos:mantrachain:2222", "service-id", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrDidDocumentNotFound, "did document at %s not found", "did:cosmos:mantrachain:2222")
			},
		},
		{

			"PASS: can delete service from did document",
			func() {
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
					types.WithServices(
						types.NewService(
							"service-id",
							"UserCredential",
							"cash/multihash",
						),
					),
				)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))
				req = *types.NewMsgDeleteService(didDoc.Id, "service-id", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = nil
			},
		},
		{
			"FAIL: cannot remove an invalid serviceID",
			func() {

				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)

				serviceID := ""

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				req = *types.NewMsgDeleteService(didDoc.Id, serviceID, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrInvalidState, "the did document doesn't have services associated")
			},
		},
		{
			"FAIL: unauthorized (wrong relationship)",
			func() {
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								"did:cosmos:mantrachain:subject",
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.CapabilityInvocation},
							nil,
						),
					),
				)

				serviceID := "service-id"

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				req = *types.NewMsgDeleteService(didDoc.Id, serviceID, "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrUnauthorized, "signer account %s not authorized to update the target did document at %s", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8", didDoc.Id)
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf(tc.name), func() {
			tc.malleate()

			_, err := server.DeleteService(sdk.WrapSDKContext(suite.ctx), &req)

			if errExp == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(errExp.Error(), err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleMsgAddController() {
	var (
		req    types.MsgAddController
		errExp error
	)

	server := NewMsgServerImpl(suite.keeper)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"FAIL: cannot add controller, did doesn't exist",
			func() {
				req = *types.NewMsgAddController(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrDidDocumentNotFound, "did document at %s not found", "did:cosmos:mantrachain:subject")
			},
		},
		{
			"FAIL: controller is not a valid did",
			func() {
				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)
				// create the did
				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				// try adding a service
				req = *types.NewMsgAddController("did:cosmos:mantrachain:subject", "", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrap(types.ErrInvalidDIDFormat, "did document controller validation error ''")
			},
		},
		{
			"FAIL: signer not authorized to change controller",
			func() {

				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.CapabilityInvocation, types.CapabilityDelegation},
							nil,
						),
					),
				)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)

				req = *types.NewMsgAddController(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
					"cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2", // does not match the pub key (it's the new controller)
				)

				errExp = errors.Wrapf(types.ErrUnauthorized, "signer account %s not authorized to update the target did document at %s", "cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2", didDoc.Id)
			},
		},
		{
			"FAIL: controller is not type key",
			func() {

				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)

				req = *types.NewMsgAddController(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:net:foochain:whatever",
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8", // does not match the pub key (it's the new controller)
				)

				errExp = errors.Wrapf(types.ErrInvalidInput, "did document controller 'did:cosmos:net:foochain:whatever' must be of type key")
			},
		},
		{
			"PASS: can add controller (via authentication relationship)",
			func() {
				didDoc, err := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)

				if err != nil {
					suite.FailNow("test setup failed: ", err)
				}

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))

				req = *types.NewMsgAddController(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = nil
			},
		},
		{
			"PASS: can add controller (via controller)",
			func() {
				didDoc, err := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.AssertionMethod},
							nil,
						),
					),
					types.WithControllers("did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2"),
				)

				if err != nil {
					suite.FailNow("test setup failed: ", err)
				}

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))

				req = *types.NewMsgAddController(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:key:cosmos17t8t3t6a6vpgk69perfyq930593sa8dn4kzsdf",
					"cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2") // this is the controller

				errExp = nil
			},
		},
		{
			"PASS: controller already added (duplicated)",
			func() {
				didDoc, err := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)

				if err != nil {
					suite.FailNow("test setup failed: ", err)
				}

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))

				req = *types.NewMsgAddController(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = nil
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()

			_, err := server.AddController(sdk.WrapSDKContext(suite.ctx), &req)

			if errExp == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(errExp.Error(), err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHandleMsgDeleteController() {
	var (
		req    types.MsgDeleteController
		errExp error
	)

	server := NewMsgServerImpl(suite.keeper)

	// FAIL: cannot delete controller, did doesn't exist
	// FAIL: signer not authorized to change controller
	// PASS: controller removed
	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"FAIL: cannot delete controller, did doesn't exist",
			func() {
				req = *types.NewMsgDeleteController(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = errors.Wrapf(types.ErrDidDocumentNotFound, "did document at %s not found", "did:cosmos:mantrachain:subject")
			},
		},
		{
			"FAIL: signer not authorized to change controller",
			func() {

				didDoc, _ := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.CapabilityInvocation, types.CapabilityDelegation},
							nil,
						),
					),
				)

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)

				req = *types.NewMsgDeleteController(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
					"cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2", // does not match the pub key (it's the new controller)
				)

				errExp = errors.Wrapf(types.ErrUnauthorized, "signer account %s not authorized to update the target did document at %s", "cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2", didDoc.Id)
			},
		},
		{
			"PASS: can delete controller (via authentication relationship)",
			func() {
				didDoc, err := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.Authentication},
							nil,
						),
					),
				)

				if err != nil {
					suite.FailNow("test setup failed: ", err)
				}

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))

				req = *types.NewMsgDeleteController(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
					"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8")
				errExp = nil
			},
		}, {
			"PASS: can delete controller (via controller)",
			func() {
				didDoc, err := types.NewDidDocument(
					"did:cosmos:mantrachain:subject",
					types.WithVerifications(
						types.NewVerification(
							types.NewVerificationMethod(
								"did:cosmos:mantrachain:subject#key-1",
								types.DID("did:cosmos:mantrachain:subject"),
								types.NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{types.AssertionMethod},
							nil,
						),
					),
					types.WithControllers(
						"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
						"did:cosmos:key:cosmos17t8t3t6a6vpgk69perfyq930593sa8dn4kzsdf",
					),
				)

				if err != nil {
					suite.FailNow("test setup failed: ", err)
				}

				suite.keeper.SetDidDocument(suite.ctx, []byte(didDoc.Id), didDoc)
				suite.keeper.SetDidMetadata(suite.ctx, []byte(didDoc.Id), types.NewDidMetadata([]byte{1}, time.Now()))

				req = *types.NewMsgDeleteController(
					"did:cosmos:mantrachain:subject",
					"did:cosmos:key:cosmos17t8t3t6a6vpgk69perfyq930593sa8dn4kzsdf",
					"cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2", // this is the controller
				)

				errExp = nil
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()

			_, err := server.DeleteController(sdk.WrapSDKContext(suite.ctx), &req)

			if errExp == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Equal(errExp.Error(), err.Error())
			}
		})
	}
}
