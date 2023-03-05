package main

import (
	"fmt"
	"net/http"

	"github.com/Rha02/block-beard/src/blockchain"
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
