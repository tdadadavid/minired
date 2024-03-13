package lib

import (
	"testing"
)


func TestRESPParser(t *testing.T) {
	// good scenerio
	t.Run("It parses bulk string of any length", func(t *testing.T) {
		// request := "$6\r\nsixtyo\r\n"
		// result := ParseResp(request)

		// assert.Contains(t, result, "sixtyo")
	})

}

// input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
// 	resp := NewRESP(strings.NewReader(input))

// 	val, err := resp.Read()
// 	if err != nil {
// 		fmt.Println("Error: ", err.Error())
// 		os.Exit(1)
// 	}

// 	fmt.Println("Value: ", val)