package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestRESPParser(t *testing.T) {
	// good scenerio
	t.Run("It parses bulk string of any length", func(t *testing.T) {
		request := "$6\r\nsixtyo\r\n"
		result := ParseResp(request)

		assert.Contains(t, result, "sixtyo")
	})

}

func TestRESPParserFail(t *testing.T) {
	// bad scenerio
	t.Run("It stops the program when a non-bulk string command is sent", func(t *testing.T) {
		request := "#6\r\nsixtyo\r\n"
		result := ParseResp(request)

		assert.Contains(t, result, "sixtyo")
	})
	//TODO: check if its is possible how to handle errors not returned.
}