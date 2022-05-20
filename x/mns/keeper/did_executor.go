package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/LimeChain/mantrachain/x/mns/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DidExecutor struct {
	id         string
	signer     sdk.Address
	pubKeyHex  string
	pubKeyType string
	didId      string
	ctx        sdk.Context
}

func NewDidExecutor(ctx sdk.Context, id string, signer sdk.Address, pubKeyHex string, pubKeyType string) *DidExecutor {
	encoded := sha256.Sum256([]byte(fmt.Sprintf("%s&ts=%d", id, ctx.BlockTime().Unix())))
	didId := hex.EncodeToString(encoded[:])

	return &DidExecutor{
		id:         didId,
		signer:     signer,
		pubKeyHex:  pubKeyHex,
		pubKeyType: pubKeyType,
		ctx:        ctx,
	}
}

func (c *DidExecutor) SetDid(ctx sdk.Context, didKeeper types.DidKeeper) (bool, error) {
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
