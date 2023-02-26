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
	fmt.Println(bc.toString())

	bc.AddBlock(1, []byte("Second Block"))

	fmt.Println(bc.toString())

	bc.AddBlock(2, []byte("Third Block"))

	fmt.Println(bc.toString())
}
