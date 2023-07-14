package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"mantrachain/x/did/types"
)

// deriveVMType derive the verification method type from a public key
func deriveVMType(pubKeyType string) (vmType types.VerificationMaterialType, err error) {
	switch pubKeyType {
	case "ed25519":
		vmType = types.DIDVMethodTypeEd25519VerificationKey2018
	case "secp256k1":
		vmType = types.DIDVMethodTypeEcdsaSecp256k1VerificationKey2019
	default:
		err = types.ErrKeyFormatNotSupported
	}
	return
}

func (k Keeper) CreateNewDidDocument(ctx sdk.Context, id string, controller string) (string, error) {
	// check that the did is not already taken
	_, found := k.GetDidDocument(ctx, []byte(id))
	if found {
		err := sdkerrors.Wrapf(types.ErrDidDocumentFound, "a document with did %s already exists", id)
		k.Logger(ctx).Error(err.Error())
		return "", err
	}

	k.Logger(ctx).Info("request to create a did document", "target did", id)

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

	k.Logger(ctx).Info("created did document", "did", id, "controller", controller)

	// emit the event
	if err := ctx.EventManager().EmitTypedEvents(types.NewDidDocumentCreatedEvent(id, controller)); err != nil {
		k.Logger(ctx).Error("failed to emit DidDocumentCreatedEvent", "did", id, "controller", controller, "err", err)
	}

	return didDocument.Id, nil
}

func (k Keeper) ForceRemoveDidDocumentIfExists(ctx sdk.Context, id string) (bool, error) {
	k.Logger(ctx).Info("request to delete a did document if exists", "target did", id)

	found := k.HasDidDocument(ctx, []byte(id))
	if !found {
		return false, nil
	}

	k.DeleteDidDocument(ctx, []byte(id))

	return true, nil
}
