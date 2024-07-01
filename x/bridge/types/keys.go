package types

const (
	// ModuleName defines the module name
	ModuleName = "bridge"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	ParamsKey = []byte("p_bridge")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
