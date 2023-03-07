package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rha02/block-beard/src/blockchain"
	"github.com/Rha02/block-beard/src/utils"
	"github.com/Rha02/block-beard/src/wallet"
)

var cache = make(map[string]*blockchain.Blockchain)

type Server struct {
	port uint16
}

func NewServer(port uint16) *Server {
	return &Server{port}
}

func (s *Server) Port() uint16 {
	return s.port
}

func (s *Server) Start() {
	http.HandleFunc("/", s.GetChainHandler)
	http.HandleFunc("/transaction", s.TransactionsHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *Server) GetBlockchain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minerWallet := wallet.NewWallet()
		bc = blockchain.NewBlockchain(minerWallet.GetAddress(), s.Port())
		cache["blockchain"] = bc
	}
	return bc
}

func (s *Server) GetChainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	bc := s.GetBlockchain()
	m, _ := bc.MarshalJSON()

	w.Header().Set("Content-Type", "application/json")
	w.Write(m)
}

func (s *Server) TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var t blockchain.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Error decoding transaction request:", err)
			return
		}
		if !t.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("Invalid transaction request: missing fields")
			return
		}
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)
		bc := s.GetBlockchain()

		w.Header().Set("Content-Type", "application/json")

		if !bc.CreateTransaction(*t.SenderAddress, *t.RecipientAddress, *t.Amount, publicKey, signature) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.JsonStatus("Transaction failed"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(utils.JsonStatus("Transaction successful"))
		return
	}

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		bc := s.GetBlockchain()
		transactions := bc.GetTransactions()
		m, _ := json.Marshal(struct {
			Transactions []*blockchain.Transaction
		}{
			Transactions: transactions,
		})
		w.Write(m)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
