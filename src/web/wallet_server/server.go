package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/Rha02/block-beard/src/blockchain"
	"github.com/Rha02/block-beard/src/utils"
	"github.com/Rha02/block-beard/src/wallet"
)

const tempDir = "./wallet_server/templates"

type Server struct {
	port    uint16
	gateway string
}

func NewServer(port uint16, gateway string) *Server {
	return &Server{port, gateway}
}

func (s *Server) Port() uint16 {
	return s.port
}

func (s *Server) Gateway() string {
	return s.gateway
}

func (s *Server) Index(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	t := template.Must(template.ParseFiles(path.Join(tempDir, "index.html")))
	t.Execute(rw, nil)
}

func (s *Server) Wallet(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	myWallet := wallet.NewWallet()
	walletJSON, _ := myWallet.MarshalJSON()
	rw.Write(walletJSON)
}

func (s *Server) PostTransactionHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var t wallet.TransactionRequest
	err := decoder.Decode(&t)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(utils.JsonStatus("Invalid transaction"))
		println("Error decoding")
		return
	}

	if !t.IsValid() {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(utils.JsonStatus("Invalid transaction: missing fields"))
		println("Error: missing fields")
		return
	}

	publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
	privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
	amount, err := strconv.ParseFloat(*t.Amount, 32)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(utils.JsonStatus("Invalid transaction: invalid amount"))
		println("Error: invalid amount")
		return
	}
	amount32 := float32(amount)

	rw.Header().Add("Content-Type", "application/json")

	transaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderAddress, *t.RecipientAddress, amount32)
	signature := transaction.GenerateSignature()
	signatureStr := signature.ToString()

	tr := blockchain.TransactionRequest{
		SenderPublicKey:  t.SenderPublicKey,
		SenderAddress:    t.SenderAddress,
		RecipientAddress: t.RecipientAddress,
		Amount:           &amount32,
		Signature:        &signatureStr,
	}

	m, _ := json.Marshal(tr)
	buf := bytes.NewBuffer(m)

	res, err := http.Post(s.Gateway()+"/transactions", "application/json", buf)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(utils.JsonStatus("Error posting transaction to blockchain"))
		println("Error posting transaction to blockchain")
		return
	}
	if res.StatusCode == http.StatusCreated {
		rw.WriteHeader(http.StatusCreated)
		rw.Write(utils.JsonStatus("Transaction posted to blockchain"))
		println("Transaction posted to blockchain")
		return
	}
	println("Error posting transaction to blockchain: " + res.Status)
	rw.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) WalletAmountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	bcAddress := r.URL.Query().Get("blockchain_address")
	endpoint := fmt.Sprintf("%s/amount?blockchain_address=%s", s.Gateway(), bcAddress)

	res, err := http.Get(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JsonStatus("Error getting amount from blockchain"))
		return
	}

	if res.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(res.Body)
		var amount blockchain.AmountResponse
		err := decoder.Decode(&amount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.JsonStatus("Error decoding amount from blockchain"))
			return
		}

		w.WriteHeader(http.StatusOK)
		m, _ := json.Marshal(struct {
			Message string  `json:"message"`
			Amount  float32 `json:"amount"`
		}{
			Message: "Amount retrieved from blockchain",
			Amount:  amount.Amount,
		})
		w.Write(m)

		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(utils.JsonStatus("Error getting amount from blockchain"))
}

func (s *Server) Start() {
	http.HandleFunc("/", s.Index)
	http.HandleFunc("/wallet", s.Wallet)
	http.HandleFunc("/wallet/amount", s.WalletAmountHandler)
	http.HandleFunc("/transaction", s.PostTransactionHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
