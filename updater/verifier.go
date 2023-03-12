package updater

import (
	"crypto"
	"crypto/ed25519"
	"errors"

	"github.com/inconshreveable/go-update"
)

func NewED25519Verifier() update.Verifier {
	return verifyFn(func(checksum, signature []byte, hash crypto.Hash, publicKey crypto.PublicKey) error {
		key, ok := publicKey.(ed25519.PublicKey)
		if !ok {
			return errors.New("not a valid ECDSA public key")
		}
		ok = ed25519.Verify(key, checksum, signature)
		if !ok {
			return errors.New("failed to verify ed25519 signature")
		}
		return nil
	})
}

type verifyFn func([]byte, []byte, crypto.Hash, crypto.PublicKey) error

func (fn verifyFn) VerifySignature(checksum []byte, signature []byte, hash crypto.Hash, publicKey crypto.PublicKey) error {
	return fn(checksum, signature, hash, publicKey)
}
