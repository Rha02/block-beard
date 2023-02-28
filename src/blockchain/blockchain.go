package blockchain

import (
	"fmt"
	"strings"
)

var MiningDifficulty = 2

// Blockchain is a struct for the blockchain.
type Blockchain struct {
	pool  []*Transaction
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
	block := NewBlock(nonce, prevHash, bc.pool)
	bc.pool = []*Transaction{}
	bc.chain = append(bc.chain, block)
	return block
}

// AddTransaction() creates a transaction and adds it to the pool.
func (bc *Blockchain) AddTransaction(sender, recipient string, amount float32) {
	t := NewTransaction(sender, recipient, amount)
	bc.pool = append(bc.pool, t)
}

// GetLastBlock() returns a pointer to the last block in the blockchain.
func (bc *Blockchain) GetLastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// CopyPool() returns a copy of the transaction pool.
func (bc *Blockchain) CopyPool() []*Transaction {
	var res []*Transaction
	for _, t := range bc.pool {
		res = append(res, NewTransaction(
			t.senderAddress,
			t.recipientAddress,
			t.amount,
		))
	}
	return res
}

// ValidateProof() takes a nonce, a previous hash, a list of transactions, and a difficulty and returns whether the proof is valid.
func (bc *Blockchain) ValidateProof(nonce int, prevHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeroes := strings.Repeat("0", difficulty)
	b := Block{
		prevHash:     prevHash,
		timestamp:    0,
		transactions: transactions,
		nonce:        nonce,
	}
	hashStr := fmt.Sprintf("%x", b.Hash())
	return hashStr[:difficulty] == zeroes
}

// ProofOfWork() returns a valid nonce for the current pool of transactions.
func (bc *Blockchain) ProofOfWork() int {
	nonce := 0
	for !bc.ValidateProof(nonce, bc.GetLastBlock().Hash(), bc.CopyPool(), MiningDifficulty) {
		nonce++
	}
	return nonce
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
