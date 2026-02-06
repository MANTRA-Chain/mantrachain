package ibc_middleware

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/MANTRA-Chain/mantrachain/v8/app/evmutil"

	evmibc "github.com/cosmos/evm/ibc"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ porttypes.IBCModule = UnwrapERC20IBCModule{}
var _ porttypes.PacketDataUnmarshaler = UnwrapERC20IBCModule{}

type evmCaller interface {
	CallEVMWithData(
		ctx sdk.Context,
		from common.Address,
		contract *common.Address,
		data []byte,
		commit bool,
		gasCap *big.Int,
	) (*evmtypes.MsgEthereumTxResponse, error)
}

// UnwrapERC20IBCModule optionally unwraps ERC20 wrappers on IBC recv when memo is `{"mantra":{"unwrap":true}}`.
type UnwrapERC20IBCModule struct {
	app         porttypes.IBCModule
	erc20Keeper *erc20keeper.Keeper
	evmCaller   evmCaller
}

func NewUnwrapERC20IBCModule(app porttypes.IBCModule, erc20Keeper *erc20keeper.Keeper, evmCaller evmCaller) UnwrapERC20IBCModule {
	if app == nil {
		panic("underlying application cannot be nil")
	}
	if erc20Keeper == nil {
		panic("erc20 keeper cannot be nil")
	}
	if evmCaller == nil {
		panic("evm caller cannot be nil")
	}
	return UnwrapERC20IBCModule{
		app:         app,
		erc20Keeper: erc20Keeper,
		evmCaller:   evmCaller,
	}
}

func (im UnwrapERC20IBCModule) OnRecvPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	ack := im.app.OnRecvPacket(ctx, channelVersion, packet, relayer)
	if !ack.Success() {
		return ack
	}

	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return ack
	}
	if !shouldUnwrapFromIBCMemo(data.Memo) {
		return ack
	}

	receiverHex, ok := receiverToHexAddress(data.Receiver)
	if !ok {
		return ack
	}

	token := transfertypes.Token{
		Denom:  transfertypes.ExtractDenomFromPath(data.Denom),
		Amount: data.Amount,
	}
	coin := evmibc.GetReceivedCoin(packet, token)

	// unwrap only for native ERC20 token-pairs (e.g. wmantraUSD wrapper)
	pairID := im.erc20Keeper.GetTokenPairID(ctx, coin.Denom)
	pair, found := im.erc20Keeper.GetTokenPair(ctx, pairID)
	if !found || !pair.Enabled || !pair.IsNativeERC20() {
		return ack
	}

	im.tryUnwrap(ctx, receiverHex, pair.GetERC20Contract(), coin.Amount.BigInt())
	return ack
}

func (im UnwrapERC20IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return im.app.OnAcknowledgementPacket(ctx, channelVersion, packet, acknowledgement, relayer)
}

func (im UnwrapERC20IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return im.app.OnTimeoutPacket(ctx, channelVersion, packet, relayer)
}

// UnmarshalPacketData implements the PacketDataUnmarshaler interface.
func (im UnwrapERC20IBCModule) UnmarshalPacketData(
	ctx sdk.Context,
	portID string,
	channelID string,
	bz []byte,
) (interface{}, string, error) {
	if unmarshaler, ok := im.app.(porttypes.PacketDataUnmarshaler); ok {
		return unmarshaler.UnmarshalPacketData(ctx, portID, channelID, bz)
	}
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return nil, "", err
	}
	return data, transfertypes.V1, nil
}

func (im UnwrapERC20IBCModule) OnChanOpenInit(
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

func (im UnwrapERC20IBCModule) OnChanOpenTry(
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

func (im UnwrapERC20IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	return im.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (im UnwrapERC20IBCModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanOpenConfirm(ctx, portID, channelID)
}

func (im UnwrapERC20IBCModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanCloseInit(ctx, portID, channelID)
}

func (im UnwrapERC20IBCModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanCloseConfirm(ctx, portID, channelID)
}

func receiverToHexAddress(receiver string) (common.Address, bool) {
	if common.IsHexAddress(receiver) {
		return common.HexToAddress(receiver), true
	}
	addr, err := sdk.AccAddressFromBech32(receiver)
	if err != nil {
		return common.Address{}, false
	}
	return common.BytesToAddress(addr.Bytes()), true
}

func shouldUnwrapFromIBCMemo(memo string) bool {
	if memo == "" {
		return false
	}
	var parsed struct {
		Mantra struct {
			Unwrap bool `json:"unwrap"`
		} `json:"mantra"`
	}
	if err := json.Unmarshal([]byte(memo), &parsed); err != nil {
		return false
	}
	return parsed.Mantra.Unwrap
}

func (im UnwrapERC20IBCModule) tryUnwrap(ctx sdk.Context, receiver common.Address, wrapper common.Address, amountWad *big.Int) {
	if amountWad == nil || amountWad.Sign() <= 0 {
		return
	}
	if amountWad.Cmp(evmutil.MinWithdrawAmountWad) < 0 {
		return
	}

	if _, err := im.evmGetUnderlying(ctx, receiver, wrapper); err != nil {
		return
	}

	_, _ = im.evmWithdrawTo(ctx, receiver, wrapper, receiver, amountWad)
}

func (im UnwrapERC20IBCModule) evmGetUnderlying(ctx sdk.Context, caller common.Address, wrapper common.Address) (common.Address, error) {
	data := make([]byte, 4)
	copy(data, evmutil.SelectorUnderlyingGetter[:])

	res, err := im.evmCaller.CallEVMWithData(ctx, caller, &wrapper, data, false, nil)
	if res == nil {
		return common.Address{}, err
	}
	if err != nil {
		return common.Address{}, err
	}
	if len(res.Ret) < 32 {
		return common.Address{}, fmt.Errorf("underlying() returned short data")
	}
	underlying := common.BytesToAddress(res.Ret[12:32])
	if underlying == (common.Address{}) {
		return common.Address{}, fmt.Errorf("underlying() returned zero address")
	}
	return underlying, nil
}

func (im UnwrapERC20IBCModule) evmWithdrawTo(ctx sdk.Context, caller common.Address, wrapper common.Address, to common.Address, amount *big.Int) ([]byte, error) {
	if amount == nil || amount.Sign() <= 0 {
		return nil, fmt.Errorf("invalid withdraw amount")
	}
	amtU256, overflow := uint256.FromBig(amount)
	if overflow {
		return nil, fmt.Errorf("withdraw amount overflows uint256")
	}

	data := make([]byte, 4+32+32)
	copy(data[:4], evmutil.SelectorWithdrawToAddressUint256[:])
	copy(data[4:4+32], common.LeftPadBytes(to.Bytes(), 32))
	amt32 := amtU256.Bytes32()
	copy(data[4+32:], amt32[:])

	res, err := im.evmCaller.CallEVMWithData(ctx, caller, &wrapper, data, true, nil)
	if res == nil {
		return nil, err
	}
	return res.Ret, err
}
