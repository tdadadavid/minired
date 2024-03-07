package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Reference study: https://redis.io/docs/reference/protocol-spec/

// this is a simple RESP parser.
func main() {
	input := "$5\r\nahmed\r\n"
	ParseResp(input)
}

func ParseResp(input string) {
	// given a string like this '$5/r/n/ahmed/r/n'
	// we need to parse this, this string according to RESP
	// is a bulk string, (a bulk string starts with $ and the length of the value (5))

	reader := bufio.NewReader(strings.NewReader(input))

	// read first byte from the buffer
	firstCharBuf, _ := reader.ReadByte()

	// now if the first byte read is not '$' then we reject it
	if firstCharBuf != '$' {
		fmt.Println("Error: Invalid type provided. Only bulk strings are supported")
		os.Exit(1)
	}

	size, _ := reader.ReadByte()

	inputSize, _ := strconv.ParseInt(string(size), 10, 64)

	// current position in parsing: ['$', '5', '/r', '/n', 'ahmed', '/r', '/n']
	//																			    ^
	reader.ReadByte()

	// current position in parsing: ['$', '5', '/r', '/n', 'ahmed', '/r', '/n']
	//																						     ^
	reader.ReadByte()


	name := make([]byte, inputSize)
	reader.Read(name)

	fmt.Println(string(name))

}