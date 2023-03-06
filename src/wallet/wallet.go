package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// Wallet is a struct for a wallet.
type Wallet struct {
	address    string
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// NewWallet() returns a pointer to a new wallet.
func NewWallet() *Wallet {
	// Create a new wallet
	w := new(Wallet)

	// Generate a new private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	w.privateKey = privateKey

	// Set the public key
	w.publicKey = &privateKey.PublicKey

	// Generate SHA256 hash of the public key
	h := sha256.New()
	h.Write(w.publicKey.X.Bytes())
	h.Write(w.publicKey.Y.Bytes())
	digest := h.Sum(nil)

	// Generate RIPEMD160 hash of the SHA256 hash
	h2 := ripemd160.New()
	h2.Write(digest)
	digest2 := h2.Sum(nil)

	// Add the network byte to the beginning of the hash
	digest3 := append([]byte{0x00}, digest2...)

	// Generate SHA256 hash of the hash
	h3 := sha256.New()
	h3.Write(digest3)
	digest4 := h3.Sum(nil)

	// Generate SHA256 hash of the hash
	h4 := sha256.New()
	h4.Write(digest4)
	digest5 := h4.Sum(nil)

	// Take the first 4 bytes of the hash and append it to the end of the hash
	digest6 := append(digest3, digest5[:4]...)

	// Convert the hash to a base58 string
	w.address = base58.Encode(digest6)

	return w
}

// MarshalJSON() returns the JSON representation of the wallet.
func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
		Address    string `json:"address"`
	}{
		PrivateKey: w.GetPrivateKeyStr(),
		PublicKey:  w.GetPublicKeyStr(),
		Address:    w.GetAddress(),
	})
}

// GetAddress() returns the address of the wallet.
func (w *Wallet) GetAddress() string {
	return w.address
}

// GetPrivateKey() returns the private key of the wallet.
func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

// GetPrivateKeyStr() returns the private key of the wallet as a string.
func (w *Wallet) GetPrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

// GetPublicKey() returns the public key of the wallet.
func (w *Wallet) GetPublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

// GetPublicKeyStr() returns the public key of the wallet as a string.
func (w *Wallet) GetPublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}
