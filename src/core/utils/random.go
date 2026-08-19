package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHex(n uint) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
