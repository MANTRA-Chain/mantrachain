package keeper

import (
	"encoding/json"
	"fmt"

	types "github.com/LimeChain/mantrachain/x/mns/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VerificationMethod struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	Controller      string `json:"controller"`
	PubKeyMultiBase string `json:"public_key_multibase"`
}

type DidDoc struct {
	Id                 string               `json:"id"`
	VerificationMethod []VerificationMethod `json:"verification_method"`
	Authentication     string               `json:"authentication"`
}

type DidController struct {
	didId  string
	ctx    sdk.Context
	didDoc DidDoc
}

func NewDidController(ctx sdk.Context, id string) *DidController {
	return &DidController{
		didId: id,
		ctx:   ctx,
	}
}

func (c *DidController) GenDidId() *DidController {
	c.didDoc.Id = fmt.Sprintf("did:%s:%s:%s", types.DidMethod, types.DidNamespace, c.didId)
	return c
}

func (c *DidController) GenDidVerMethod(pubKeyHex string, accAddress string, vmType string) *DidController {
	c.didDoc.VerificationMethod = []VerificationMethod{
		{
			Id:              fmt.Sprintf("%s#%s", c.didDoc.Id, accAddress),
			Type:            vmType,
			Controller:      c.didDoc.Id,
			PubKeyMultiBase: pubKeyHex,
		},
	}
	return c
}

func (c *DidController) GenDidAuth() *DidController {
	c.didDoc.Authentication = c.didDoc.VerificationMethod[0].Id
	return c
}

func (c *DidController) SetDid() error {
	return nil
}

func (c *DidController) GetDidDoc() ([]byte, error) {
	return json.Marshal(c.didDoc)
}
