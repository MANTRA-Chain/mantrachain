package types

const (
	// ModuleName defines the module name
	ModuleName = "guard"

	// RouterKey is the message route for slashing
	RouterKey = ModuleName
)

var (
	ParamsKey = []byte("p_guard")

	accountPrivilegesStoreKey          = "account-privileges-store"
	requiredPrivilegesStoreKey         = "required-privileges-store"
	whitelistTransfersAccAddrsStoreKey = "whitelist-transfers-acc-addrs-store"

	whitelistTransfersAccAddrsIndex = "whitelist-transfers-acc-addrs-id"

	Delimiter   = []byte{0x00}
	Placeholder = []byte{0x01}
)

const (
	GuardTransferCoinsKey = "guard-transfer-coins-value"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetWhitelistTransfersAccAddrsIndex(id []byte) []byte {
	key := make([]byte, len(whitelistTransfersAccAddrsIndex)+len(Delimiter)+len(id)+len(Delimiter))
	copy(key, whitelistTransfersAccAddrsIndex)
	copy(key[len(whitelistTransfersAccAddrsIndex):], Delimiter)
	copy(key[len(whitelistTransfersAccAddrsIndex)+len(Delimiter):], id)
	copy(key[len(whitelistTransfersAccAddrsIndex)+len(Delimiter)+len(id):], Delimiter)
	return key
}

func WhitelistTransfersAccAddrsStoreKey() []byte {
	key := make([]byte, len(whitelistTransfersAccAddrsStoreKey)+len(Delimiter))
	copy(key, whitelistTransfersAccAddrsStoreKey)
	copy(key[len(whitelistTransfersAccAddrsStoreKey):], Delimiter)

	return key
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
