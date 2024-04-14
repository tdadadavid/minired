package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAOFWrite(t *testing.T){
	var path = "minired_test.aof"
	t.Cleanup(func() {
		os.Remove(path)
	})

	t.Run("It writes redis commands to the aof file", func(t *testing.T) {
		aof, _ := NewAppendOnlyFile(path)
		commands := Value{
			Array: []Value {
				{Typ: "bulk", Bulk: "set"},
				{Typ: "bulk", Bulk: "key"},
				{Typ: "bulk", Bulk: "value"},
			},
		}

		err := aof.Write(commands)
		assert.Nil(t, err)
	})
	
}