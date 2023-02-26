package main

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
	transactions []string
	nonce        int
}

// NewBlock() takes a nonce and a previous hash and returns a pointer to a new block.
func NewBlock(nonce int, prevHash [32]byte) *Block {
	return &Block{
		prevHash:  prevHash,
		timestamp: time.Now().UnixNano(),
		nonce:     nonce,
	}
}

func (b *Block) Hash() [32]byte {
	// marshall the block to json
	res, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}

	return sha256.Sum256(res)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrevHash     [32]byte `json:"prevHash"`
		Timestamp    int64    `json:"timestamp"`
		Transactions []string `json:"transactions"`
		Nonce        int      `json:"nonce"`
	}{
		PrevHash:     b.prevHash,
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
		Nonce:        b.nonce,
	})
}

// ToString() returns a developer-friendly string representation of the block.
func (b *Block) toString() string {
	return fmt.Sprintf(
		"Prev. hash: %x, Timestamp: %d, Transactions: %s, Nonce: %d",
		b.prevHash, b.timestamp, b.transactions, b.nonce,
	)
}

// Blockchain is a struct for the blockchain.
type Blockchain struct {
	chain []*Block
}

// NewBlockchain() returns a pointer to a new blockchain
func NewBlockchain() *Blockchain {
	b := new(Block)
	bc := new(Blockchain)
	bc.AddBlock(0, b.Hash())
	return bc
}

// AddBlock() takes a nonce and a previous hash and adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(nonce int, prevHash [32]byte) *Block {
	block := NewBlock(nonce, prevHash)
	bc.chain = append(bc.chain, block)
	return block
}

func (bc *Blockchain) GetLastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// ToString() returns a developer-friendly string representation of the blockchain.
func (bc *Blockchain) toString() string {
	var res string
	for i, block := range bc.chain {
		res += fmt.Sprintf("Block %d: %s", i, block.toString())
		if i < len(bc.chain)-1 {
			res += " | "
		}
	}
	return res
}
