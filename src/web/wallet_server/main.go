package main

import (
	"flag"
	"log"
	"os"
)

func init() {
	log.SetPrefix("Wallet Server: ")
}

func main() {
	port := flag.Uint("port", 8080, "port to listen on")
	gateway := flag.String("gateway", "http://localhost:3000", "address of the blockchain server")
	flag.Parse()

	println(os.Getwd())

	server := NewServer(uint16(*port), *gateway)
	server.Start()
}
