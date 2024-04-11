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

func TestSetCommand(t *testing.T) {
	t.Run("It set the given key to the provided value", func(t *testing.T) {
		key := "admin"
		val := "king"

		args := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: key},
				{Typ: "bulk", Bulk: val},
			},
		}

		result := set(context.Background(), args.Array)

		if KvStore.store[key] != val {
			t.Fatalf(key, "was not assigned value ", val)
		}

		assert.Equal(t, result.Str, "OK")
	})

	t.Run("It updates the key's value on every call", func(t *testing.T) {
		key := "admin"
		val1 := "king"
		val2 := "monarch"

		args := []Value{
			{
				Array: []Value{
					{Typ: "bulk", Bulk: key},
					{Typ: "bulk", Bulk: val1},
				},
			},
			{
				Array: []Value{
					{Typ: "bulk", Bulk: key},
					{Typ: "bulk", Bulk: val2},
				},
			},
		}

		for _, v := range args {
			set(context.Background(), v.Array)
		}

		assert.NotEqual(t, KvStore.store[key], val1)
		assert.Equal(t, KvStore.store[key], val2)
	})
}
