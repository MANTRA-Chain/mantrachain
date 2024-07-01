package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/did module sentinel errors
var (
	ErrInvalidSigner              = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrDidDocumentNotFound        = errors.Register(ModuleName, 1101, "did document not found")
	ErrDidDocumentFound           = errors.Register(ModuleName, 1102, "did document found")
	ErrInvalidDIDFormat           = errors.Register(ModuleName, 1103, "input not compliant with the DID specifications (crf. https://www.w3.org/TR/did-core/#did-syntax)")
	ErrInvalidDIDURLFormat        = errors.Register(ModuleName, 1104, "input not compliant with the DID URL specifications (crf. https://www.w3.org/TR/did-core/#did-url-syntax)")
	ErrInvalidRFC3986UriFormat    = errors.Register(ModuleName, 1105, "input not compliant with the RFC3986 URI specifications (crf. https://datatracker.ietf.org/doc/html/rfc3986)")
	ErrEmptyRelationships         = errors.Register(ModuleName, 1106, "a verification method should have at least one verification relationship. (cfr. https://www.w3.org/TR/did-core/#verification-relationships)")
	ErrUnauthorized               = errors.Register(ModuleName, 1107, "the requester DID's verification method is not listed in the required relationship")
	ErrInvalidState               = errors.Register(ModuleName, 1108, "the requested action is not applicable on the resource")
	ErrInvalidInput               = errors.Register(ModuleName, 1109, "input is invalid")
	ErrVerificationMethodNotFound = errors.Register(ModuleName, 1110, "verification method not found")
	ErrInvalidDidMethodFormat     = errors.Register(ModuleName, 1111, "invalid did method format")
	ErrKeyFormatNotSupported      = errors.Register(ModuleName, 1112, "key format not supported")
)
