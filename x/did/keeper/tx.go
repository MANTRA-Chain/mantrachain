package keeper

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/MANTRA-Finance/mantrachain/x/did/types"
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

func (k Keeper) CreateNewDidDocument(ctx sdk.Context, id string, signer sdk.Address, pubKeyHex string, pubKeyType string, controller sdk.Address) (string, error) {
	// check that the did is not already taken
	_, found := k.GetDidDocument(ctx, []byte(id))
	if found {
		err := sdkerrors.Wrapf(types.ErrDidDocumentFound, "a document with did %s already exists", id)
		k.Logger(ctx).Error(err.Error())
		return "", err
	}

	k.Logger(ctx).Info("request to create a did document", "target did", id)

	pubKeyHex = pubKeyHex[1:]

	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return "", err
	}

	did := types.NewChainDID(ctx.ChainID(), id)
	vmID := types.NewChainDID(ctx.ChainID(), id).NewVerificationMethodID(signer.String())

	vmType, err := deriveVMType(pubKeyType)
	if err != nil {
		return "", err
	}

	auth := types.NewVerification(
		types.NewVerificationMethod(
			vmID,
			did,
			types.NewPublicKeyMultibase(pubKeyBytes, vmType),
		),
		[]string{types.Authentication},
		nil,
	)

	verifications := types.Verifications{auth}

	didDocument, err := types.NewDidDocument(did.String(),
		types.WithVerifications(verifications...),
	)
	if err != nil {
		return "", err
	}

	didDocument.AddControllers(types.NewKeyDID(controller.String()).String())

	// persist the did document
	k.SetDidDocument(ctx, []byte(id), didDocument)

	// now create and persist the metadata
	didM := types.NewDidMetadata(ctx.TxBytes(), ctx.BlockTime())
	k.SetDidMetadata(ctx, []byte(id), didM)

	k.Logger(ctx).Info("created did document", "did", id, "controller", signer.String())

	// emit the event
	if err := ctx.EventManager().EmitTypedEvents(types.NewDidDocumentCreatedEvent(id, signer.String())); err != nil {
		k.Logger(ctx).Error("failed to emit DidDocumentCreatedEvent", "did", id, "signer", signer, "err", err)
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
