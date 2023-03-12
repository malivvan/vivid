package main

import (
	"bufio"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/malivvan/vivid/internal/crypto"
)

func main() {

	// Generate key
	var (
		err   error
		b     []byte
		block *pem.Block
		pub   ed25519.PublicKey
		priv  ed25519.PrivateKey
	)
	pub, priv, err = ed25519.GenerateKey(rand.Reader)
	check(err)

	// Encode
	b, err = x509.MarshalPKCS8PrivateKey(priv)
	check(err)
	block = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	}
	privBytes := pem.EncodeToMemory(block)
	b, err = x509.MarshalPKIXPublicKey(pub)
	check(err)
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}
	pubBytes := pem.EncodeToMemory(block)

	// Encrypt
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter password: ")
	scanner.Scan()
	pass1 := scanner.Text()
	fmt.Scanln(pass1)
	fmt.Print("ReEnter password: ")
	scanner.Scan()
	pass2 := scanner.Text()
	fmt.Scanln(pass2)
	if pass1 != pass2 {
		log.Fatal("passwords do not match")
	}
	encPrivBytes, err := crypto.Encrypt(crypto.ScryptLevel[0], pass1, privBytes)
	check(err)

	// Write
	check(ioutil.WriteFile("../signer/key.pem", encPrivBytes, 0600))
	check(ioutil.WriteFile("../../cmd/pub.pem", pubBytes, 0644))
	fmt.Println(hex.EncodeToString([]byte(pub)))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
