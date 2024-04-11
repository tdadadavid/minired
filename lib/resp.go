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
	NULL = '_'
)

// will hold the request arguements and command
// it will be used in the serialization/desrialization of reqeust
type Value struct {
	typ string // holds the datatype of the value from the requests.
	num int32 // holds all integer requests.
	str byte // holds all string requests.
	bulk string // holds all bulk string requests.
	array []Value // holds all array requests
	error string
}


type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{ reader: bufio.NewReader(rd) }
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

// Convert respsonse into RESP type.
func (v Value) Marshal() ([]byte)  {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshallError()
	default:
		return []byte{}
	}
}

// Structure of RESP "error":
// _[Carriage Return Line Feed]
// doc: https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-errors
func (v Value) marshallNull() []byte {
	var bytes []byte
	bytes = append(bytes, NULL)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// Structure of RESP "nll":
// -[Error Message][Carriage Return Line Feed]
// doc: https://redis.io/docs/latest/develop/reference/protocol-spec/#nulls
func (v Value) marshallError() []byte {
	var bytes []byte

	bytes = append(bytes, ERROR)
	bytes = append(bytes, []byte(v.error)...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}


// Structure of RESP "array":
// *[len-of-array][Carriage Return Line Feed][firstElement]...[elementN][Carriage Return Line Feed]
// doc: https://redis.io/docs/latest/develop/reference/protocol-spec/#arrays
func (v Value) marshalArray() []byte {
	var bytes []byte

	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(v.array))...)
	bytes = append(bytes, '\r', '\n')
	
	for i := 0; i < len(v.array); i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes;
}

// Structure of RESP "string":
// +[string][Carriage Return Line Feed]
// doc: https://redis.io/docs/latest/develop/reference/protocol-spec/#simple-strings
func (v Value) marshalString() []byte {
	var result []byte

	result = append(result, STRING)
	result = append(result, v.str)
	result = append(result, '\r' , '\n')

	return result
}

// The structure of RESP "bulk":
// $[len-of-the-bulk-sting][Carriage Return Line Feed][bulk-string][Carriage-Return-Line-Feed]
// doc: https://redis.io/docs/latest/develop/reference/protocol-spec/#bulk-strings
func (v Value) marshalBulk() []byte {
	var bytes []byte

	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
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
