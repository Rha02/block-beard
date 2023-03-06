package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"text/template"

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

}

func (s *Server) Start() {
	http.HandleFunc("/", s.Index)
	http.HandleFunc("/wallet", s.Wallet)
	http.HandleFunc("/transaction", s.PostTransactionHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
