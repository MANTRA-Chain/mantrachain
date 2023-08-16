package amm

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ Order   = (*BaseOrder)(nil)
	_ Orderer = (*BaseOrderer)(nil)

	DefaultOrderer = BaseOrderer{}
)

// OrderDirection specifies an order direction, either buy or sell.
type OrderDirection int

// OrderDirection enumerations.
const (
	Buy OrderDirection = iota + 1
	Sell
)

func (dir OrderDirection) String() string {
	switch dir {
	case Buy:
		return "Buy"
	case Sell:
		return "Sell"
	default:
		return fmt.Sprintf("OrderDirection(%d)", dir)
	}
}

type Orderer interface {
	Order(dir OrderDirection, price sdk.Dec, amt math.Int) Order
}

// BaseOrderer creates new BaseOrder with sufficient offer coin amount
// considering price and amount.
type BaseOrderer struct{}

func (orderer BaseOrderer) Order(dir OrderDirection, price sdk.Dec, amt math.Int) Order {
	return NewBaseOrder(dir, price, amt, OfferCoinAmount(dir, price, amt))
}

// Order is the universal interface of an order.
type Order interface {
	GetDirection() OrderDirection
	// GetBatchId returns the batch id where the order was created.
	// Batch id of 0 means the current batch.
	GetBatchId() uint64
	GetPrice() sdk.Dec
	GetAmount() math.Int // The original order amount
	GetOfferCoinAmount() math.Int
	GetPaidOfferCoinAmount() math.Int
	SetPaidOfferCoinAmount(amt math.Int)
	GetReceivedDemandCoinAmount() math.Int
	SetReceivedDemandCoinAmount(amt math.Int)
	GetOpenAmount() math.Int
	SetOpenAmount(amt math.Int)
	IsMatched() bool
	// HasPriority returns true if the order has higher priority
	// than the other order.
	HasPriority(other Order) bool
	String() string
}

// BaseOrder is the base struct for an Order.
type BaseOrder struct {
	Direction       OrderDirection
	Price           sdk.Dec
	Amount          math.Int
	OfferCoinAmount math.Int

	// Match info
	OpenAmount               math.Int
	PaidOfferCoinAmount      math.Int
	ReceivedDemandCoinAmount math.Int
}

// NewBaseOrder returns a new BaseOrder.
func NewBaseOrder(dir OrderDirection, price sdk.Dec, amt, offerCoinAmt math.Int) *BaseOrder {
	return &BaseOrder{
		Direction:                dir,
		Price:                    price,
		Amount:                   amt,
		OfferCoinAmount:          offerCoinAmt,
		OpenAmount:               amt,
		PaidOfferCoinAmount:      sdk.ZeroInt(),
		ReceivedDemandCoinAmount: sdk.ZeroInt(),
	}
}

// GetDirection returns the order direction.
func (order *BaseOrder) GetDirection() OrderDirection {
	return order.Direction
}

func (order *BaseOrder) GetBatchId() uint64 {
	return 0
}

// GetPrice returns the order price.
func (order *BaseOrder) GetPrice() sdk.Dec {
	return order.Price
}

// GetAmount returns the order amount.
func (order *BaseOrder) GetAmount() math.Int {
	return order.Amount
}

func (order *BaseOrder) GetOfferCoinAmount() math.Int {
	return order.OfferCoinAmount
}

func (order *BaseOrder) GetPaidOfferCoinAmount() math.Int {
	return order.PaidOfferCoinAmount
}

func (order *BaseOrder) SetPaidOfferCoinAmount(amt math.Int) {
	order.PaidOfferCoinAmount = amt
}

func (order *BaseOrder) GetReceivedDemandCoinAmount() math.Int {
	return order.ReceivedDemandCoinAmount
}

func (order *BaseOrder) SetReceivedDemandCoinAmount(amt math.Int) {
	order.ReceivedDemandCoinAmount = amt
}

func (order *BaseOrder) GetOpenAmount() math.Int {
	return order.OpenAmount
}

func (order *BaseOrder) SetOpenAmount(amt math.Int) {
	order.OpenAmount = amt
}

func (order *BaseOrder) IsMatched() bool {
	return order.OpenAmount.LT(order.Amount)
}

// HasPriority returns whether the order has higher priority than
// the other order.
func (order *BaseOrder) HasPriority(other Order) bool {
	return order.Amount.GT(other.GetAmount())
}

func (order *BaseOrder) String() string {
	return fmt.Sprintf("BaseOrder(%s,%s,%s)", order.Direction, order.Price, order.Amount)
}
