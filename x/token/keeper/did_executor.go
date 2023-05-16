package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidExecutor struct {
	signer     sdk.Address
	pubKeyHex  string
	pubKeyType string
	controller sdk.Address
	ctx        sdk.Context
	dk         types.DidKeeper
}

func NewDidExecutor(ctx sdk.Context, signer sdk.Address, pubKeyHex string, pubKeyType string, controller sdk.Address, dk types.DidKeeper) *DidExecutor {

	return &DidExecutor{
		signer:     signer,
		pubKeyHex:  pubKeyHex,
		pubKeyType: pubKeyType,
		controller: controller,
		ctx:        ctx,
		dk:         dk,
	}
}

func (c *DidExecutor) CreateDidForNft(id []byte) (string, error) {
	encoded := sha256.Sum256([]byte(fmt.Sprintf("%s/nft/%s", types.ModuleName, id)))
	didId := hex.EncodeToString(encoded[:])

	res, err := c.dk.CreateNewDidDocument(c.ctx, didId, c.signer, c.pubKeyHex, c.pubKeyType, c.controller)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *DidExecutor) ForceDeleteDidOfNftIfExists(id []byte) (bool, error) {
	encoded := sha256.Sum256([]byte(fmt.Sprintf("%s/nft/%s", types.ModuleName, id)))
	didId := hex.EncodeToString(encoded[:])

	return c.dk.ForceRemoveDidDocumentIfExists(c.ctx, didId)
}
