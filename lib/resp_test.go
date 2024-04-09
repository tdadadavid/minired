package lib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestReadBulk_WhenCorrectSyntacIsNotPassed(t *testing.T) {
	// good scenerio
	t.Run("It fails when a wrong bulk string syntax is sent from client.", func(t *testing.T) {
		wrong_bulk_string := "$n8\r\nsixtyo\r\n"
		reader := strings.NewReader(wrong_bulk_string)
		resp := NewResp(reader)
		_, err  := resp.readBulk()
        
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})


}

func TestReadBulk_WhenCorrectSyntaxIsPassed(t *testing.T) {
	t.Run("It returns the actual string when the correct bulk string is sent from client", func(t *testing.T) {
		// we need to remove the number because the "Read" function already parses it
		right_bulk_string := "6\r\nsixtyo\r\n" 
		reader := strings.NewReader(right_bulk_string)
		resp := NewResp(reader)
		result, err  := resp.readBulk()
        
		assert.Nil(t, err)
		assert.NotNil(t, result.bulk)
		assert.Equal(t, result.bulk, "sixtyo")
	})
}

func TestBulkString_ExpectValueTypeToBeBulk(t *testing.T) {
	t.Run("The type of Value when after parsing is 'bulk'", func(t *testing.T) {
		right_bulk_string := "6\r\nsixtyo\r\n" 
		reader := strings.NewReader(right_bulk_string)
		resp := NewResp(reader)
		result, err  := resp.readBulk()

		assert.Nil(t, err)
		assert.NotNil(t, result.bulk)
		assert.Equal(t, result.typ, "bulk")
	})
}