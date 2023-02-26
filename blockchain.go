package main

import (
	"fmt"
	"time"
)

// Block is a struct for a block in the blockchain.
type Block struct {
	prevHash     []byte
	timestamp    int64
	transactions []string
	nonce        int
}

// NewBlock() takes a nonce and a previous hash and returns a pointer to a new block.
func NewBlock(nonce int, prevHash []byte) *Block {
	return &Block{
		prevHash:  prevHash,
		timestamp: time.Now().UnixNano(),
		nonce:     nonce,
	}
}

// ToString() returns a developer-friendly string representation of the block.
func (b *Block) toString() string {
	return fmt.Sprintf(
		"Prev. hash: %s, Timestamp: %d, Transactions: %s, Nonce: %d",
		string(b.prevHash), b.timestamp, b.transactions, b.nonce,
	)
}

// Blockchain is a struct for the blockchain.
type Blockchain struct {
	chain []*Block
}

// NewBlockchain() returns a pointer to a new blockchain
func NewBlockchain() *Blockchain {
	return &Blockchain{
		chain: []*Block{NewBlock(0, []byte("First Block"))},
	}
}

// AddBlock() takes a nonce and a previous hash and adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(nonce int, prevHash []byte) *Block {
	block := NewBlock(nonce, prevHash)
	bc.chain = append(bc.chain, block)
	return block
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
