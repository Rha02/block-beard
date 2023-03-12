package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
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

func (b *Block) GetPrevHash() [32]byte {
	return b.prevHash
}

func (b *Block) GetTimestamp() int64 {
	return b.timestamp
}

func (b *Block) GetTransactions() []*Transaction {
	return b.transactions
}

func (b *Block) GetNonce() int {
	return b.nonce
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

// UnmarshalJSON() takes a json representation of a block and returns a pointer to a new block.
func (b *Block) UnmarshalJSON(data []byte) error {
	var prevHash string

	tmp := &struct {
		PrevHash     *string         `json:"prevHash"`
		Timestamp    *int64          `json:"timestamp"`
		Transactions *[]*Transaction `json:"transactions"`
		Nonce        *int            `json:"nonce"`
	}{
		PrevHash:     &prevHash,
		Timestamp:    &b.timestamp,
		Transactions: &b.transactions,
		Nonce:        &b.nonce,
	}
	if err := json.Unmarshal(data, tmp); err != nil {
		return err
	}

	decodedPrevHash, _ := hex.DecodeString(prevHash)
	copy(b.prevHash[:], decodedPrevHash[:32])

	return nil
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
