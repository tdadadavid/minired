package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingCommand(t *testing.T) {
	tests := []struct {
		desc    string
		args    string
		expects string
	}{
		{desc: "", expects: "PONG"},
		{desc: "It returns the arguement passed to it when called", args: "hello world", expects: "hello world"},
	}

	t.Run("It returns PONG when no arguement is passed", func(t *testing.T) {

		value := &Value{}
		result := ping(value.Array)

		assert.Equal(t, tests[0].expects, result.Str)

	})

	t.Run("It returns the arguement passed to it when called", func(t *testing.T) {
		value := &Value{
			Array: []Value{
				{Typ: "string", Bulk: tests[1].args},
			},
		}
		result := ping(value.Array)

		assert.Equal(t, "bulk", result.Typ)
		assert.Equal(t, tests[1].expects, result.Bulk)

	})
}
