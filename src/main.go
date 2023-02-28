package main

import (
	"fmt"
	"log"

	"github.com/Rha02/block-beard/src/blockchain"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	bc := blockchain.NewBlockchain()
	bc.AddTransaction("Bruce Wayne", "Clark Kent", 100)

	nonce := bc.ProofOfWork()
	bc.AddBlock(nonce, bc.GetLastBlock().Hash())

	nonce = bc.ProofOfWork()
	bc.AddBlock(nonce, bc.GetLastBlock().Hash())

	fmt.Println(bc.ToString())
}
