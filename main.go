package main

import (
	"fmt"
	"go_redis/lib"
	"net"
)

func main() {
	fmt.Println("Server listening on port 6379")


	//create a tcp listener on port 6379
	listener, err := net.Listen("tcp", ":6379")
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
		resp := lib.NewResp(conn)

		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(value)
		// for any request the user makes return a pong.
		conn.Write([]byte("+PONG\r\n"))
	}
		
}