package main

import (
	"context"
	"fmt"
	"minired/lib"
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

	path := "minired.aof"
	aof, err := lib.NewAppendOnlyFile(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	aof.Read(func(value lib.Value) {
		command := strings.ToLower(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := lib.CommandHandlers[command]
		if !ok {
			fmt.Println("Command not supported: [", command, "]")
		}
		handler(context.Background(), args)
	})

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

		writer := lib.NewWriter(conn)

		// extract the command & arguements from the request
		// redis commands are case-insensitive [https://redis.io/docs/latest/commands/command/]
		command := strings.ToLower(value.Array[0].Bulk)
		args := value.Array[1:]

		if command == "set" || command == "hset" {
			aof.Write(value)
		}

		// get the handler for the command .
		handler, ok := lib.CommandHandlers[command]
		if !ok {
			fmt.Println("Command not supported: [", command, "]")
			writer.Write(lib.Value{Typ: "string", Str: ""})
			continue
		}
		// and feed it the arguements
		result := handler(context.Background(), args)

		writer.Write(result)
	}
}
