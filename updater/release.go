package updater

import (
	"crypto"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/inconshreveable/go-update"
)

type ReleaseVersion struct {
	name         string
	tag          string
	binaryURL    string
	checksumURL  string
	signatureURL string
}

func (rv *ReleaseVersion) Name() string {
	return rv.name
}

func (rv *ReleaseVersion) Tag() string {
	return rv.tag
}

func (rv *ReleaseVersion) Apply(pubkey string) error {

	// get hash bytes
	checksumResp, err := http.Get(rv.checksumURL)
	if err != nil {
		return err
	}
	defer checksumResp.Body.Close()
	checksumBytes, err := ioutil.ReadAll(checksumResp.Body)
	if err != nil {
		return err
	}
	checksum, err := hex.DecodeString(string(checksumBytes))
	if err != nil {
		return err
	}

	// get signature bytes
	signatureResp, err := http.Get(rv.checksumURL)
	if err != nil {
		return err
	}
	defer signatureResp.Body.Close()
	signatureBytes, err := ioutil.ReadAll(signatureResp.Body)
	if err != nil {
		return err
	}
	signature, err := hex.DecodeString(string(signatureBytes))
	if err != nil {
		return err
	}

	// get binary body
	binaryResp, err := http.Get(rv.binaryURL)
	if err != nil {
		return err
	}
	defer binaryResp.Body.Close()

	// set options
	opts := update.Options{
		Hash:      crypto.SHA256,
		Signature: signature,
		Checksum:  checksum,
		Verifier:  NewED25519Verifier(),
	}
	err = opts.SetPublicKeyPEM([]byte(`
	-----BEGIN PUBLIC KEY-----
	` + pubkey + `
	-----END PUBLIC KEY-----
	`))
	if err != nil {
		return err
	}

	// apply update
	err = update.Apply(binaryResp.Body, opts)
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			return errors.New("Failed to rollback from bad update: " + rerr.Error())
		}
	}
	return err

}
