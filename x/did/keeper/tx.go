package keeper

import (
	"cosmossdk.io/errors"
	"github.com/MANTRA-Finance/mantrachain/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateNewDidDocument(ctx sdk.Context, id string, controller string) (string, error) {
	// check that the did is not already taken
	_, found := k.GetDidDocument(ctx, []byte(id))
	if found {
		err := errors.Wrapf(types.ErrDidDocumentFound, "a document with did %s already exists", id)
		k.logger.Error(err.Error())
		return "", err
	}

	k.logger.Info("request to create a did document", "target did", id)

	did := types.NewChainDID(ctx.ChainID(), id)

	didDocument, err := types.NewDidDocument(did.String(),
		types.WithControllers(types.NewKeyDID(controller).String()),
	)
	if err != nil {
		return "", err
	}

	// persist the did document
	k.SetDidDocument(ctx, []byte(id), didDocument)

	// now create and persist the metadata
	didM := types.NewDidMetadata(ctx.TxBytes(), ctx.BlockTime())
	k.SetDidMetadata(ctx, []byte(id), didM)

	k.logger.Info("created did document", "did", id, "controller", controller)

	// emit the event
	if err := ctx.EventManager().EmitTypedEvents(types.NewDidDocumentCreatedEvent(id, controller)); err != nil {
		k.logger.Error("failed to emit DidDocumentCreatedEvent", "did", id, "controller", controller, "err", err)
	}

	return didDocument.Id, nil
}

func (k Keeper) ForceRemoveDidDocumentIfExists(ctx sdk.Context, id string) (bool, error) {
	k.logger.Info("request to delete a did document if exists", "target did", id)

	found := k.HasDidDocument(ctx, []byte(id))
	if !found {
		return false, nil
	}

	k.DeleteDidDocument(ctx, []byte(id))

	return true, nil
}
