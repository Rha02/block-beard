package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"

	"github.com/Rha02/block-beard/src/utils"
)

// Transaction is a struct for a transaction.
type Transaction struct {
	senderPrivateKey *ecdsa.PrivateKey
	senderPublicKey  *ecdsa.PublicKey
	senderAddress    string
	recipientAddress string
	amount           float32
}

// NewTransaction creates a new transaction.
func NewTransaction(
	privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey,
	senderAddress, recipientAddress string, amount float32,
) *Transaction {
	return &Transaction{
		senderPrivateKey: privateKey,
		senderPublicKey:  publicKey,
		senderAddress:    senderAddress,
		recipientAddress: recipientAddress,
		amount:           amount,
	}
}

// MarshalJSON is a custom JSON marshaller for the Transaction struct.
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_address"`
		Recipient string  `json:"recipient_address"`
		Amount    float32 `json:"amount"`
	}{
		Sender:    t.senderAddress,
		Recipient: t.recipientAddress,
		Amount:    t.amount,
	})
}

// GenerateSignature() generates a signature for the transaction.
func (t *Transaction) GenerateSignature() *utils.Signature {
	m, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	h := sha256.Sum256(m)

	r, s, err := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	if err != nil {
		panic(err)
	}

	return &utils.Signature{R: r, S: s}
}
