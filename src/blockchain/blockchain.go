package blockchain

import (
	"fmt"
)

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

// GetLastBlock() returns a pointer to the last block in the blockchain.
func (bc *Blockchain) GetLastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// ToString() returns a developer-friendly string representation of the blockchain.
func (bc *Blockchain) ToString() string {
	var res string
	for i, block := range bc.chain {
		res += fmt.Sprintf("Block %d: %s", i, block.ToString())
		if i < len(bc.chain)-1 {
			res += " | "
		}
	}
	return res
}
