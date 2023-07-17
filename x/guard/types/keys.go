package types

const (
	// ModuleName defines the module name
	ModuleName = "guard"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_guard"
)

var (
	accountPrivilegesStoreKey  = "account-privileges-store"
	requiredPrivilegesStoreKey = "required-privileges-store"

	Delimiter   = []byte{0x00}
	Placeholder = []byte{0x01}
)

const (
	GuardTransferCoinsKey = "guard-transfer-coins-value"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func RequiredPrivilegesStoreKey(kind []byte) []byte {
	key := make([]byte, len(requiredPrivilegesStoreKey)+len(Delimiter)+len(kind)+len(Delimiter))
	copy(key, requiredPrivilegesStoreKey)
	copy(key[len(requiredPrivilegesStoreKey):], Delimiter)
	copy(key[len(requiredPrivilegesStoreKey)+len(Delimiter):], kind)
	copy(key[len(requiredPrivilegesStoreKey)+len(Delimiter)+len(kind):], Delimiter)

	return key
}

func AccountPrivilegesStoreKey() []byte {
	key := make([]byte, len(accountPrivilegesStoreKey)+len(Delimiter))
	copy(key, accountPrivilegesStoreKey)
	copy(key[len(accountPrivilegesStoreKey):], Delimiter)

	return key
}
