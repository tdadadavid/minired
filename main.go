package main

import (
	"log"
	server "minired/lib"
)

func main() {
	server := server.NewServer(":6379")
	log.Fatal(server.Start())
}
