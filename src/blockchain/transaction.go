package blockchain

import (
	"encoding/json"
	"fmt"
)

// Transaction is a struct for a transaction in the blockchain.
type Transaction struct {
	senderAddress    string
	recipientAddress string
	amount           float32
}

// NewTransaction() takes a sender, recipient, and amount and returns a pointer to a new transaction.
func NewTransaction(sender, recipient string, amount float32) *Transaction {
	return &Transaction{
		senderAddress:    sender,
		recipientAddress: recipient,
		amount:           amount,
	}
}

// ToString() returns a developer-friendly string representation of the transaction.
func (t *Transaction) ToString() string {
	return fmt.Sprintf(
		"Sender: %s, Recipient: %s, Amount: %f",
		t.senderAddress, t.recipientAddress, t.amount,
	)
}

// MarshalJSON() returns a json representation of the transaction.
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string  `json:"sender_address"`
		RecipientAddress string  `json:"recipient_address"`
		Amount           float32 `json:"amount"`
	}{
		SenderAddress:    t.senderAddress,
		RecipientAddress: t.recipientAddress,
		Amount:           t.amount,
	})
}

type TransactionRequest struct {
	SenderAddress    *string  `json:"sender_address"`
	RecipientAddress *string  `json:"recipient_address"`
	Amount           *float32 `json:"amount"`
	SenderPublicKey  *string  `json:"sender_public_key"`
	Signature        *string  `json:"signature"`
}

func (tr *TransactionRequest) IsValid() bool {
	return tr.SenderAddress != nil && tr.RecipientAddress != nil && tr.Amount != nil && tr.SenderPublicKey != nil && tr.Signature != nil
}
