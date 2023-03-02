package main

import (
	"fmt"
	"log"

	"github.com/Rha02/block-beard/src/wallet"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	w := wallet.NewWallet("Rha02")
	fmt.Println(w.GetPrivateKeyStr())
	fmt.Println(w.GetPublicKeyStr())
	fmt.Println(w.GetAddress())
}
