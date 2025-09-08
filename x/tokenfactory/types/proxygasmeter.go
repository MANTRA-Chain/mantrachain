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
func NewProxyGasMeter(gasMeter storetypes.GasMeter, limit storetypes.Gas) storetypes.GasMeter {
	gasLimit := min(limit, gasMeter.GasRemaining())
	base := storetypes.NewGasMeter(gasLimit)

	if limit >= gasMeter.GasRemaining() {
		return gasMeter
	}

	return &ProxyGasMeter{
		GasMeter: base,
		parent:   gasMeter,
	}
}

func (pgm ProxyGasMeter) ConsumeGas(amount storetypes.Gas, descriptor string) {
	pgm.GasMeter.ConsumeGas(amount, descriptor)
	pgm.parent.ConsumeGas(amount, descriptor)
}

func (pgm ProxyGasMeter) String() string {
	return fmt.Sprintf("ProxyGasMeter{consumed: %d, limit: %d}", pgm.GasConsumed(), pgm.Limit)
}
