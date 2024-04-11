package lib

import (
	"context"
	"sync"
)

type KeyValue struct {
	mu    sync.RWMutex
	store map[string]string
}

var kvStore KeyValue = KeyValue{
	store: map[string]string{},
	mu:    sync.RWMutex{},
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
	kvStore.mu.Lock()
	defer kvStore.mu.Unlock()

	key := args[0].Bulk
	value := args[1].Bulk

	kvStore.store[key] = value
	return Value{Typ: "string", Str: "OK"}
}

func get(_ context.Context, args []Value) Value {
	kvStore.mu.RLock()
	defer kvStore.mu.RUnlock()

	key := args[0].Bulk

	value := kvStore.store[key]

	return Value{Typ: "string", Str: value}
}

// the string is the command while the func is the handler
var CommandHandlers = map[string]func(ctx context.Context, val []Value) Value{
	"ping": ping,
	"set":  set,
	"get":  get,
}
