package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 3000, "port to listen on")
	flag.Parse()

	log.Printf("Starting server on port %d", *port)

	server := NewServer(uint16(*port))
	server.Start()
}
