package main

import (
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	bc := NewBlockchain()

	bc.AddBlock(0, bc.GetLastBlock().Hash())

	bc.AddBlock(1, bc.GetLastBlock().Hash())

	fmt.Println(bc.toString())
}
