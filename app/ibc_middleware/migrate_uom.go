package ibc_middleware

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
)

const (
	UOM          = "uom"
	AMANTRA      = "amantra"
	SCALE_FACTOR = 4_000_000_000_000 // 4 * 10^12
)

var _ porttypes.IBCModule = MigrateUomIBCModule{}

type MigrateUomIBCModule struct {
	// Since this is the last middleware in the stack, `app` is the core `transfer` IBC module.
	app        porttypes.IBCModule
	bankkeeper bankkeeper.Keeper
}

func NewMigrateUomIBCModule(app porttypes.IBCModule, bankkeeper bankkeeper.Keeper) MigrateUomIBCModule {
	return MigrateUomIBCModule{
		app,
		bankkeeper,
	}
}

// OnChanOpenInit implements the IBCModule interface
func (im MigrateUomIBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	return im.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, counterparty, version)
}

// OnChanOpenTry implements the IBCModule interface
func (im MigrateUomIBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	return im.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, counterparty, counterpartyVersion)
}

// OnChanOpenAck implements the IBCModule interface
func (im MigrateUomIBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	return im.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

// OnChanOpenConfirm implements the IBCModule interface
func (im MigrateUomIBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return im.app.OnChanOpenConfirm(ctx, portID, channelID)
}

// OnChanCloseInit implements the IBCModule interface
func (im MigrateUomIBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return im.app.OnChanCloseInit(ctx, portID, channelID)
}

// OnChanCloseConfirm implements the IBCModule interface
func (im MigrateUomIBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return im.app.OnChanCloseConfirm(ctx, portID, channelID)
}

// OnRecvPacket implements the IBCModule interface
func (im MigrateUomIBCModule) OnRecvPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	// First, call the underlying app's OnRecvPacket to handle the transfer.
	ack := im.app.OnRecvPacket(ctx, channelVersion, packet, relayer)
	if !ack.Success() {
		return ack
	}

	// The transfer was successful. Now, perform the uom to amantra migration.
	var data transfertypes.InternalTransferRepresentation
	data, ackErr := transfertypes.UnmarshalPacketData(packet.GetData(), channelVersion, "")
	if ackErr != nil {
		ack = channeltypes.NewErrorAcknowledgement(ackErr)
		return ack
	}

	token := data.Token
	var uomCoin sdk.Coin

	if token.Denom.HasPrefix(packet.SourcePort, packet.SourceChannel) {
		// remove prefix added by sender chain
		token.Denom.Trace = token.Denom.Trace[1:]

		mantraDenom := token.Denom.IBCDenom()
		if mantraDenom != UOM {
			// Not a 'uom' token, no migration needed.
			return ack
		}
		transferAmount, ok := math.NewIntFromString(token.Amount)
		if !ok {
			return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "unable to parse transfer amount: %s", token.Amount))
		}
		uomCoin = sdk.NewCoin(UOM, transferAmount)
	}

	// Get the recipient's address.
	recipient, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(err, "cannot parse recipient address"))
	}

	// 1. Burn the received 'uom' voucher from the recipient's account.
	if err := im.bankkeeper.SendCoinsFromAccountToModule(ctx, recipient, transfertypes.ModuleName, sdk.NewCoins(uomCoin)); err != nil {
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(err, "failed to send uom to migration module"))
	}
	if err := im.bankkeeper.BurnCoins(ctx, transfertypes.ModuleName, sdk.NewCoins(uomCoin)); err != nil {
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(err, "failed to burn uom from migration module"))
	}

	// 2. Mint 4e12 times the amount in 'amantra'.
	scalingFactor := math.NewInt(SCALE_FACTOR) // 4 * 10^12
	amantraAmount := uomCoin.Amount.Mul(scalingFactor)
	amantraCoin := sdk.NewCoin(AMANTRA, amantraAmount)

	if err := im.bankkeeper.MintCoins(ctx, transfertypes.ModuleName, sdk.NewCoins(amantraCoin)); err != nil {
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(err, "failed to mint amantra in migration module"))
	}

	if err := im.bankkeeper.SendCoinsFromModuleToAccount(ctx, transfertypes.ModuleName, recipient, sdk.NewCoins(amantraCoin)); err != nil {
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(err, "failed to send amantra from migration module to recipient"))
	}

	return ack
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im MigrateUomIBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return im.app.OnAcknowledgementPacket(ctx, channelVersion, packet, acknowledgement, relayer)
}

// OnTimeoutPacket implements the IBCModule interface
func (im MigrateUomIBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return im.app.OnTimeoutPacket(ctx, channelVersion, packet, relayer)
}
