package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

// msg types
const (
	TypeMsgCreateDidDocument = "create-did"
)

var (
	_ legacytx.LegacyMsg = &MsgCreateDidDocument{}
	_ sdk.Msg            = &MsgCreateDidDocument{}
)

// NewMsgCreateDidDocument creates a new MsgCreateDidDocument instance
func NewMsgCreateDidDocument(
	id string,
	verifications []*Verification,
	services []*Service,
	signerAccount string,
) *MsgCreateDidDocument {
	return &MsgCreateDidDocument{
		Id:            id,
		Verifications: verifications,
		Services:      services,
		Signer:        signerAccount,
	}
}

// Route implements sdk.Msg
func (MsgCreateDidDocument) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgCreateDidDocument) Type() string {
	return TypeMsgCreateDidDocument
}

func (msg *MsgCreateDidDocument) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// --------------------------
// UPDATE IDENTIFIER
// --------------------------

// msg types
const (
	TypeMsgUpdateDidDocument = "update-did"
)

var (
	_ legacytx.LegacyMsg = &MsgUpdateDidDocument{}
	_ sdk.Msg            = &MsgUpdateDidDocument{}
)

func NewMsgUpdateDidDocument(
	didDoc *DidDocument,
	signerAccount string,
) *MsgUpdateDidDocument {
	return &MsgUpdateDidDocument{
		Doc:    didDoc,
		Signer: signerAccount,
	}
}

// Route implements sdk.Msg
func (MsgUpdateDidDocument) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgUpdateDidDocument) Type() string {
	return TypeMsgUpdateDidDocument
}

func (msg *MsgUpdateDidDocument) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// --------------------------
// ADD VERIFICATION
// --------------------------
// msg types
const (
	TypeMsgAddVerification = "add-verification"
)

var (
	_ legacytx.LegacyMsg = &MsgAddVerification{}
	_ sdk.Msg            = &MsgAddVerification{}
)

// NewMsgAddVerification creates a new MsgAddVerification instance
func NewMsgAddVerification(
	id string,
	verification *Verification,
	signerAccount string,
) *MsgAddVerification {
	return &MsgAddVerification{
		Id:           id,
		Verification: verification,
		Signer:       signerAccount,
	}
}

// Route implements sdk.Msg
func (MsgAddVerification) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgAddVerification) Type() string {
	return TypeMsgAddVerification
}

func (msg *MsgAddVerification) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// --------------------------
// REVOKE VERIFICATION
// --------------------------

// msg types
const (
	TypeMsgRevokeVerification = "revoke-verification"
)

var (
	_ legacytx.LegacyMsg = &MsgRevokeVerification{}
	_ sdk.Msg            = &MsgRevokeVerification{}
)

// NewMsgRevokeVerification creates a new MsgRevokeVerification instance
func NewMsgRevokeVerification(
	id string,
	methodID string,
	signerAccount string,
) *MsgRevokeVerification {
	return &MsgRevokeVerification{
		Id:       id,
		MethodId: methodID,
		Signer:   signerAccount,
	}
}

// Route implements sdk.Msg
func (MsgRevokeVerification) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgRevokeVerification) Type() string {
	return TypeMsgRevokeVerification
}

func (msg *MsgRevokeVerification) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// --------------------------
// SET VERIFICATION RELATIONSHIPS
// --------------------------
// msg types
const (
	TypeMsgSetVerificationRelationships = "set-verification-relationships"
)

var (
	_ legacytx.LegacyMsg = &MsgSetVerificationRelationships{}
	_ sdk.Msg            = &MsgSetVerificationRelationships{}
)

func NewMsgSetVerificationRelationships(
	id string,
	methodID string,
	relationships []string,
	signerAccount string,
) *MsgSetVerificationRelationships {
	return &MsgSetVerificationRelationships{
		Id:            id,
		MethodId:      methodID,
		Relationships: relationships,
		Signer:        signerAccount,
	}
}

// Route implements sdk.Msg
func (MsgSetVerificationRelationships) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgSetVerificationRelationships) Type() string {
	return TypeMsgSetVerificationRelationships
}

func (msg *MsgSetVerificationRelationships) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// --------------------------
// ADD SERVICE
// --------------------------

// msg types
const (
	TypeMsgAddService = "add-service"
)

var (
	_ legacytx.LegacyMsg = &MsgAddService{}
	_ sdk.Msg            = &MsgAddService{}
)

// NewMsgAddService creates a new MsgAddService instance
func NewMsgAddService(
	id string,
	service *Service,
	signerAccount string,
) *MsgAddService {
	return &MsgAddService{
		Id:          id,
		ServiceData: service,
		Signer:      signerAccount,
	}
}

// Route implements sdk.Msg
func (MsgAddService) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgAddService) Type() string {
	return TypeMsgAddService
}

func (msg *MsgAddService) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// --------------------------
// DELETE SERVICE
// --------------------------

// msg types
const (
	TypeMsgDeleteService = "delete-service"
)

var (
	_ legacytx.LegacyMsg = &MsgDeleteService{}
	_ sdk.Msg            = &MsgDeleteService{}
)

func NewMsgDeleteService(
	id string,
	serviceID string,
	signerAccount string,
) *MsgDeleteService {
	return &MsgDeleteService{
		Id:        id,
		ServiceId: serviceID,
		Signer:    signerAccount,
	}
}

// Route implements sdk.Msg
func (MsgDeleteService) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgDeleteService) Type() string {
	return TypeMsgDeleteService
}

func (msg *MsgDeleteService) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// --------------------------
// ADD CONTROLLER
// --------------------------

// msg types
const (
	TypeMsgAddController = "add-controller"
)

var (
	_ legacytx.LegacyMsg = &MsgAddController{}
	_ sdk.Msg            = &MsgAddController{}
)

func NewMsgAddController(
	id string,
	controllerDID string,
	signerAccount string,
) *MsgAddController {
	return &MsgAddController{
		Id:            id,
		ControllerDid: controllerDID,
		Signer:        signerAccount,
	}
}

// Route implements sdk.Msg
func (MsgAddController) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgAddController) Type() string {
	return TypeMsgAddController
}

func (msg *MsgAddController) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// --------------------------
// DELETE CONTROLLER
// --------------------------

// msg types
const (
	TypeMsgDeleteController = "delete-controller"
)

var (
	_ legacytx.LegacyMsg = &MsgDeleteController{}
	_ sdk.Msg            = &MsgDeleteController{}
)

func NewMsgDeleteController(
	id string,
	controllerDID string,
	signerAccount string,
) *MsgDeleteController {
	return &MsgDeleteController{
		Id:            id,
		ControllerDid: controllerDID,
		Signer:        signerAccount,
	}
}

// Route implements sdk.Msg
func (MsgDeleteController) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgDeleteController) Type() string {
	return TypeMsgDeleteController
}

func (msg *MsgDeleteController) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
