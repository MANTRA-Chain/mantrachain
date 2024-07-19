package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgCreateDidDocument_Route(t *testing.T) {
	assert.Equalf(t, ModuleName, MsgCreateDidDocument{}.Route(), "Route()")
}

func TestMsgCreateDidDocument_Type(t *testing.T) {
	assert.Equalf(t, TypeMsgCreateDidDocument, MsgCreateDidDocument{}.Type(), "Type()")
}

func TestMsgUpdateDidDocument_Route(t *testing.T) {
	assert.Equalf(t, ModuleName, MsgUpdateDidDocument{}.Route(), "Route()")
}

func TestMsgUpdateDidDocument_Type(t *testing.T) {
	assert.Equalf(t, TypeMsgUpdateDidDocument, MsgUpdateDidDocument{}.Type(), "Type()")
}

func TestMsgAddVerification_Route(t *testing.T) {
	assert.Equalf(t, ModuleName, MsgAddVerification{}.Route(), "Route()")
}

func TestMsgAddVerification_Type(t *testing.T) {
	assert.Equalf(t, TypeMsgAddVerification, MsgAddVerification{}.Type(), "Type()")
}

func TestMsgRevokeVerification_Route(t *testing.T) {
	assert.Equalf(t, ModuleName, MsgRevokeVerification{}.Route(), "Route()")
}

func TestMsgRevokeVerification_Type(t *testing.T) {
	assert.Equalf(t, TypeMsgRevokeVerification, MsgRevokeVerification{}.Type(), "Type()")
}

func TestMsgSetVerificationRelationships_Route(t *testing.T) {
	assert.Equalf(t, ModuleName, MsgSetVerificationRelationships{}.Route(), "Route()")
}

func TestMsgSetVerificationRelationships_Type(t *testing.T) {
	assert.Equalf(t, TypeMsgSetVerificationRelationships, MsgSetVerificationRelationships{}.Type(), "Type()")
}

func TestMsgDeleteService_Route(t *testing.T) {
	assert.Equalf(t, ModuleName, MsgDeleteService{}.Route(), "Route()")
}

func TestMsgDeleteService_Type(t *testing.T) {
	assert.Equalf(t, TypeMsgDeleteService, MsgDeleteService{}.Type(), "Type()")
}

func TestMsgAddService_Route(t *testing.T) {
	assert.Equalf(t, ModuleName, MsgAddService{}.Route(), "Route()")
}

func TestMsgAddService_Type(t *testing.T) {
	assert.Equalf(t, TypeMsgAddService, MsgAddService{}.Type(), "Type()")
}

func TestMsgAddController_Route(t *testing.T) {
	assert.Equalf(t, ModuleName, MsgAddController{}.Route(), "Route()")
}

func TestMsgAddController_Type(t *testing.T) {
	assert.Equalf(t, TypeMsgAddController, MsgAddController{}.Type(), "Type()")
}

func TestMsgDeleteController_Route(t *testing.T) {
	assert.Equalf(t, ModuleName, MsgDeleteController{}.Route(), "Route()")
}

func TestMsgDeleteController_Type(t *testing.T) {
	assert.Equalf(t, TypeMsgDeleteController, MsgDeleteController{}.Type(), "Type()")
}
