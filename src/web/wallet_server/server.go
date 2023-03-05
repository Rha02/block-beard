package main

import (
	"fmt"
	"net/http"
	"path"
	"text/template"
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

func (s *Server) Start() {
	http.HandleFunc("/", s.Index)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
