package ibc_middleware

import (
	"context"
	"strings"

	erc20types "github.com/cosmos/evm/x/erc20/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"

	"cosmossdk.io/core/address"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"

	"github.com/ethereum/go-ethereum/common"
)

const (
	EventTypeAutoConvertERC20Coin = "auto_convert_erc20_coin"
)

type autoConvertAccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, account sdk.AccountI)
}

type autoConvertBankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type autoConvertERC20MsgServer interface {
	ConvertCoin(goCtx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
}

var (
	_ porttypes.IBCModule             = AutoConvertERC20CoinIBCModule{}
	_ porttypes.PacketDataUnmarshaler = AutoConvertERC20CoinIBCModule{}
)

type AutoConvertERC20CoinIBCModule struct {
	app           porttypes.IBCModule
	accountKeeper autoConvertAccountKeeper
	bankKeeper    autoConvertBankKeeper
	erc20Server   autoConvertERC20MsgServer
	addrCodec     address.Codec
}

func NewAutoConvertERC20CoinIBCModule(
	app porttypes.IBCModule,
	accountKeeper autoConvertAccountKeeper,
	bankKeeper autoConvertBankKeeper,
	erc20Server autoConvertERC20MsgServer,
	addrCodec address.Codec,
) AutoConvertERC20CoinIBCModule {
	return AutoConvertERC20CoinIBCModule{
		app:           app,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		erc20Server:   erc20Server,
		addrCodec:     addrCodec,
	}
}

func (im AutoConvertERC20CoinIBCModule) UnmarshalPacketData(ctx sdk.Context, portID string, channelID string, bz []byte) (interface{}, string, error) {
	if unmarshaler, ok := im.app.(porttypes.PacketDataUnmarshaler); ok {
		return unmarshaler.UnmarshalPacketData(ctx, portID, channelID, bz)
	}
	return nil, "", errorsmod.Wrap(channeltypes.ErrInvalidChannel, "underlying IBC module does not support packet data unmarshaling")
}

func (im AutoConvertERC20CoinIBCModule) OnChanOpenInit(
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

func (im AutoConvertERC20CoinIBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	return im.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, counterparty, counterpartyVersion)
}

func (im AutoConvertERC20CoinIBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID string,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	return im.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (im AutoConvertERC20CoinIBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID string,
	channelID string,
) error {
	return im.app.OnChanOpenConfirm(ctx, portID, channelID)
}

func (im AutoConvertERC20CoinIBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID string,
	channelID string,
) error {
	return im.app.OnChanCloseInit(ctx, portID, channelID)
}

func (im AutoConvertERC20CoinIBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID string,
	channelID string,
) error {
	return im.app.OnChanCloseConfirm(ctx, portID, channelID)
}

func (im AutoConvertERC20CoinIBCModule) OnRecvPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	ack := im.app.OnRecvPacket(ctx, channelVersion, packet, relayer)
	if !ack.Success() {
		return ack
	}

	// Only auto-convert on successful receives.
	var data transfertypes.InternalTransferRepresentation
	data, ackErr := transfertypes.UnmarshalPacketData(packet.GetData(), channelVersion, "")
	if ackErr != nil {
		return channeltypes.NewErrorAcknowledgement(ackErr)
	}

	// Only auto-convert when the receiver is an EVM hex address.
	if !common.IsHexAddress(data.Receiver) {
		return ack
	}

	token := data.Token
	// For returning tokens that originated here, the receiver is credited in the base denom
	// (e.g. `erc20:<addr>`), not an `ibc/<hash>` voucher.
	denomStr := token.Denom.String()
	baseDenom := denomStr
	if i := strings.LastIndex(denomStr, "/"); i >= 0 {
		baseDenom = denomStr[i+1:]
	}
	if !strings.HasPrefix(baseDenom, "erc20:") {
		return ack
	}

	transferAmount, ok := math.NewIntFromString(token.Amount)
	if !ok {
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "unable to parse transfer amount: %s", token.Amount))
	}
	coin := sdk.NewCoin(baseDenom, transferAmount)

	// Resolve recipient bytes; fall back to EVM bytes if the codec can't parse hex.
	recipientBz, err := im.addrCodec.StringToBytes(data.Receiver)
	if err != nil {
		if common.IsHexAddress(data.Receiver) {
			recipientBz = common.HexToAddress(data.Receiver).Bytes()
		} else {
			return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(err, "invalid recipient"))
		}
	}
	recipient := sdk.AccAddress(recipientBz)

	// Only attempt conversion if the receiver actually has the base coin.
	if bal := im.bankKeeper.GetBalance(ctx, recipient, baseDenom); bal.Amount.LT(transferAmount) {
		return ack
	}

	// Ensure the account exists; ConvertCoin burns from the sender's bank balance.
	if im.accountKeeper.GetAccount(ctx, recipient) == nil {
		acc := im.accountKeeper.NewAccountWithAddress(ctx, recipient)
		im.accountKeeper.SetAccount(ctx, acc)
	}

	receiverEvm := common.HexToAddress(data.Receiver)
	msg := erc20types.NewMsgConvertCoin(coin, receiverEvm, recipient)
	if _, err := im.erc20Server.ConvertCoin(ctx, msg); err != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"auto_convert_erc20_coin_failed",
				sdk.NewAttribute("denom", baseDenom),
				sdk.NewAttribute("amount", token.Amount),
				sdk.NewAttribute("receiver", receiverEvm.Hex()),
				sdk.NewAttribute("error", err.Error()),
			),
		)
		return ack
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeAutoConvertERC20Coin,
			sdk.NewAttribute("denom", baseDenom),
			sdk.NewAttribute("amount", token.Amount),
			sdk.NewAttribute("receiver", receiverEvm.Hex()),
		),
	)

	return ack
}

func (im AutoConvertERC20CoinIBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return im.app.OnAcknowledgementPacket(ctx, channelVersion, packet, acknowledgement, relayer)
}

func (im AutoConvertERC20CoinIBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return im.app.OnTimeoutPacket(ctx, channelVersion, packet, relayer)
}
