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

	bc.AddBlock(0, bc.GetLastBlock().Hash())

	bc.AddBlock(1, bc.GetLastBlock().Hash())

	fmt.Println(bc.ToString())
}
