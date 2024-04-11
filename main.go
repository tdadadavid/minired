package main

import (
	"context"
	"fmt"
	"go_redis/lib"
	"net"
	"strings"
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

		// if the request sent is not an array type ignore it
		if value.Typ != "array" {
			continue
		}

		// ignore empty request
		if len(value.Array) == 0 {
			continue
		}

		// extract the command & arguements from the request
		// redis commands are case-insensitive [https://redis.io/docs/latest/commands/command/]
		command := strings.ToLower(value.Array[0].Bulk)
		args := value.Array[1:]

		// get the handler for the command .
		handler := lib.CommandHandlers[command]

		// and feed it the arguements
		result := handler(context.Background(), args)

		writer := lib.NewWriter(conn)
		writer.Write(result)
	}

}
