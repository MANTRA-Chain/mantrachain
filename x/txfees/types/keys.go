package types

const (
	// ModuleName defines the module name
	ModuleName = "txfees"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var ParamsKey = []byte("p_txfees")

func KeyPrefix(p string) []byte {
	return []byte(p)
}
