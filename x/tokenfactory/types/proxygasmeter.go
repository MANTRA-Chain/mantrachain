package types

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
)

var _ storetypes.GasMeter = &ProxyGasMeter{}

// ProxyGasMeter wraps another GasMeter, but enforces a lower gas limit.
// Gas consumption is delegated to the wrapped GasMeter, so it won't risk losing gas accounting compared to standalone
// gas meter.
type ProxyGasMeter struct {
	storetypes.GasMeter

	parent storetypes.GasMeter
}

// NewProxyGasMeter returns a new ProxyGasMeter which is like a basic gas meter with minimum of new limit and remaining gas
// of the parent gas meter, it also delegate the gas consumption to parent gas meter in real time, so it won't risk
// losing gas accounting in face of panics or other unexpected errors.
//
// If limit is greater than or equal to the remaining gas, no wrapping is needed and the original gas meter is returned.
func NewProxyGasMeter(gasMeter storetypes.GasMeter, limit storetypes.Gas) storetypes.GasMeter {
	if gasMeter == nil {
		panic("cannot create ProxyGasMeter with nil parent gas meter")
	}

	if limit >= gasMeter.GasRemaining() {
		return gasMeter
	}

	base := storetypes.NewGasMeter(limit)
	return &ProxyGasMeter{
		GasMeter: base,
		parent:   gasMeter,
	}
}

func (pgm ProxyGasMeter) RefundGas(amount storetypes.Gas, descriptor string) {
	pgm.parent.RefundGas(amount, descriptor)
	pgm.GasMeter.RefundGas(amount, descriptor)
}

func (pgm ProxyGasMeter) ConsumeGas(amount storetypes.Gas, descriptor string) {
	pgm.parent.ConsumeGas(amount, descriptor)
	pgm.GasMeter.ConsumeGas(amount, descriptor)
}

func (pgm ProxyGasMeter) String() string {
	return fmt.Sprintf("ProxyGasMeter{consumed: %d, limit: %d}", pgm.GasConsumed(), pgm.Limit())
}
