package lib

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingCommand(t *testing.T) {
	tests := []struct {
		args    string
		expects string
	}{
		{expects: "PONG"},
		{args: "hello world", expects: "hello world"},
	}

	t.Run("It returns PONG when no arguement is passed", func(t *testing.T) {

		value := &Value{}
		result := ping(context.Background(), value.Array)

		assert.Equal(t, "string", result.Typ)
		assert.Equal(t, tests[0].expects, result.Str)

	})

	t.Run("It returns the arguement passed to it when called", func(t *testing.T) {
		value := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: tests[1].args},
			},
		}
		result := ping(context.Background(), value.Array)

		assert.Equal(t, "bulk", result.Typ)
		assert.Equal(t, tests[1].expects, result.Bulk)

	})
}
