package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/LimeChain/mantrachain/x/mdb/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidExecutor struct {
	signer     sdk.Address
	pubKeyHex  string
	pubKeyType string
	didId      string
	ctx        sdk.Context
	didKeeper  types.DidKeeper
}

func NewDidExecutor(ctx sdk.Context, signer sdk.Address, pubKeyHex string, pubKeyType string, didKeeper types.DidKeeper) *DidExecutor {

	return &DidExecutor{
		signer:     signer,
		pubKeyHex:  pubKeyHex,
		pubKeyType: pubKeyType,
		ctx:        ctx,
		didKeeper:  didKeeper,
	}
}

func (c *DidExecutor) SetDid(id string) (bool, error) {
	encoded := sha256.Sum256([]byte(fmt.Sprintf("%s&ts=%d", id, c.ctx.BlockTime().Unix())))
	didId := hex.EncodeToString(encoded[:])

	res, err := c.didKeeper.SetNewDidDocument(c.ctx, didId, c.signer, c.pubKeyHex, c.pubKeyType)
	if err != nil {
		return false, err
	}
	c.didId = res
	return true, nil
}

func (c *DidExecutor) GetDidId() string {
	return c.didId
}
