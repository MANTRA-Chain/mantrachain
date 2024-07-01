package conv

import (
	"encoding/binary"
	"fmt"
)

type KeyPrefixBuilder struct {
	Key []byte
}

func (k *KeyPrefixBuilder) AInt(n uint64) {
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, n)
	k.Key = append(k.Key, indexBytes...)
}

func (k *KeyPrefixBuilder) AString(s string) {
	k.Key = append(k.Key, []byte(s)...)
}

func GetByteKey(keys ...interface{}) []byte {
	builder := KeyPrefixBuilder{Key: []byte{}}
	for _, key := range keys {
		switch v := key.(type) {
		default:
			panic(fmt.Sprintf("Unsupported Key Type: %T with value: %#v", v, key))
		case uint64:
			builder.AInt(v)
		case string:
			builder.AString(v)
		case []byte:
			builder.Key = append(builder.Key, v...)
		}
	}
	return builder.Key
}
