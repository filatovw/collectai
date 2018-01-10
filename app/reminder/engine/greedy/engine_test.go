package greedy_test

import (
	"testing"

	"github.com/filatovw/collectai/app/reminder/engine/greedy"
	"github.com/stretchr/testify/assert"
)

func TestEngineEmpty(t *testing.T) {
	check := func(input [][]string, host string) {
		engine := greedy.Engine{}
		err := engine.Init(input, host)
		assert.Error(t, err)
	}
	check([][]string{}, "localhost")
	check([][]string{[]string{"aaa"}}, "")
	check([][]string{[]string{"me@mail.com", "message to me", "not an offset format"}}, "localhost")
}
