package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Block is a struct for a block in the blockchain.
type Block struct {
	prevHash     [32]byte
	timestamp    int64
	transactions []*Transaction
	nonce        int
}

// NewBlock() takes a nonce and a previous hash and returns a pointer to a new block.
func NewBlock(nonce int, prevHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		prevHash:     prevHash,
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		transactions: transactions,
	}
}

// Hash() returns the hash of the block.
func (b *Block) Hash() [32]byte {
	// marshall the block to json
	res, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}

	return sha256.Sum256(res)
}

// MarshalJSON() returns a json representation of the block.
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrevHash     string         `json:"prevHash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
		Nonce        int            `json:"nonce"`
	}{
		PrevHash:     fmt.Sprintf("%x", b.prevHash),
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
		Nonce:        b.nonce,
	})
}

// ToString() returns a developer-friendly string representation of the block.
func (b *Block) ToString() string {
	var transactions string
	for _, t := range b.transactions {
		transactions += t.ToString() + ", "
	}

	return fmt.Sprintf(
		"Prev. hash: %x, Timestamp: %d, Transactions: {%s}, Nonce: %d",
		b.prevHash, b.timestamp, transactions, b.nonce,
	)
}
