package main

import (
	"log"
	"minired/server"
)

func main() {
	server := server.NewServer(":6379")
	log.Fatal(server.Start())
}
