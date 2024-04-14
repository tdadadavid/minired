package lib

import (
	"context"
	"fmt"
	"sync"
)

// the string is the command while the func is the handler
var CommandHandlers = map[string]func(ctx context.Context, val []Value) Value{
	"ping":    ping,
	"set":     set,
	"get":     get,
	"hset":    hset,
	"hget":    hget,
	"hgetall": hgetall,
}

type SimpleStore struct {
	mu        sync.RWMutex
	kvStore   map[string]string
	hashStore map[string]map[string]string
}

// for testing purposes.
var KvStore SimpleStore = SimpleStore{
	kvStore:   map[string]string{},
	hashStore: map[string]map[string]string{},
	mu:        sync.RWMutex{},
}

// doc: https://redis.io/docs/latest/commands/ping/
func ping(_ context.Context, args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "string", Str: "PONG"}
	}

	// the syntax of the "PING" command with an arg is
	// PING arguement? (without space)
	// PING hello world will respond with "hello",
	// discarding the world, if you want the world use
	// PING "hello world"
	return Value{Typ: "bulk", Bulk: args[0].Bulk}
}

// doc: https://redis.io/docs/latest/commands/set/
func set(_ context.Context, args []Value) Value {
	KvStore.mu.Lock()
	defer KvStore.mu.Unlock()

	if len(args) != 2 {
		return Value{Typ: "error", Str: "ERR incorrect number of arguements for the 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk
	KvStore.kvStore[key] = value

	return Value{Typ: "string", Str: "OK"}
}

func get(_ context.Context, args []Value) Value {
	KvStore.mu.RLock()
	defer KvStore.mu.RUnlock()

	if len(args) != 1 {
		return Value{Typ: "error", Str: "ERR incorrect number of arguements for the 'get' command"}
	}

	key := args[0].Bulk

	value, ok := KvStore.kvStore[key]
	if !ok {
		return Value{Typ: "string", Str: "nil"}
	}

	return Value{Typ: "string", Str: value}
}

// doc: https://redis.io/docs/latest/commands/hset/
func hset(_ context.Context, args []Value) Value {
	KvStore.mu.Lock()
	defer KvStore.mu.Unlock()

	if len(args) < 3 {
		return Value{Typ: "error", Str: "ERR incorrect number of arguements for the 'hset' command"}
	}

	hashKey := args[0].Bulk

	store := make(map[string]string)
	values := args[1:]
	values_len := len(values)

	// it is possible client sends key without value.
	if values_len%2 != 0 {
		return Value{Typ: "error", Str: "ERR incorrect number of arguements for the 'hset' command"}
	}

	for i := 0; i < values_len; i++ {
		store[values[i].Bulk] = values[i+1].Bulk
		i++
	}

	KvStore.hashStore[hashKey] = store

	return Value{Typ: "string", Str: fmt.Sprint(len(store))}
}

// doc: https://redis.io/docs/latest/commands/hget/
func hget(_ context.Context, args []Value) Value {
	KvStore.mu.Lock()
	defer KvStore.mu.Unlock()

	if len(args) < 2 {
		return Value{Typ: "error", Str: "ERR incorrect number of arguements for the 'hget' command"}
	}

	hashKey := args[0].Bulk
	subKey := args[1].Bulk

	result, ok := KvStore.hashStore[hashKey][subKey]
	if !ok {
		return Value{Typ: "string", Str: "nil"}
	}

	return Value{Typ: "string", Str: result}
}

// doc: https://redis.io/docs/latest/commands/hgetall/
func hgetall(_ context.Context, args []Value) Value {
	KvStore.mu.RLock()
	defer KvStore.mu.RUnlock()

	if len(args) < 1 {
		return Value{Typ: "error", Str: "ERR incorrect number of arguements for the 'hgetall' command"}
	}

	hashKey := args[0].Bulk

	values, ok := KvStore.hashStore[hashKey]
	if !ok {
		return Value{Typ: "string", Str: "nil"}
	}

	results := []Value{}
	for k, v := range values {
		results = append(results, Value{Typ: "string", Str: k})
		results = append(results, Value{Typ: "string", Str: v})
	}

	return Value{Typ: "array", Array: results}
}
