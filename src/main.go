package main

import (
	"fmt"
	"log"

	"github.com/Rha02/block-beard/src/blockchain"
	"github.com/Rha02/block-beard/src/wallet"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	wm := wallet.NewWallet()
	w1 := wallet.NewWallet()
	w2 := wallet.NewWallet()

	t := wallet.NewTransaction(w1.GetPrivateKey(), w1.GetPublicKey(), w1.GetAddress(), w2.GetAddress(), 1.0)

	bc := blockchain.NewBlockchain(wm.GetAddress())

	success := bc.AddTransaction(w1.GetAddress(), w2.GetAddress(), 1.0, w1.GetPublicKey(), t.GenerateSignature())
	fmt.Println(success)

	bc.Mine()
	fmt.Println(bc.ToString())

	fmt.Println("Balance of w1:", bc.GetBalance(w1.GetAddress()))
	fmt.Println("Balance of w2:", bc.GetBalance(w2.GetAddress()))
	fmt.Println("Balance of wm:", bc.GetBalance(wm.GetAddress()))
}
