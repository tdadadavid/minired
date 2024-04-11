package lib

// doc: https://redis.io/docs/latest/commands/ping/
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "string", Str: "PONG"}
	}

	// the syntax of the "PING" command with an arg is 
	// PING arguement? (without space)
	// PING hello world will respond with "hello", 
	// discarding the world, if you want the world use
	// PING "hello world"
	return Value{Typ: "bulk", Bulk: args[0].Bulk }
}

// the string is the command while the func is the handler
var CommandHandlers = map[string]func([]Value) Value{
	"ping": ping,
}
