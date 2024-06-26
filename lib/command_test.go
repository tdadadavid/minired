package lib

import (
	"context"
	"math/rand"
	"strconv"
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

		if KvStore.kvStore[key] != val {
			t.Fatalf(key, "was not assigned value ", val)
		}

		assert.Equal(t, result.Str, "OK")
	})

	t.Run("It returns an error value when incorrect number of args is sents", func(t *testing.T) {
		key := "admin"

		args := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: key},
			},
		}

		result := set(context.Background(), args.Array)

		assert.Contains(t, result.Str, "ERR")
		assert.Equal(t, result.Typ, "error")
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

		assert.NotEqual(t, KvStore.kvStore[key], val1)
		assert.Equal(t, KvStore.kvStore[key], val2)
	})
}

func TestGetCommand(t *testing.T) {
	t.Run("It returns an error when incorrect number of args is sent", func(t *testing.T) {

		args := &Value{
			Array: []Value{},
		}

		result := get(context.Background(), args.Array)

		assert.Contains(t, result.Str, "ERR")
		assert.Equal(t, result.Typ, "error")
	})

	t.Run("It returns nil if a value has not been set for the provided key", func(t *testing.T) {
		key := strconv.Itoa(rand.Intn(50)) //random key

		args := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: key},
			},
		}

		result := get(context.Background(), args.Array)

		assert.Contains(t, result.Str, "nil")
		assert.Equal(t, result.Typ, "string")
	})

	t.Run("It retrieves the value of the key set", func(t *testing.T) {
		key := "admin"
		val := "king"

		args := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: key},
				{Typ: "bulk", Bulk: val},
			},
		}
		_ = set(context.Background(), args.Array)

		args = &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: key},
			},
		}
		result := get(context.Background(), args.Array)

		assert.Equal(t, result.Str, val)
		assert.Equal(t, result.Typ, "string")
	})
}

func TestHSetCommand(t *testing.T) {
	t.Run("it returns an error when the wrong number of key value pairs for the hash key is provided by the client", func(t *testing.T) {
		hashKey := "application"
		role := "admin"
		person := "king"
		role2 := "super-admin"

		args := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: hashKey},
				{Typ: "bulk", Bulk: role},
				{Typ: "bulk", Bulk: person},
				{Typ: "bulk", Bulk: role2}, // there is no value to the key entry
			},
		}

		result := hset(context.Background(), args.Array)

		assert.Contains(t, result.Str, "ERR")
		assert.Equal(t, result.Typ, "error")
	})

	t.Run("it returns an error when the wrong number of arguement is sent from the client", func(t *testing.T) {
		args := &Value{
			Array: []Value{},
		}

		result := hset(context.Background(), args.Array)

		assert.Contains(t, result.Str, "ERR")
		assert.Equal(t, result.Typ, "error")
	})

	t.Run("It set the given hash key to the provided 'key-value' value", func(t *testing.T) {

		hashKey := "application"
		role := "admin"
		person := "king"

		args := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: hashKey},
				{Typ: "bulk", Bulk: role},
				{Typ: "bulk", Bulk: person},
			},
		}

		result := hset(context.Background(), args.Array)

		if KvStore.hashStore[hashKey] == nil {
			t.Fatal("values was not stored in the hash store", KvStore.hashStore)
		}

		if KvStore.hashStore[hashKey][role] != person {
			t.Fatalf(role, "was not assigned value ", person)
		}

		assert.Len(t, KvStore.hashStore[hashKey], 1)
		assert.Equal(t, result.Str, "1")
	})

	t.Run("It updates the key's value on every call", func(t *testing.T) {
		hashKey := "admin"
		field := "status"
		value := "monarch"
		final_value := "king"

		args := []Value{
			{
				Array: []Value{
					{Typ: "bulk", Bulk: hashKey},
					{Typ: "bulk", Bulk: field},
					{Typ: "bulk", Bulk: value},
				},
			},
			{
				Array: []Value{
					{Typ: "bulk", Bulk: hashKey},
					{Typ: "bulk", Bulk: field},
					{Typ: "bulk", Bulk: final_value},
				},
			},
		}

		for _, v := range args {
			hset(context.Background(), v.Array)
		}

		assert.NotEqual(t, KvStore.hashStore[hashKey][field], value)
		assert.Equal(t, KvStore.hashStore[hashKey][field], final_value)
	})
}

func TestHGetCommand(t *testing.T) {
	t.Run("It returns an error when incorrect number of args is sent from the client", func(t *testing.T) {
		args := &Value{
			Array: []Value{},
		}

		result := hget(context.Background(), args.Array)

		assert.Contains(t, result.Str, "ERR")
		assert.Equal(t, result.Typ, "error")
	})

	t.Run("It returns nil if a value has not been set for the provided key", func(t *testing.T) {
		hashKey := strconv.Itoa(rand.Intn(50)) //random key
		key := strconv.Itoa(rand.Intn(50))     //random key

		args := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: hashKey},
				{Typ: "bulk", Bulk: key},
			},
		}

		result := hget(context.Background(), args.Array)

		assert.Contains(t, result.Str, "nil")
		assert.Equal(t, result.Typ, "string")
	})

	t.Run("It retrieves the value of the key set", func(t *testing.T) {
		hashKey := strconv.Itoa(rand.Intn(50)) //random key
		key := "admin"
		val := "king"

		args := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: hashKey},
				{Typ: "bulk", Bulk: key},
				{Typ: "bulk", Bulk: val},
			},
		}
		_ = hset(context.Background(), args.Array)

		args = &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: hashKey},
				{Typ: "bulk", Bulk: key},
			},
		}
		result := hget(context.Background(), args.Array)

		assert.Equal(t, result.Str, val)
		assert.Equal(t, result.Typ, "string")
	})
}

func TestHGetAllCommand(t *testing.T) {
	t.Run("It returns an error value when incorrect number of args is sents", func(t *testing.T) {
		args := &Value{
			Array: []Value{},
		}

		result := hgetall(context.Background(), args.Array)

		assert.Contains(t, result.Str, "ERR")
		assert.Equal(t, result.Typ, "error")
	})

	t.Run("It returns nil if a value has not been set for the provided key", func(t *testing.T) {
		hashKey := strconv.Itoa(rand.Intn(50)) //random key

		args := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: hashKey},
			},
		}

		result := hgetall(context.Background(), args.Array)

		assert.Contains(t, result.Str, "nil")
		assert.Equal(t, result.Typ, "string")
	})

	t.Run("It returns all the key value pairs when the hash key exists", func(t *testing.T) {
		hashKey := strconv.Itoa(rand.Intn(50)) //random key
		key := "admin"
		val := "king"

		args := &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: hashKey},
				{Typ: "bulk", Bulk: key},
				{Typ: "bulk", Bulk: val},
			},
		}
		_ = hset(context.Background(), args.Array)

		args = &Value{
			Array: []Value{
				{Typ: "bulk", Bulk: hashKey},
			},
		}

		result := hgetall(context.Background(), args.Array)

		assert.Len(t, result.Array, 2)
		assert.Equal(t, result.Typ, "array")
	})
}
