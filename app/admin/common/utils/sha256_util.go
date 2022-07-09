package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(data string, salt string) string {
	s := sha256.New()
	s.Write([]byte(data + salt))
	return hex.EncodeToString(s.Sum(nil))
}
