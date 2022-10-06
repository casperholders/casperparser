package utils

import (
	"encoding/hex"
	"golang.org/x/crypto/blake2b"
	"strings"
)

const (
	StrKeyTagEd25519        = "ed25519"
	StrKeyTagSecp256k1      = "secp256k1"
	Separator          byte = 0
)

func AccountHash(publicKey string) string {
	tag := ""
	strippedKey := ""

	if strings.HasPrefix(publicKey, "01") {
		tag = StrKeyTagEd25519
		strippedKey = strings.TrimPrefix(publicKey, "01")
	} else if strings.HasPrefix(publicKey, "02") {
		tag = StrKeyTagSecp256k1
		strippedKey = strings.TrimPrefix(publicKey, "02")
	} else {
		return ""
	}
	publicKeyBytes, err := hex.DecodeString(strippedKey)

	if err != nil {
		return ""
	}

	buffer := append([]byte(tag), Separator)
	buffer = append(buffer, publicKeyBytes...)

	hash := blake2b.Sum256(buffer)

	return hex.EncodeToString(hash[:])
}
