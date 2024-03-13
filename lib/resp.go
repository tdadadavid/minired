package lib

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
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


type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{ reader: bufio.NewReader(rd) }
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		byte, err := r.reader.ReadByte()
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


func (r *Resp) readInteger() (num int, n int, err error) {
	line, n, err := r.readLine()
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

func (r *Resp) Read() (Value, error) {
	resp_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch resp_type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Println("Unknown type: ", string(resp_type))
		return Value{}, nil
	}
}


func (r *Resp) readArray() (Value, error) {
	value := Value{}
	value.typ = "array"

	// get the len of the array by reading the next character.
	// ["*", "2", ""]
	//        ^
	// the 2 specifies the number of elements in the array request.
	arr_len, _, err := r.readInteger()
	if err != nil {
		return value, err
	}

	value.array = make([]Value, 0)

	for i := 0; i < arr_len; i++ {
		// recursion happens here [read the next stream of bytes.]
		val, err := r.Read() 
		if err != nil {
			return value, err
		}

		value.array = append(value.array, val)
	}

	return value, nil
}

func (r *Resp) readBulk() (Value, error) {
	val := Value{}
	val.typ = "bulk"

	// read next byte to know the length of the string
	// indicates the length of the string 
	size, _, err := r.readInteger()
	if err != nil {
		return val, err
	}

	bulk := make([]byte, size)
	r.reader.Read(bulk)
	val.bulk = string(bulk)

	//read trailing line [CLRF]
	r.readLine()

	return val, nil
}
