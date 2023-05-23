package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChainDID(t *testing.T) {

	tests := []struct {
		did   string
		chain string
		want  DID
	}{
		{
			"subject",
			"mantrachain",
			DID("did:cosmos:net:mantrachain:subject"),
		},
		{
			"",
			"mantrachain",
			DID("did:cosmos:net:mantrachain:"),
		},
		{
			"cosmos1uam3kpjdx3wksx46lzq6y628wwyzv0guuren75",
			"cosmoshub",
			DID("did:cosmos:net:cosmoshub:cosmos1uam3kpjdx3wksx46lzq6y628wwyzv0guuren75"),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestDID#", i), func(t *testing.T) {
			if got := NewChainDID(tt.chain, tt.did); got != tt.want {
				t.Errorf("DID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKeyDID(t *testing.T) {

	tests := []struct {
		name    string
		account string
		want    DID
	}{
		{
			"PASS: valid account",
			"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
			"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewKeyDID(tt.account), "NewKeyDID(%v)", tt.account)
		})
	}
}

func TestDID_NewVerificationMethodID(t *testing.T) {

	tests := []struct {
		name string
		did  DID
		vmID string
		want string
	}{
		{
			"PASS: generated vmId for network DID",
			DID("did:cosmos:net:foochain:whatever"),
			"123456",
			"did:cosmos:net:foochain:whatever#123456",
		},
		{
			"PASS: generated vmId for key DID",
			DID("did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
			"123456",
			"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8#123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.did.NewVerificationMethodID(tt.vmID), "NewVerificationMethodID(%v)", tt.vmID)
		})
	}
}

func TestIsValidDID(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"did:mantrachain:net:subject", true},
		{"did:mantrachain:cosmos1uam3kpjdx3wksx46lzq6y628wwyzv0guuren75", true},
		{"did:mantrachain:cosmos1uam3kpjdx3wksx46lzq6y628wwyzv0guuren75#key-1", false},
		{"did:subject", false},
		{"DID:cash:subject", false},
		{"d1d:cash:subject", false},
		{"d1d:#:subject", false},
		{"", false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestIsValidDID#", i), func(t *testing.T) {
			if got := IsValidDID(tt.input); got != tt.want {
				t.Errorf("IsValidDID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidDIDURL(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"did:mantrachain:subject", true},
		{"did:mantrachain:cosmos1uam3kpjdx3wksx46lzq6y628wwyzv0guuren75", true},
		{"did:mantrachain:cosmos1uam3kpjdx3wksx46lzq6y628wwyzv0guuren75#key-1", true},
		{"did:mantrachain:cosmos1uam3kpjdx3wksx46lzq6y628wwyzv0guuren75?key=1", true},
		{"did:cosmos:net:cosmoscash-testnet:575d062c-d110-42a9-9c04-cb1ff8c01f06#Z46DAL1MrJlVW_WmJ19WY8AeIpGeFOWl49Qwhvsnn2M", true},
		{"did:subject", false},
		{"DID:cash:subject", false},
		{"d1d:cash:subject", false},
		{"d1d:#:subject", false},
		{"", false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestIsValidDIDURL#", i), func(t *testing.T) {
			if got := IsValidDIDURL(tt.input); got != tt.want {
				t.Errorf("IsValidDIDURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidRFC3986Uri(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			"[][àséf",
			true,
		},
		{
			"# \u007e // / / ///// //// // / / ??? ?? 不",
			true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestIsValidRFC3986Uri#", i), func(t *testing.T) {
			if got := IsValidRFC3986Uri(tt.input); got != tt.want {
				t.Errorf("IsValidRFC3986Uri() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidDIDDocument(t *testing.T) {
	tests := []struct {
		name  string
		didFn func() *DidDocument
		want  bool
	}{
		{
			"PASS: document is valid",
			func() *DidDocument {
				return &DidDocument{
					Context: []string{contextDIDBase},
					Id:      "did:cosmos:net:mantrachain:1",
				}
			},
			true, // all good
		},
		{
			"FAIL: empty context",
			func() *DidDocument {
				return &DidDocument{
					Context: []string{},
					Id:      "did:cosmos:net:mantrachain:1",
				}
			},
			false, // missing context
		},
		{
			"PASS: minimal did document",
			func() *DidDocument {
				dd, err := NewDidDocument("did:cosmos:mantrachain:1")
				assert.NoError(t, err)
				return &dd
			},
			true, // all good
		},
		{
			"FAIL: empty did",
			func() *DidDocument {
				return &DidDocument{
					Context: []string{contextDIDBase},
					Id:      "",
				}
			},
			false, // empty id
		},
		{
			"FAIL: nil did document",
			func() *DidDocument {
				return nil
			},
			false, // nil pointer
		},
		{
			"PASS: did with valid controller",
			func() *DidDocument {
				dd, err := NewDidDocument("did:cosmos:key:cas:1", WithControllers(
					"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
				))
				assert.NoError(t, err)
				return &dd
			},
			true,
		},
		{
			"FAIL: invalid controller",
			func() *DidDocument {
				return &DidDocument{
					Context: []string{contextDIDBase},
					Id:      "did:cosmos:net:foochain:1",
					Controller: []string{
						"did:cosmos:net:foochain:whatever",
					},
				}
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint("TestIsValidDIDDocument#", tt.name), func(t *testing.T) {
			if got := IsValidDIDDocument(tt.didFn()); got != tt.want {
				t.Errorf("TestIsValidDIDDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidDIDMetadata(t *testing.T) {

	tests := []struct {
		didMetaFn func() *DidMetadata
		want      bool
	}{
		{
			func() *DidMetadata {
				now := time.Now()
				return &DidMetadata{
					VersionId: "d95daac05a36f93d1494208d02d1522d758466c62ea6b64c50b78999d2021f51",
					Created:   &now,
				}
			},
			true, // all good
		},
		{
			func() *DidMetadata {
				now := time.Now()
				return &DidMetadata{
					VersionId: "",
					Created:   &now,
				}
			},
			false, // missing version
		},
		{
			func() *DidMetadata {
				now := time.Now()
				return &DidMetadata{
					VersionId: "d95daac05a36f93d1494208d02d1522d758466c62ea6b64c50b78999d2021f51",
					Updated:   &now,
				}
			},
			false, // null created
		},
		{
			func() *DidMetadata {
				var now time.Time
				return &DidMetadata{
					VersionId: "d95daac05a36f93d1494208d02d1522d758466c62ea6b64c50b78999d2021f51",
					Created:   &now,
				}
			},
			false, // zero created
		},
		{
			func() *DidMetadata {
				return nil
			},
			false, // nil pointer
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestIsValidDIDMetadata#", i), func(t *testing.T) {
			if got := IsValidDIDMetadata(tt.didMetaFn()); got != tt.want {
				t.Errorf("TestIsValidDIDMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateVerification(t *testing.T) {
	tests := []struct {
		v       *Verification
		wantErr bool
	}{
		{
			v: NewVerification(
				NewVerificationMethod(
					"did:mantrachain:subject#key-1",
					"did:mantrachain:subject",
					NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
				),
				nil,
				nil,
			),
			wantErr: true, // relationships are nil
		},
		{
			v:       nil,
			wantErr: true,
		},
		{
			v: NewVerification(
				NewVerificationMethod(
					"did:mantrachain:subject#key-1",
					DID("did:mantrachain:subject"),
					NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
				),
				[]string{string(AssertionMethod)},
				nil,
			),
			wantErr: false,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestValidateVerification#", i), func(t *testing.T) {
			if err := ValidateVerification(tt.v); (err != nil) != tt.wantErr {
				t.Errorf("TestValidateVerification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateService(t *testing.T) {

	tests := []struct {
		s       *Service
		wantErr bool
	}{
		{
			s:       NewService("agent:abc", "DIDCommMessaging", "https://agent.abc/abc"),
			wantErr: false,
		},
		{
			s:       nil,
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestValidateService#", i), func(t *testing.T) {
			if err := ValidateService(tt.s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"    a    ", false},
		{"\t", true},
		{"\n", true},
		{"   ", true},
		{"  \t \n", true},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestIsEmpty#", i), func(t *testing.T) {
			if got := IsEmpty(tt.input); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDidDocument(t *testing.T) {
	type params struct {
		id      string
		options []DidDocumentOption
	}
	tests := []struct {
		params  params
		wantDid DidDocument
		wantErr bool
	}{
		{
			params: params{
				"did:mantrachain:subject",
				[]DidDocumentOption{
					WithVerifications(
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1",
								"did:mantrachain:subject",
								NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
								KeyAgreement,
							},
							[]string{
								"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
							},
						),
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1", // duplicate key
								"did:mantrachain:subject",
								NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
								KeyAgreement,
							},
							[]string{},
						),
					),
					WithServices(&Service{
						"agent:xyz",
						"DIDCommMessaging",
						"https://agent.xyz/1234",
					}),
				},
			},
			wantDid: DidDocument{},
			wantErr: true, // the error is caused by duplicated verification method id
		},
		{
			params: params{
				"did:mantrachain:subject",
				[]DidDocumentOption{
					WithVerifications(
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1",
								"did:mantrachain:subject",
								NewPublicKeyMultibase([]byte("02503c8ace59c085b15c5f9c2474e9235bcb9694f07516bdc06f7caec788c3dd2c"), DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
								KeyAgreement,
							},
							[]string{
								"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
							},
						),
					),
					WithServices(
						&Service{
							"agent:xyz",
							"DIDCommMessaging",
							"https://agent.xyz/1234",
						},
						&Service{
							"agent:xyz",
							"DIDCommMessaging",
							"https://agent.xyz/1234",
						},
					),
				},
			},
			wantDid: DidDocument{},
			wantErr: true, //duplicated service id
		},
		{
			wantErr: true, // invalid did
			params: params{
				id:      "something not right",
				options: []DidDocumentOption{},
			},
			wantDid: DidDocument{},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestNewDidDocument#", i), func(t *testing.T) {
			gotDid, err := NewDidDocument(tt.params.id, tt.params.options...)

			if tt.wantErr {
				require.NotNil(t, err, "test: TestNewDidDocument#%v", i)
				return
			}

			require.Nil(t, err, "test: TestNewDidDocument#%v", i)
			assert.Equal(t, tt.wantDid, gotDid)
		})
	}
}

func TestDidDocument_AddControllers(t *testing.T) {

	tests := []struct {
		name                string
		malleate            func() DidDocument
		controllers         []string
		expectedControllers []string
		wantErr             error
	}{
		{
			"PASS: controllers added",
			func() DidDocument {
				dd, _ := NewDidDocument("did:mantrachain:subject",
					WithControllers(
						"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
						"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8", // duplicate controllers
					),
				)
				return dd
			},
			[]string{
				"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
			},
			[]string{
				"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
				"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
			},
			nil,
		},
		{
			"FAIL: invalid controller did",
			func() DidDocument {
				dd, _ := NewDidDocument("did:mantrachain:subject",
					WithControllers(
						"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
					),
				)
				return dd
			},
			[]string{
				"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
				"not a did:cosmos:key:cosmos100000000000000000000000000000000000004",
			},
			[]string{},
			sdkerrors.Wrapf(ErrInvalidDIDFormat, "did document controller validation error 'not a did:cosmos:key:cosmos100000000000000000000000000000000000004'"),
		},
		{
			"FAIL: controller did is not type key",
			func() DidDocument {
				dd, _ := NewDidDocument("did:mantrachain:subject",
					WithControllers(
						"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
					),
				)
				return dd
			},
			[]string{
				"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
				"did:cosmos:net:foochain:1234",
			},
			[]string{},
			sdkerrors.Wrapf(ErrInvalidInput, "did document controller 'did:cosmos:net:foochain:1234' must be of type key"),
		},
		{
			"PASS: controllers empty",
			func() DidDocument {
				dd, _ := NewDidDocument("did:mantrachain:subject",
					WithControllers(
						"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
						"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
					),
				)
				return dd
			},
			nil,
			[]string{
				"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
				"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
			},
			nil,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestDidDocument_AddControllers#", i), func(t *testing.T) {
			didDoc := tt.malleate()
			err := didDoc.AddControllers(tt.controllers...)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			}
		})
	}
}

func TestDidDocument_DeleteControllers(t *testing.T) {

	tests := []struct {
		name                string
		malleate            func() DidDocument
		controllers         []string
		expectedControllers []string
		wantErr             error
	}{
		{
			"PASS: controllers deleted",
			func() DidDocument {
				dd, _ := NewDidDocument("did:mantrachain:subject",
					WithControllers(
						"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
						"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
					),
				)
				return dd
			},
			[]string{
				"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
			},
			[]string{
				"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
			},
			nil,
		},
		{
			"FAIL: invalid controller did",
			func() DidDocument {
				dd, _ := NewDidDocument("did:mantrachain:subject",
					WithControllers(
						"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
					),
				)
				return dd
			},
			[]string{
				"not a did:cosmos:key:cosmos100000000000000000000000000000000000004",
			},
			[]string{},
			sdkerrors.Wrapf(ErrInvalidDIDFormat, "did document controller validation error 'not a did:cosmos:key:cosmos100000000000000000000000000000000000004'"),
		},
		{
			"PASS: controllers empty",
			func() DidDocument {
				dd, _ := NewDidDocument("did:mantrachain:subject",
					WithControllers(
						"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
						"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
					),
				)
				return dd
			},
			nil,
			[]string{
				"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
				"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
			},
			nil,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestDidDocument_AddControllers#", i), func(t *testing.T) {
			didDoc := tt.malleate()
			err := didDoc.DeleteControllers(tt.controllers...)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			}
		})
	}
}

func TestDidDocument_AddVerifications(t *testing.T) {
	type params struct {
		malleate      func() DidDocument // build a did document
		verifications []*Verification    // input list of verifications
	}
	tests := []struct {
		params  params
		wantDid DidDocument // expected result
		wantErr bool
	}{
		{
			wantErr: true, // duplicated existing method id
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1",
								"did:mantrachain:subject",
								NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
								KeyAgreement,
								KeyAgreement, // test duplicated relationship
							},
							[]string{
								"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
							},
						),
					))
					return d
				},
				[]*Verification{
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{
							string(CapabilityDelegation),
						},
						[]string{
							"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
						},
					),
				},
			},
			wantDid: DidDocument{},
		},
		{
			wantErr: true, // duplicated new method id
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1",
								"did:mantrachain:subject",
								NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
								KeyAgreement,
								KeyAgreement, // test duplicated relationship
							},
							[]string{
								"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
							},
						),
					))
					return d
				},
				[]*Verification{
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-2",
							"did:mantrachain:subject",
							NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{
							KeyAgreement,
						},
						[]string{
							"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
						},
					),
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-2",
							"did:mantrachain:subject",
							NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{
							Authentication,
						},
						[]string{
							"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
						},
					),
				},
			},
			wantDid: DidDocument{},
		},
		{
			wantErr: true, // fail validation
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1",
								"did:mantrachain:subject",
								NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
								KeyAgreement,
								KeyAgreement, // test duplicated relationship
							},
							[]string{
								"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
							},
						),
					))
					return d
				},
				[]*Verification{
					{
						[]string{
							string(Authentication),
							string(KeyAgreement),
							string(KeyAgreement), // test duplicated relationship
						},
						&VerificationMethod{
							"invalid method url",
							DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
							"did:mantrachain:subject",
							&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
						},
						[]string{
							"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
						},
					},
				},
			},
			wantDid: DidDocument{},
		},
		{
			wantErr: true, // verification relationship does not exists
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject")
					return d
				},
				[]*Verification{
					{
						[]string{
							Authentication,
							"UNSUPPORTED RELATIONSHIP",
							KeyAgreement,
						},
						&VerificationMethod{
							"did:mantrachain:subject#key1",
							DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
							"did:mantrachain:subject",
							&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
						},
						[]string{
							"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
						},
					},
				},
			},
			wantDid: DidDocument{},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestDidDocument_AddVerifications#", i), func(t *testing.T) {
			gotDid := tt.params.malleate()

			err := gotDid.AddVerifications(tt.params.verifications...)

			if tt.wantErr {
				require.NotNil(t, err, "test: TestDidDocument_AddVerifications#%v", i)
				return
			}

			require.Nil(t, err, "test: TestDidDocument_AddVerifications#%v", i)
			assert.Equal(t, tt.wantDid, gotDid)
		})
	}
}

func TestDidDocument_RevokeVerification(t *testing.T) {
	type params struct {
		malleate func() DidDocument // build a did document
		methodID string             // input list of verifications
	}
	tests := []struct {
		params  params
		wantDid DidDocument // expected result
		wantErr bool
	}{
		{
			wantErr: false,
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject",
						WithVerifications(
							NewVerification(
								VerificationMethod{
									"did:mantrachain:subject#key-1",
									DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
									"did:mantrachain:subject",
									&VerificationMethod_PublicKeyMultibase{"F03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
								},
								[]string{
									Authentication,
									KeyAgreement,
								},
								nil,
							),
						),
					)
					return d
				},
				"did:mantrachain:subject#key-1",
			},
			wantDid: DidDocument{
				Context: []string{
					contextDIDBase,
				},
				Id: "did:mantrachain:subject",
			},
		},
		{
			wantErr: false,
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject",
						WithVerifications(
							NewVerification(
								VerificationMethod{
									"did:mantrachain:subject#key-1",
									DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
									"did:mantrachain:subject",
									&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
								},
								[]string{
									Authentication,
									KeyAgreement,
								},
								nil,
							),
							NewVerification(
								VerificationMethod{
									"did:mantrachain:subject#key-2",
									DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
									"did:mantrachain:subject",
									&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
								},
								[]string{
									Authentication,
									CapabilityInvocation,
								},
								nil,
							),
							NewVerification(
								VerificationMethod{
									"did:mantrachain:subject#key-3",
									DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
									"did:mantrachain:subject",
									&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
								},
								[]string{
									Authentication,
									KeyAgreement,
									AssertionMethod,
								},
								nil,
							),
						),
					)
					return d
				},
				"did:mantrachain:subject#key-2",
			},
			wantDid: DidDocument{
				Context: []string{
					contextDIDBase,
				},
				Id:         "did:mantrachain:subject",
				Controller: nil,
				VerificationMethod: []*VerificationMethod{
					{
						"did:mantrachain:subject#key-1",
						DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
						"did:mantrachain:subject",
						&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
					},
					{
						"did:mantrachain:subject#key-3",
						DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
						"did:mantrachain:subject",
						&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
					},
				},
				Service:         nil,
				Authentication:  []string{"did:mantrachain:subject#key-1", "did:mantrachain:subject#key-3"},
				KeyAgreement:    []string{"did:mantrachain:subject#key-1", "did:mantrachain:subject#key-3"},
				AssertionMethod: []string{"did:mantrachain:subject#key-3"},
			},
		},
		{
			wantErr: true, // verification method not found
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject",
						WithVerifications(
							NewVerification(
								NewVerificationMethod(
									"did:mantrachain:subject#key-1",
									"did:mantrachain:subject",
									NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
								),
								[]string{
									Authentication,
									KeyAgreement,
								},
								nil,
							),
							NewVerification(
								NewVerificationMethod(
									"did:mantrachain:subject#key-2",
									"did:mantrachain:subject",
									NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
								),
								[]string{
									Authentication,
									CapabilityInvocation,
								},
								[]string{
									"https://gpg.jsld.org/contexts/lds-gpg2020-v0.0.jsonld",
								},
							),
						),
					)
					return d
				},
				"did:mantrachain:subject#key-3",
			},
			wantDid: DidDocument{},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestDidDocument_RevokeVerification#", i), func(t *testing.T) {
			gotDid := tt.params.malleate()

			err := gotDid.RevokeVerification(tt.params.methodID)

			if tt.wantErr {
				require.NotNil(t, err, "test: TestDidDocument_RevokeVerification#%v", i)
				return
			}

			require.Nil(t, err, "test: TestDidDocument_RevokeVerification#%v", i)

			assert.Equal(t, tt.wantDid, gotDid)
		})
	}
}

func TestDidDocument_SetVerificationRelationships(t *testing.T) {
	type params struct {
		malleate      func() DidDocument
		methodID      string
		relationships []string
	}
	tests := []struct {
		params  params
		wantDid DidDocument // expected result
		wantErr bool
	}{
		{
			wantErr: true, // empty relationships
			params: params{
				malleate: func() DidDocument {
					dd, _ := NewDidDocument("did:mantrachain:subject")
					return dd
				},
				methodID:      "did:mantrachain:subject#key-1",
				relationships: []string{},
			},
			wantDid: DidDocument{
				Context: []string{contextDIDBase},
				Id:      "did:mantrachain:subject",
			},
		},
		{
			wantErr: true, //invalid method id
			params: params{
				malleate: func() DidDocument {
					dd, _ := NewDidDocument("did:mantrachain:subject")
					return dd
				},
				methodID:      "did:mantrachain:subject#key-1 invalid ",
				relationships: []string{},
			},
			wantDid: DidDocument{},
		},
		{
			wantErr: false,
			params: params{
				malleate: func() DidDocument {
					dd, _ := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							VerificationMethod{
								"did:mantrachain:subject#key-1",
								DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
								"did:mantrachain:subject",
								&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
							},
							[]string{
								Authentication,
								KeyAgreement,
							},
							[]string{},
						),
					))
					return dd
				},
				methodID: "did:mantrachain:subject#key-1",
				relationships: []string{
					string(AssertionMethod),
					string(AssertionMethod), // test duplicated relationship
					string(AssertionMethod), // test duplicated relationship
					string(AssertionMethod), // test duplicated relationship
				},
			},

			wantDid: DidDocument{
				Context: []string{contextDIDBase},
				Id:      "did:mantrachain:subject",
				VerificationMethod: []*VerificationMethod{
					{
						"did:mantrachain:subject#key-1",
						DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
						"did:mantrachain:subject",
						&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
					},
				},
				AssertionMethod: []string{"did:mantrachain:subject#key-1"},
			},
		},
		{
			wantErr: false, // different delete scenarios
			params: params{
				malleate: func() DidDocument {
					dd, _ := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							VerificationMethod{
								"did:mantrachain:subject#key-1",
								DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
								"did:mantrachain:subject",
								&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
							},
							[]string{
								Authentication,
								KeyAgreement,
							},
							[]string{},
						),
						NewVerification(
							VerificationMethod{
								"did:mantrachain:subject#key-2",
								DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
								"did:mantrachain:subject",
								&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
							},
							[]string{
								Authentication,
							},
							[]string{},
						),
					))
					return dd
				},
				methodID:      "did:mantrachain:subject#key-1",
				relationships: []string{string(AssertionMethod)},
			},
			wantDid: DidDocument{
				Context: []string{contextDIDBase},
				Id:      "did:mantrachain:subject",
				VerificationMethod: []*VerificationMethod{
					{
						"did:mantrachain:subject#key-1",
						DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
						"did:mantrachain:subject",
						&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
					},
					{
						"did:mantrachain:subject#key-2",
						DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
						"did:mantrachain:subject",
						&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
					},
				},
				Authentication:  []string{"did:mantrachain:subject#key-2"},
				AssertionMethod: []string{"did:mantrachain:subject#key-1"},
			},
		},
		{
			wantErr: false, // different delete scenarios
			params: params{
				malleate: func() DidDocument {
					dd, _ := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							VerificationMethod{
								"did:mantrachain:subject#key-2",
								DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
								"did:mantrachain:subject",
								&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
							},
							[]string{
								Authentication,
							},
							[]string{},
						),
						NewVerification(
							VerificationMethod{
								"did:mantrachain:subject#key-3",
								DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
								"did:mantrachain:subject",
								&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
							},
							[]string{
								Authentication,
							},
							[]string{},
						),
						NewVerification(
							VerificationMethod{
								"did:mantrachain:subject#key-1",
								DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
								"did:mantrachain:subject",
								&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
							},
							[]string{
								Authentication,
								KeyAgreement,
							},
							[]string{},
						),
					))
					return dd
				},
				methodID:      "did:mantrachain:subject#key-1",
				relationships: []string{string(AssertionMethod)},
			},
			wantDid: DidDocument{
				Context: []string{contextDIDBase},
				Id:      "did:mantrachain:subject",
				VerificationMethod: []*VerificationMethod{
					{
						"did:mantrachain:subject#key-2",
						DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
						"did:mantrachain:subject",
						&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
					},
					{
						"did:mantrachain:subject#key-3",
						DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
						"did:mantrachain:subject",
						&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
					},
					{
						"did:mantrachain:subject#key-1",
						DIDVMethodTypeEcdsaSecp256k1VerificationKey2019.String(),
						"did:mantrachain:subject",
						&VerificationMethod_PublicKeyHex{"03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"},
					},
				},

				Authentication:  []string{"did:mantrachain:subject#key-2", "did:mantrachain:subject#key-3"},
				AssertionMethod: []string{"did:mantrachain:subject#key-1"},
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestDidDocument_SetVerificationRelationships#", i), func(t *testing.T) {
			didDoc := tt.params.malleate()
			err := didDoc.SetVerificationRelationships(tt.params.methodID, tt.params.relationships...)

			if tt.wantErr {
				require.NotNil(t, err, "test: TestDidDocument_SetVerificationRelationships#%v", i)
				return
			}

			require.Nil(t, err, "test: TestDidDocument_SetVerificationRelationships#%v", i)
			assert.Equal(t, tt.wantDid, didDoc)

		})
	}
}

func TestDidDocument_HasRelationship(t *testing.T) {

	type params struct {
		didFn         func() DidDocument
		signer        BlockchainAccountID
		relationships []string
	}
	tests := []struct {
		name                    string
		params                  params
		expectedHasRelationship bool
	}{
		{
			name:                    "PASS: has relationships (multibase)",
			expectedHasRelationship: true,
			params: params{
				didFn: func() DidDocument {
					dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1",
								"did:mantrachain:subject",
								NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								string(Authentication),
								string(KeyAgreement),
							},
							nil,
						),
					))
					assert.NoError(t, err)
					return dd
				},
				signer: NewBlockchainAccountID("mantrachain", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
				relationships: []string{
					string(AssertionMethod),
					string(Authentication),
				},
			},
		},
		{
			name:                    "PASS: relationships missing (multibase, blockchainaccountid, hex)",
			expectedHasRelationship: false,
			params: params{
				didFn: func() DidDocument {
					dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1",
								"did:mantrachain:subject",
								NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
								KeyAgreement,
							},
							nil,
						),
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:controller-1#key-2",
								"did:mantrachain:controller-1",
								NewBlockchainAccountID("mantrachain", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
							),
							[]string{
								CapabilityDelegation,
							},
							nil,
						),
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-3",
								"did:mantrachain:subject",
								NewPublicKeyHex([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
								KeyAgreement,
							},
							nil,
						),
					))
					assert.NoError(t, err)
					return dd
				},
				signer: NewBlockchainAccountID("mantrachain", "subject"),
				relationships: []string{
					string(CapabilityDelegation),
				},
			},
		},
		{
			name:                    "PASS: relationships missing (blockchainaccountid)",
			expectedHasRelationship: false,
			params: params{
				didFn: func() DidDocument {
					dd, err := NewDidDocument("did:mantrachain:subject")
					assert.NoError(t, err)
					return dd
				},
				signer: NewBlockchainAccountID("mantrachain", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
				relationships: []string{
					string(CapabilityDelegation),
				},
			},
		},
		{
			name:                    "PASS: relationships missing (Multibase)",
			expectedHasRelationship: false,
			params: params{
				didFn: func() DidDocument {
					dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1",
								"did:mantrachain:subject",
								NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
								KeyAgreement,
							},
							nil,
						),
					))
					assert.NoError(t, err)
					return dd
				},
				signer:        NewBlockchainAccountID("mantrachain", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
				relationships: nil,
			},
		},
		{
			name:                    "PASS: has relationship (BlockchainAccountID)",
			expectedHasRelationship: true,
			params: params{
				didFn: func() DidDocument {
					dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1",
								"did:mantrachain:subject",
								NewPublicKeyMultibase([]byte("00dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7"), DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
							},
							nil,
						),
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-2",
								"did:mantrachain:subject",
								NewBlockchainAccountID("mantrachain", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
							),
							[]string{
								KeyAgreement,
							},
							nil,
						),
					))
					assert.NoError(t, err)
					return dd
				},
				signer: NewBlockchainAccountID("mantrachain", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
				relationships: []string{
					string(KeyAgreement),
				},
			},
		},
		{
			name:                    "PASS:  missing relationship (PublicKeyHex)",
			expectedHasRelationship: false,
			params: params{
				didFn: func() DidDocument {
					dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
						NewVerification(
							NewVerificationMethod(
								"did:mantrachain:subject#key-1",
								"did:mantrachain:subject",
								NewPublicKeyHex([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
							),
							[]string{
								Authentication,
								KeyAgreement,
							},
							nil,
						),
					))
					assert.NoError(t, err)
					return dd
				},
				signer:        NewBlockchainAccountID("mantrachain", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
				relationships: nil,
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestDidDocument_SetVerificationRelationships#", i), func(t *testing.T) {
			didDoc := tt.params.didFn()
			gotHasRelationship := didDoc.HasRelationship(tt.params.signer, tt.params.relationships...)
			assert.Equal(t, tt.expectedHasRelationship, gotHasRelationship)
		})
	}
}

func TestDidDocument_AddServices(t *testing.T) {
	type params struct {
		malleate func() DidDocument // build a did document
		services []*Service         // input list of verifications
	}
	tests := []struct {
		params  params
		wantDid DidDocument // expected result
		wantErr bool
	}{
		{
			wantErr: false,
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject")
					return d
				},
				[]*Service{
					NewService(
						"agent:abc",
						"DIDCommMessaging",
						"https://agent.abc/1234",
					),
					NewService(
						"agent:xyz",
						"DIDCommMessaging",
						"https://agent.xyz/1234",
					),
				},
			},
			wantDid: DidDocument{
				Context: []string{contextDIDBase},
				Id:      "did:mantrachain:subject",
				Service: []*Service{
					NewService(
						"agent:abc",
						"DIDCommMessaging",
						"https://agent.abc/1234",
					),
					NewService(
						"agent:xyz",
						"DIDCommMessaging",
						"https://agent.xyz/1234",
					),
				},
			},
		},
		{
			wantErr: true, // duplicated existing service id
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument(
						"did:mantrachain:subject",
						WithServices(
							NewService(
								"agent:xyz",
								"DIDCommMessaging",
								"https://agent.xyz/1234",
							),
						),
					)
					return d
				},
				[]*Service{
					{
						"agent:abc",
						"DIDCommMessaging",
						"https://agent.abc/1234",
					}, {
						"agent:xyz",
						"DIDCommMessaging",
						"https://agent.xyz/1234",
					},
				},
			},
			wantDid: DidDocument{},
		},
		{
			wantErr: true, // duplicated new service id
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject")
					return d
				},
				[]*Service{
					{
						"agent:xyz",
						"DIDCommMessaging",
						"https://agent.xyz/1234",
					}, {
						"agent:xyz",
						"DIDCommMessaging",
						"https://agent.xyz/1234",
					},
				},
			},
			wantDid: DidDocument{},
		},
		{
			wantErr: true, // fail validation
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject")
					return d
				},
				[]*Service{
					{
						"agent:abc",
						"DIDCommMessaging",
						"https://agent.abc/1234",
					}, {
						"",
						"DIDCommMessaging",
						"https://agent.xyz/1234",
					},
				},
			},
			wantDid: DidDocument{},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestDidDocument_AddServices#", i), func(t *testing.T) {
			gotDid := tt.params.malleate()

			err := gotDid.AddServices(tt.params.services...)

			if tt.wantErr {
				require.NotNil(t, err, "test: TestDidDocument_AddServices#%v", i)
				return
			}

			require.Nil(t, err, "test: TestDidDocument_AddServices#%v", i)
			assert.Equal(t, tt.wantDid, gotDid)
		})
	}
}

func TestDidDocument_DeleteService(t *testing.T) {
	type params struct {
		didFn    func() DidDocument // build a did document
		methodID string             // input list of verifications
	}
	tests := []struct {
		params  params
		wantDid DidDocument // expected result
		wantErr bool
	}{
		{
			wantErr: false,
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject",
						WithServices(
							&Service{
								"agent:abc",
								"DIDCommMessaging",
								"https://agent.abc/1234",
							},
						),
					)
					return d
				},
				"agent:abc",
			},
			wantDid: DidDocument{
				Context: []string{contextDIDBase},
				Id:      "did:mantrachain:subject",
			},
		},
		{
			wantErr: false,
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject",
						WithServices(
							&Service{
								"agent:zyz",
								"DIDCommMessaging",
								"https://agent.abc/1234",
							},
							&Service{
								"agent:abc",
								"DIDCommMessaging",
								"https://agent.abc/1234",
							},
						),
					)
					return d
				},
				"agent:abc",
			},
			wantDid: DidDocument{
				Context: []string{contextDIDBase},
				Id:      "did:mantrachain:subject",
				Service: []*Service{
					{
						"agent:zyz",
						"DIDCommMessaging",
						"https://agent.abc/1234",
					},
				},
			},
		},
		{
			wantErr: false,
			params: params{
				func() DidDocument {
					d, _ := NewDidDocument("did:mantrachain:subject",
						WithServices(
							&Service{
								"agent:zyz",
								"DIDCommMessaging",
								"https://agent.abc/1234",
							},
							&Service{
								"agent:abc",
								"DIDCommMessaging",
								"https://agent.abc/1234",
							},
							&Service{
								"agent:007",
								"DIDCommMessaging",
								"https://agent.abc/007",
							},
						),
					)
					return d
				},
				"agent:abc",
			},
			wantDid: DidDocument{
				Context: []string{contextDIDBase},
				Id:      "did:mantrachain:subject",
				Service: []*Service{
					{
						"agent:zyz",
						"DIDCommMessaging",
						"https://agent.abc/1234",
					}, {
						"agent:007",
						"DIDCommMessaging",
						"https://agent.abc/007",
					},
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint("TestDidDocument_DeleteService#", i), func(t *testing.T) {
			gotDid := tt.params.didFn()

			gotDid.DeleteService(tt.params.methodID)

			assert.Equal(t, tt.wantDid, gotDid)
		})
	}
}

func TestBlockchainAccountID_GetAddress(t *testing.T) {
	tests := []struct {
		name string
		baID BlockchainAccountID
		want string
	}{
		{
			"PASS: can get address",
			BlockchainAccountID("cosmos:foochain:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
			"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
		},
		{
			// TODO: this should result in an error
			"PASS: address is empty",
			BlockchainAccountID("cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
			"",
		},
		{
			// TODO: this should result in an error
			"PASS: can get address (but address is wrong)",
			BlockchainAccountID("cosmos:foochain:whatever"),
			"whatever",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.baID.GetAddress(), "GetAddress()")
		})
	}
}

func TestNewPublicKeyMultibaseFromHex(t *testing.T) {
	type args struct {
		pubKeyHex string
		vmType    VerificationMaterialType
	}
	tests := []struct {
		name    string
		args    args
		wantPkm PublicKeyMultibase
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"PASS: key match",
			args{
				pubKeyHex: "03dfd0a469806d66a23c7c948f55c129467d6d0974a222ef6e24a538fa6882f3d7",
				vmType:    DIDVMethodTypeEcdsaSecp256k1VerificationKey2019,
			},
			PublicKeyMultibase{
				data:   []byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215},
				vmType: DIDVMethodTypeEcdsaSecp256k1VerificationKey2019,
			},
			assert.NoError,
		},
		{
			"FAIL: invalid hex key",
			args{
				pubKeyHex: "not hex string",
				vmType:    DIDVMethodTypeEcdsaSecp256k1VerificationKey2019,
			},
			PublicKeyMultibase{
				data:   nil,
				vmType: "",
			},
			assert.Error, // TODO: check the error message
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPkm, err := NewPublicKeyMultibaseFromHex(tt.args.pubKeyHex, tt.args.vmType)
			if !tt.wantErr(t, err, fmt.Sprintf("NewPublicKeyMultibaseFromHex(%v, %v)", tt.args.pubKeyHex, tt.args.vmType)) {
				return
			}
			assert.Equalf(t, tt.wantPkm, gotPkm, "NewPublicKeyMultibaseFromHex(%v, %v)", tt.args.pubKeyHex, tt.args.vmType)

		})
	}
}

func TestDidDocument_HasPublicKey(t *testing.T) {

	tests := []struct {
		name   string
		didFn  func() DidDocument
		pubkey func() types.PubKey
		want   bool
	}{
		{
			"PASS: has public key (multibase)",
			func() DidDocument {
				dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewPublicKeyMultibase([]byte{2, 201, 95, 248, 187, 133, 206, 97, 166, 70, 229, 226, 88, 124, 29, 43, 70, 3, 244, 225, 19, 128, 44, 132, 110, 15, 15, 35, 40, 189, 237, 71, 245}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{
							Authentication,
							KeyAgreement,
						},
						nil,
					),
				))
				assert.NoError(t, err)
				return dd
			},
			func() types.PubKey {
				var pk types.PubKey
				c := simapp.MakeTestEncodingConfig().Marshaler
				err := c.UnmarshalInterfaceJSON([]byte(`{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Aslf+LuFzmGmRuXiWHwdK0YD9OETgCyEbg8PIyi97Uf1"}`), &pk)
				assert.NoError(t, err)
				return pk

			},
			true,
		},
		{
			"PASS: doesn't have public key (multibase)",
			func() DidDocument {
				dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{
							Authentication,
							KeyAgreement,
						},
						nil,
					),
				))
				assert.NoError(t, err)
				return dd
			},
			func() types.PubKey {
				var pk types.PubKey
				c := simapp.MakeTestEncodingConfig().Marshaler
				err := c.UnmarshalInterfaceJSON([]byte(`{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Aslf+LuFzmGmRuXiWHwdK0YD9OETgCyEbg8PIyi97Uf1"}`), &pk)
				assert.NoError(t, err)
				return pk

			},
			false,
		},
		{
			"PASS: has public key (blockchainAccount)",
			func() DidDocument {
				dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewBlockchainAccountID("foochain", "cosmos17t8t3t6a6vpgk69perfyq930593sa8dn4kzsdf"),
						),
						[]string{
							Authentication,
							KeyAgreement,
						},
						nil,
					),
				))
				assert.NoError(t, err)
				return dd
			},
			func() types.PubKey {
				var pk types.PubKey
				c := simapp.MakeTestEncodingConfig().Marshaler
				err := c.UnmarshalInterfaceJSON([]byte(`{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Aslf+LuFzmGmRuXiWHwdK0YD9OETgCyEbg8PIyi97Uf1"}`), &pk)
				assert.NoError(t, err)
				return pk

			},
			true,
		},
		{
			"PASS: doesn't have public key (blockchainAccountId)",
			func() DidDocument {
				dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewBlockchainAccountID("foochain", "cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2"),
						),
						[]string{
							Authentication,
							KeyAgreement,
						},
						nil,
					),
				))
				assert.NoError(t, err)
				return dd
			},
			func() types.PubKey {
				var pk types.PubKey
				c := simapp.MakeTestEncodingConfig().Marshaler
				err := c.UnmarshalInterfaceJSON([]byte(`{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Aslf+LuFzmGmRuXiWHwdK0YD9OETgCyEbg8PIyi97Uf1"}`), &pk)
				assert.NoError(t, err)
				return pk

			},
			false,
		},
		{
			"PASS: has public key (publicKeyHex)",
			func() DidDocument {
				dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewPublicKeyHex([]byte{2, 201, 95, 248, 187, 133, 206, 97, 166, 70, 229, 226, 88, 124, 29, 43, 70, 3, 244, 225, 19, 128, 44, 132, 110, 15, 15, 35, 40, 189, 237, 71, 245}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{
							Authentication,
							KeyAgreement,
						},
						nil,
					),
				))
				assert.NoError(t, err)
				return dd
			},
			func() types.PubKey {
				var pk types.PubKey
				c := simapp.MakeTestEncodingConfig().Marshaler
				err := c.UnmarshalInterfaceJSON([]byte(`{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Aslf+LuFzmGmRuXiWHwdK0YD9OETgCyEbg8PIyi97Uf1"}`), &pk)
				assert.NoError(t, err)
				return pk

			},
			true,
		},
		{
			"PASS: doesn't have public key (pubKeyHex)",
			func() DidDocument {
				dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewPublicKeyHex([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{
							Authentication,
							KeyAgreement,
						},
						nil,
					),
				))
				assert.NoError(t, err)
				return dd
			},
			func() types.PubKey {
				var pk types.PubKey
				c := simapp.MakeTestEncodingConfig().Marshaler
				err := c.UnmarshalInterfaceJSON([]byte(`{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Aslf+LuFzmGmRuXiWHwdK0YD9OETgCyEbg8PIyi97Uf1"}`), &pk)
				assert.NoError(t, err)
				return pk

			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			didDoc := tt.didFn()
			pubKey := tt.pubkey()
			assert.Equalf(t, tt.want, didDoc.HasPublicKey(pubKey), "HasPublicKey(%v)", pubKey)
		})
	}
}

func TestDidDocument_GetVerificationMethodBlockchainAddress(t *testing.T) {
	tests := []struct {
		name        string
		didFn       func() DidDocument
		methodID    string
		wantAddress string
		wantErr     error
	}{
		{
			"PASS: get address (PublicKeyMultibase)",
			func() DidDocument {
				dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewPublicKeyMultibase([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{
							Authentication,
							KeyAgreement,
						},
						nil,
					),
				))
				assert.NoError(t, err)
				return dd
			},
			"did:mantrachain:subject#key-1",
			"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
			nil,
		},
		{
			"PASS: get address (PublicKeyHex)",
			func() DidDocument {
				dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewPublicKeyHex([]byte{3, 223, 208, 164, 105, 128, 109, 102, 162, 60, 124, 148, 143, 85, 193, 41, 70, 125, 109, 9, 116, 162, 34, 239, 110, 36, 165, 56, 250, 104, 130, 243, 215}, DIDVMethodTypeEcdsaSecp256k1VerificationKey2019),
						),
						[]string{
							Authentication,
							KeyAgreement,
						},
						nil,
					),
				))
				assert.NoError(t, err)
				return dd
			},
			"did:mantrachain:subject#key-1",
			"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
			nil,
		},
		{
			"PASS: get address (BlockchainAccountID)",
			func() DidDocument {
				dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewBlockchainAccountID("foochain", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
						),
						[]string{
							Authentication,
							KeyAgreement,
						},
						nil,
					),
				))
				assert.NoError(t, err)
				return dd
			},
			"did:mantrachain:subject#key-1",
			"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
			nil,
		},
		{
			"PASS: get address (BlockchainAccountID)",
			func() DidDocument {
				dd, err := NewDidDocument("did:mantrachain:subject", WithVerifications(
					NewVerification(
						NewVerificationMethod(
							"did:mantrachain:subject#key-1",
							"did:mantrachain:subject",
							NewBlockchainAccountID("foochain", "cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8"),
						),
						[]string{
							Authentication,
							KeyAgreement,
						},
						nil,
					),
				))
				assert.NoError(t, err)
				return dd
			},
			"did:mantrachain:subject#key-2",
			"cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
			ErrVerificationMethodNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			didDoc := tt.didFn()
			gotAddress, err := didDoc.GetVerificationMethodBlockchainAddress(tt.methodID)
			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantAddress, gotAddress, "GetVerificationMethodBlockchainAddress(%v)", tt.methodID)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			}
		})
	}
}

func TestDidDocument_HasController(t *testing.T) {

	tests := []struct {
		name          string
		didFn         func() DidDocument
		controllerDID string
		want          bool
	}{
		{
			"PASS: controller found",
			func() DidDocument {
				dd, err := NewDidDocument(
					"did:mantrachain:subject",
					WithControllers(
						"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
						"did:cosmos:key:cosmos17t8t3t6a6vpgk69perfyq930593sa8dn4kzsdf",
					),
				)
				assert.NoError(t, err)
				return dd
			},
			"did:cosmos:key:cosmos17t8t3t6a6vpgk69perfyq930593sa8dn4kzsdf",
			true,
		},
		{
			"PASS: controller not found",
			func() DidDocument {
				dd, err := NewDidDocument(
					"did:mantrachain:subject",
					WithControllers(
						"did:cosmos:key:cosmos1lvl2s8x4pta5f96appxrwn3mypsvumukvk7ck2",
						"did:cosmos:key:cosmos17t8t3t6a6vpgk69perfyq930593sa8dn4kzsdf",
					),
				)
				assert.NoError(t, err)
				return dd
			},
			"did:cosmos:key:cosmos1sl48sj2jjed7enrv3lzzplr9wc2f5js5tzjph8",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			didDoc := tt.didFn()
			assert.Equalf(t, tt.want, didDoc.HasController(DID(tt.controllerDID)), "HasController(%v)", tt.controllerDID)
		})
	}
}
