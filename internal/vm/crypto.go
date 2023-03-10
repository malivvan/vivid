package vm

import (
	"github.com/malivvan/vlang/internal/crypto"
)

var scryptLevel = 0

func encrypt(password string, data []byte) []byte {
	encryptedData, err := crypto.Encrypt(crypto.ScryptLevel[0], password, data)
	if err != nil {
		panic(err)
	}
	return encryptedData
}

func decrypt(password string, data []byte) []byte {
	decryptedData, _, err := crypto.Decrypt(password, data)
	if err != nil {
		panic(err)
	}
	return decryptedData
}
