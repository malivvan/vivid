package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/malivvan/vivid/internal/crypto"
)

var DistPath = ".." + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "dist" + string(os.PathSeparator)

func main() {
	// Read key
	keyPEMEnc, err := ioutil.ReadFile("key.pem")
	check(err)

	// Read password
	fmt.Printf("Password: ")
	pass, err := gopass.GetPasswd()
	check(err)

	// Decrypt
	keyPEM, _, err := crypto.Decrypt(string(pass), keyPEMEnc)
	check(err)

	// Decode
	p, _ := pem.Decode(keyPEM)
	if p == nil {
		log.Fatal("no pem block found")
	}
	key, err := x509.ParsePKCS8PrivateKey(p.Bytes)
	check(err)
	edKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		log.Fatal("key is not ed25519 key")
	}
	fmt.Println(hex.EncodeToString([]byte(edKey.Public().(ed25519.PublicKey))))

	// Sign and hash
	check(filepath.Walk(DistPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".sha256" || ext == ".sig" {
			return nil
		}

		hash(path)
		sign(path, edKey)

		return nil
	}))
}

func hash(path string) {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	check(err)
	hasher := sha256.New()
	hasher.Write(buf)
	hashHex := hex.EncodeToString(hasher.Sum(nil))
	check(ioutil.WriteFile(path+".sha256", []byte(hashHex), 0644))
	fmt.Println("hashed", strings.TrimPrefix(path, DistPath))
}

func sign(path string, key ed25519.PrivateKey) {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	check(err)
	signature := ed25519.Sign(key, buf)
	signatureHex := hex.EncodeToString(signature)
	check(ioutil.WriteFile(path+".sig", []byte(signatureHex), 0644))
	fmt.Println("signed", strings.TrimPrefix(path, DistPath))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
