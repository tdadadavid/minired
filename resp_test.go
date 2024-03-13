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