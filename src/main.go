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
	bc := blockchain.NewBlockchain("miner_address")
	bc.AddTransaction("Bruce Wayne", "Clark Kent", 100)
	bc.Mine()

	bc.AddTransaction("Clark Kent", "Bruce Wayne", 70)
	bc.AddTransaction("Bruce Wayne", "Clark Kent", 50)
	bc.Mine()

	fmt.Println(bc.ToString())

	fmt.Printf("Balance of Bruce Wayne: %f\n", bc.GetBalance("Bruce Wayne"))
	fmt.Printf("Balance of Clark Kent: %f\n", bc.GetBalance("Clark Kent"))
	fmt.Printf("Balance of miner: %f", bc.GetBalance("miner_address"))
}
