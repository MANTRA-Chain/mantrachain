package keeper

import (
	"github.com/LimeChain/mantrachain/x/mns/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidExecutor struct {
	id         string
	signer     sdk.Address
	pubKeyHex  string
	pubKeyType string
	didId      string
}

func NewDidExecutor(id string, signer sdk.Address, pubKeyHex string, pubKeyType string) *DidExecutor {
	return &DidExecutor{
		id:         id,
		signer:     signer,
		pubKeyHex:  pubKeyHex,
		pubKeyType: pubKeyType,
	}
}

func (c *DidExecutor) SetDid(ctx sdk.Context, didKeeper types.DidDocumentKeeper) (bool, error) {
	didId, err := didKeeper.SetNewDidDocument(ctx, c.id, c.signer, c.pubKeyHex, c.pubKeyType)
	if err != nil {
		return false, err
	}
	c.didId = didId
	return true, nil
}

func (c *DidExecutor) GetDidId() string {
	return c.didId
}
