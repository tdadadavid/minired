package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Reference study: https://redis.io/docs/reference/protocol-spec/

// datatypes supported in the RESP
const (
	STRING = '+'
	BULK = '$'
	ERROR = '-'
	INTEGER = ':'
	ARRAY = '*'
)

// will hold the request arguements and command
// it will be used in the serialization/desrialization of reqeust
type Value struct {
	typ string // holds the datatype of the value from the request.
	num int32 // holds all integer request.
	str string // holds all string requests.
	bulk string // holds all bulk string requests.
	array []Value // holds all array requests
}


type RESP struct {
	reader *bufio.Reader
}

func NewRESP(reader io.Reader) *RESP {
	return &RESP{ reader: bufio.NewReader(reader) }
}

func (this *RESP) readLine() (line []byte, n int, err error) {
	for {
		byte, err := this.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		n += 1
		line = append(line, byte)

		// if the line is greater than 2 and the second-to-the-last
		// character is the 'Carriage Return' break the loop.
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}

	// return from the beginning of the string to the second-to-the-last character
	// also return the number of characters on the line just read.
	return line[:len(line)-2], n , nil
}


func (this *RESP) readInteger() (num int, n int, err error) {
	line, n, err := this.readLine()
	if err != nil {
		return 0, 0, nil
	}

	// convert the integer to a 64-bit integer in base 10.
	_64bitInteger, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(_64bitInteger), n, nil
}

func (resp *RESP) Read() (Value, error) {
	resp_type, err := resp.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch resp_type {
	case ARRAY:
		return resp.readArray()
	case BULK:
		return resp.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(resp_type))
		return Value{}, nil
	}
}

func (resp *RESP) readArray() (Value, error) {
	_, err := resp.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	value := Value{ typ: "array" }

	for {
		// read the /r/n bytes 
		resp.reader.ReadByte()
		resp.reader.ReadByte()

		_type, err := resp.reader.ReadByte()
		if err != nil {
			break
		}

		switch _type {
		case BULK:
			val, err := resp.readBulk()
			if err != nil {
				break
			}
			value.array = append(value.array, val)
		default:
			fmt.Println("Unknown type: ", string(_type))
			return Value{}, nil
		}
	}
	return value, nil
}

func (resp *RESP) readBulk() (Value, error) {
	val := Value{ typ: "bulk" }

	// read next byte to know the length of the string
	// indicates the length of the string 
	size, _ := resp.reader.ReadByte()

	inputSize, _ := strconv.ParseInt(string(size), 10, 64)

	// current position in parsing: ['$', '5', '/r', '/n', 'ahmed', '/r', '/n']
	//																			    ^
	resp.reader.ReadByte()

	// current position in parsing: ['$', '5', '/r', '/n', 'ahmed', '/r', '/n']
	//																						     ^
	resp.reader.ReadByte()


	name := make([]byte, inputSize)
	resp.reader.Read(name)

	val.str = string(name)
	return val, nil
}

func main() {
	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	resp := NewRESP(strings.NewReader(input))

	val, err := resp.Read()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Value: ", val)

}


// // this is a simple RESP parser.
// func main() {
// 	input := "$5\r\nahmed\r\n"
// 	ParseResp(input)
// }

// func ParseResp(input string) string {
// 	// given a string like this '$5/r/n/ahmed/r/n'
// 	// we need to parse this, this string according to RESP
// 	// the example given is a bulk string starts
// 	// with $ and the length of the value (5))

// 	reader := bufio.NewReader(strings.NewReader(input))

// 	// read first byte from the buffer
// 	firstCharBuf, _ := reader.ReadByte()

// 	// now if the first byte read is not '$' then we reject it
// 	// why? becuase we are first considering a bulk string.
// 	if firstCharBuf != '$' {
// 		fmt.Println("Error: Invalid type provided. Only bulk strings are supported")
// 		os.Exit(1)
// 	}

// 	// indicates the length of the string 
// 	size, _ := reader.ReadByte()

// 	inputSize, _ := strconv.ParseInt(string(size), 10, 64)

// 	// current position in parsing: ['$', '5', '/r', '/n', 'ahmed', '/r', '/n']
// 	//																			    ^
// 	reader.ReadByte()

// 	// current position in parsing: ['$', '5', '/r', '/n', 'ahmed', '/r', '/n']
// 	//																						     ^
// 	reader.ReadByte()


// 	name := make([]byte, inputSize)
// 	reader.Read(name)

// 	fmt.Println(string(name))
// 	return string(name)
// }