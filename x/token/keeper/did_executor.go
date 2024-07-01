package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/MANTRA-Finance/mantrachain/x/token/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidExecutor struct {
	controller string
	ctx        sdk.Context
	didKeeper  types.DidKeeper
}

func NewDidExecutor(ctx sdk.Context, controller string, didKeeper types.DidKeeper) *DidExecutor {

	return &DidExecutor{
		controller: controller,
		ctx:        ctx,
		didKeeper:  didKeeper,
	}
}

func (c *DidExecutor) CreateDidForNft(id []byte) (string, error) {
	encoded := sha256.Sum256([]byte(fmt.Sprintf("%s/nft/%s", types.ModuleName, id)))
	didId := hex.EncodeToString(encoded[:])

	res, err := c.didKeeper.CreateNewDidDocument(c.ctx, didId, c.controller)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *DidExecutor) ForceDeleteDidOfNftIfExists(id []byte) (bool, error) {
	encoded := sha256.Sum256([]byte(fmt.Sprintf("%s/nft/%s", types.ModuleName, id)))
	didId := hex.EncodeToString(encoded[:])

	return c.didKeeper.ForceRemoveDidDocumentIfExists(c.ctx, didId)
}
