package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Rha02/block-beard/src/utils"
)

const (
	MiningDifficulty = 2
	MINING_SENDER    = "BlockBeard"
	MINING_REWARD    = 1.0
)

// Blockchain is a struct for the blockchain.
type Blockchain struct {
	pool    []*Transaction
	chain   []*Block
	address string
	port    uint16
}

// NewBlockchain() returns a pointer to a new blockchain
func NewBlockchain(bcAddress string, port uint16) *Blockchain {
	b := new(Block)
	bc := new(Blockchain)
	bc.AddBlock(0, b.Hash())
	bc.address = bcAddress
	return bc
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Chain   []*Block
		Pool    []*Transaction
		Address string
		Port    uint16
	}{
		Chain:   bc.chain,
		Pool:    bc.pool,
		Address: bc.address,
		Port:    bc.port,
	})
}

func (bc *Blockchain) GetTransactions() []*Transaction {
	return bc.pool
}

// AddBlock() takes a nonce and a previous hash and adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(nonce int, prevHash [32]byte) *Block {
	block := NewBlock(nonce, prevHash, bc.pool)
	bc.pool = []*Transaction{}
	bc.chain = append(bc.chain, block)
	return block
}

func (bc *Blockchain) CreateTransaction(
	sender, recipient string, amount float32, senderPublicKey *ecdsa.PublicKey, signature *utils.Signature,
) bool {
	isTransacted := bc.AddTransaction(sender, recipient, amount, senderPublicKey, signature)
	return isTransacted
}

// AddTransaction() creates a transaction and adds it to the pool.
func (bc *Blockchain) AddTransaction(
	sender, recipient string, amount float32, senderPublicKey *ecdsa.PublicKey, signature *utils.Signature,
) bool {
	t := NewTransaction(sender, recipient, amount)

	if sender != MINING_SENDER && !bc.VerifyTransaction(senderPublicKey, signature, t) {
		fmt.Printf("Invalid transaction from %s", sender)
		return false
	}

	// if bc.GetBalance(sender) < amount {
	// 	fmt.Printf("Not enough funds to send %f from %s to %s\n", amount, sender, recipient)
	// 	return false
	// }
	bc.pool = append(bc.pool, t)
	return true
}

// VerifyTransaction() takes a public key, a signature, and a transaction and returns whether the transaction is valid.
func (bc *Blockchain) VerifyTransaction(
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction,
) bool {
	m, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	h := sha256.Sum256(m)
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
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

// Mine() mines a new block.
func (bc *Blockchain) Mine() {
	nonce := bc.ProofOfWork()
	bc.AddTransaction(MINING_SENDER, bc.address, MINING_REWARD, nil, nil)
	bc.AddBlock(nonce, bc.GetLastBlock().Hash())
	fmt.Println("Mined a new block successfully!")
}

// GetBalance() returns the balance of a given address.
func (bc *Blockchain) GetBalance(address string) float32 {
	var balance float32
	for _, block := range bc.chain {
		for _, transaction := range block.transactions {
			if transaction.senderAddress == address {
				balance -= transaction.amount
			}
			if transaction.recipientAddress == address {
				balance += transaction.amount
			}
		}
	}
	return balance
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
