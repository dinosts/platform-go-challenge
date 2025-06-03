package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func NewHasher(salt string) func(str string) string {
	return func(str string) string {
		hasher := sha256.New()
		hasher.Write([]byte(salt + str))
		hash := hex.EncodeToString(hasher.Sum(nil))

		return hash
	}
}
