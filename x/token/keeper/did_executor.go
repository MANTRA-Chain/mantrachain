package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/AumegaChain/aumega/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidExecutor struct {
	controller string
	ctx        sdk.Context
	dk         types.DidKeeper
}

func NewDidExecutor(ctx sdk.Context, controller string, dk types.DidKeeper) *DidExecutor {

	return &DidExecutor{
		controller: controller,
		ctx:        ctx,
		dk:         dk,
	}
}

func (c *DidExecutor) CreateDidForNft(id []byte) (string, error) {
	encoded := sha256.Sum256([]byte(fmt.Sprintf("%s/nft/%s", types.ModuleName, id)))
	didId := hex.EncodeToString(encoded[:])

	res, err := c.dk.CreateNewDidDocument(c.ctx, didId, c.controller)
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
