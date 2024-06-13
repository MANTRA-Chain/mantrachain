package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/did module sentinel errors
var (
	ErrDidDocumentNotFound        = sdkerrors.Register(ModuleName, 2, "did document not found")
	ErrDidDocumentFound           = sdkerrors.Register(ModuleName, 3, "did document found")
	ErrInvalidDIDFormat           = sdkerrors.Register(ModuleName, 4, "input not compliant with the DID specifications (crf. https://www.w3.org/TR/did-core/#did-syntax)")
	ErrInvalidDIDURLFormat        = sdkerrors.Register(ModuleName, 5, "input not compliant with the DID URL specifications (crf. https://www.w3.org/TR/did-core/#did-url-syntax)")
	ErrInvalidRFC3986UriFormat    = sdkerrors.Register(ModuleName, 6, "input not compliant with the RFC3986 URI specifications (crf. https://datatracker.ietf.org/doc/html/rfc3986)")
	ErrEmptyRelationships         = sdkerrors.Register(ModuleName, 7, "a verification method should have at least one verification relationship. (cfr. https://www.w3.org/TR/did-core/#verification-relationships)")
	ErrUnauthorized               = sdkerrors.Register(ModuleName, 8, "the requester DID's verification method is not listed in the required relationship")
	ErrInvalidState               = sdkerrors.Register(ModuleName, 9, "the requested action is not applicable on the resource")
	ErrInvalidInput               = sdkerrors.Register(ModuleName, 10, "input is invalid")
	ErrVerificationMethodNotFound = sdkerrors.Register(ModuleName, 11, "verification method not found")
	ErrInvalidDidMethodFormat     = sdkerrors.Register(ModuleName, 12, "invalid did method format")
	ErrKeyFormatNotSupported      = sdkerrors.Register(ModuleName, 13, "key format not supported")
)
