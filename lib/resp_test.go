package lib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadArray_WhenCorrectSyntazIsNotPassed(t *testing.T) {
	t.Run("It returns error when wrong RESP array syntax is sent", func(t *testing.T) {
		wrong_arr_string := "*wrong\r\n$5\r\nbayom"
		reader := strings.NewReader(wrong_arr_string)
		resp := NewResp(reader)
		_, err := resp.readArray()

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid syntax")
	})
}

func TestReadArray_WhenCorrectSyntazIsPassed(t *testing.T) {
	t.Run("The contents of the array is the  same as the RESP from the client", func(t *testing.T) {
		// ignore the RESP array identicator [*] because READ already handles that.
		resp_arr := "2\r\n$1\r\nA\r\n$1\r\nB\r\n"
		reader := strings.NewReader(resp_arr)
		resp := NewResp(reader)
		result, err := resp.readArray()

		assert.Nil(t, err)

		assert.Len(t, result.Array, 2)
		assert.Equal(t, result.Array[0].Bulk, "A")
		assert.Equal(t, result.Array[1].Bulk, "B")
	})
}

func TestReadBulk_WhenCorrectSyntacIsNotPassed(t *testing.T) {
	t.Run("It fails when a wrong bulk string syntax is sent from client.", func(t *testing.T) {
		wrong_bulk_string := "$n8\r\nsixtyo\r\n"
		reader := strings.NewReader(wrong_bulk_string)
		resp := NewResp(reader)
		_, err := resp.readBulk()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})
}

func TestReadBulk_WhenCorrectSyntaxIsPassed(t *testing.T) {
	t.Run("It returns the actual string when the correct bulk string is sent from client", func(t *testing.T) {
		// we need to remove the $ symbol because the "Read" function already parses it
		right_bulk_string := "6\r\nsixtyo\r\n"
		reader := strings.NewReader(right_bulk_string)
		resp := NewResp(reader)
		result, err := resp.readBulk()

		assert.Nil(t, err)
		assert.NotNil(t, result.Bulk)
		assert.Equal(t, result.Bulk, "sixtyo")
	})
}

func TestBulkString_ExpectValueTypeToBeBulk(t *testing.T) {
	t.Run("The type of Value when after parsing is 'bulk'", func(t *testing.T) {
		right_bulk_string := "6\r\nsixtyo\r\n"
		reader := strings.NewReader(right_bulk_string)
		resp := NewResp(reader)
		result, err := resp.readBulk()

		assert.Nil(t, err)
		assert.NotNil(t, result.Bulk)
		assert.Equal(t, result.Typ, "bulk")
	})
}
