package tokenfactory

import (
	"github.com/MANTRA-Chain/mantrachain/v7/x/tokenfactory/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/cosmos/ibc-go/v10/modules/core/api"
)

var _ api.IBCModule = (*IBCV2Module)(nil)

type IBCV2Module struct {
	app                api.IBCModule
	tokenfactoryKeeper keeper.Keeper
}

func NewIBCV2Module(app api.IBCModule, tokenfactoryKeeper keeper.Keeper) IBCV2Module {
	return IBCV2Module{
		app:                app,
		tokenfactoryKeeper: tokenfactoryKeeper,
	}
}

func (im IBCV2Module) OnSendPacket(
	ctx sdk.Context,
	sourceClient string,
	destinationClient string,
	sequence uint64,
	payload channeltypesv2.Payload,
	signer sdk.AccAddress,
) error {
	// IBC v2 does not have channel opening handshake, so we need to store the escrow address
	// on every send packet
	escrowAddress := transfertypes.GetEscrowAddress(transfertypes.PortID, sourceClient)
	im.tokenfactoryKeeper.StoreEscrowAddress(ctx, escrowAddress)
	return im.app.OnSendPacket(ctx, sourceClient, destinationClient, sequence, payload, signer)
}

func (im IBCV2Module) OnRecvPacket(
	ctx sdk.Context,
	sourceClient string,
	destinationClient string,
	sequence uint64,
	payload channeltypesv2.Payload,
	relayer sdk.AccAddress,
) channeltypesv2.RecvPacketResult {
	return im.app.OnRecvPacket(ctx, sourceClient, destinationClient, sequence, payload, relayer)
}

func (im IBCV2Module) OnTimeoutPacket(
	ctx sdk.Context,
	sourceClient string,
	destinationClient string,
	sequence uint64,
	payload channeltypesv2.Payload,
	relayer sdk.AccAddress,
) error {
	return im.app.OnTimeoutPacket(ctx, sourceClient, destinationClient, sequence, payload, relayer)
}

func (im IBCV2Module) OnAcknowledgementPacket(
	ctx sdk.Context,
	sourceClient string,
	destinationClient string,
	sequence uint64,
	acknowledgement []byte,
	payload channeltypesv2.Payload,
	relayer sdk.AccAddress,
) error {
	return im.app.OnAcknowledgementPacket(ctx, sourceClient, destinationClient, sequence, acknowledgement, payload, relayer)
}
