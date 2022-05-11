package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	types "github.com/LimeChain/mantrachain/x/mns/types"
)

func GetDomainIndex(domain string) string {
	index := sha256.Sum256([]byte(fmt.Sprintf("%s/domain/%s", types.ModuleName, domain)))
	return hex.EncodeToString(index[:])
}

func GetDomainNameIndex(domain string, domainName string) string {
	index := sha256.Sum256([]byte(fmt.Sprintf("%s/domain/%s/domainName/%s", types.ModuleName, domain, domainName)))
	return hex.EncodeToString(index[:])
}
