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
	s.GetBlockchain().Run()
	http.HandleFunc("/", s.GetChainHandler)
	http.HandleFunc("/transactions", s.TransactionsHandler)
	http.HandleFunc("/mine", s.MineHandler)
	http.HandleFunc("/mine/start", s.StartMineHandler)
	http.HandleFunc("/amount", s.AmountHandler)
	http.HandleFunc("/consensus", s.ConsensusHandler)
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

	if r.Method == http.MethodPut {
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

		if !bc.AddTransaction(*t.SenderAddress, *t.RecipientAddress, *t.Amount, publicKey, signature) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.JsonStatus("Transaction failed"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(utils.JsonStatus("Transaction successful"))
		return
	}

	if r.Method == http.MethodDelete {
		bc := s.GetBlockchain()
		bc.ClearTransactionsPool()
		w.WriteHeader(http.StatusOK)
		w.Write(utils.JsonStatus("Transactions pool cleared"))
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

func (s *Server) MineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	bc := s.GetBlockchain()
	if !bc.Mine() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JsonStatus("Mining failed"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(utils.JsonStatus("Mining successful"))
}

func (s *Server) StartMineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	bc := s.GetBlockchain()
	bc.StartMining()

	w.WriteHeader(http.StatusOK)
	w.Write(utils.JsonStatus("Mining successful"))
}

func (s *Server) AmountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	bcAddress := r.URL.Query().Get("blockchain_address")

	bc := s.GetBlockchain()
	amount := bc.GetBalance(bcAddress)

	ar := &blockchain.AmountResponse{Amount: amount}
	m, _ := ar.MarshalJSON()

	w.WriteHeader(http.StatusOK)
	w.Write(m)
}

func (s *Server) ConsensusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bc := s.GetBlockchain()
	bc.ResolveConflicts()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(utils.JsonStatus("Consensus resolved"))
}
