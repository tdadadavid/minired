package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Server listening on port 6379")


	//create a tcp listener on port 6379
	listener, err := net.Listen("tcp", "6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	//use listener to accepts any requests from client.
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	// close the connection once done with the requests
	// [remember to study tcp,http protocols]
	defer conn.Close() 

	for {
		// create a buffer in memory
		buff := make([]byte, 1024)

		// read but ignore the content from the client
		_, err := conn.Read(buff)
		if err != nil { 
			if err == io.EOF { //if there is end-of-file error break from the loop
				break
			}
			fmt.Println("Error reading from client", err.Error())
			os.Exit(1) // stop the server.
		}

		// for any request the user makes return a pong.
		conn.Write([]byte("+Ok\r\n"))
	}

	fmt.Println("Hello world.")
}