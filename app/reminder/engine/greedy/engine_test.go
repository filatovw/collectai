package greedy_test

import (
	"testing"
	"time"

	"github.com/filatovw/collectai/app/reminder/engine/greedy"
	"github.com/stretchr/testify/assert"
)

func TestEngineMalformedInput(t *testing.T) {
	check := func(input [][]string, host string) {
		engine := greedy.Engine{}
		err := engine.Init(input, host)
		assert.Error(t, err)
	}
	check([][]string{}, "localhost")
	check([][]string{[]string{"me@mail.com", "message to me", "0s-3s"}}, "")
	check([][]string{[]string{"me@mail.com", "message to me", "not an offset format"}}, "localhost")
	check([][]string{[]string{"me@mail.com", "message to me"}}, "localhost")
	check([][]string{[]string{"me@mail.com", "message to me", ""}}, "localhost")
}

func TestEngineCorrectInput(t *testing.T) {
	check := func(input [][]string, host string, expected []greedy.User) {
		engine := greedy.Engine{}
		err := engine.Init(input, host)
		assert.NoError(t, err)
		assert.Equal(t, expected, engine.Data())
		t.Logf(`%#v`, engine.Data())
	}
	check([][]string{[]string{"me@mail.com", "message to me", "0s"}}, "localhost",
		[]greedy.User{greedy.User{Offsets: []time.Duration{0}, Meta: greedy.Meta{Email: "me@mail.com", Text: "message to me"}}})

	check([][]string{[]string{"me@mail.com", "message to me", "0s-13s"}}, "localhost",
		[]greedy.User{greedy.User{Offsets: []time.Duration{0, 13000000000}, Meta: greedy.Meta{Email: "me@mail.com", Text: "message to me"}}})

	check([][]string{[]string{"me@mail.com", "message to me", "-0s-13s"}}, "localhost",
		[]greedy.User{greedy.User{Offsets: []time.Duration{0, 13000000000}, Meta: greedy.Meta{Email: "me@mail.com", Text: "message to me"}}})
}

func TestEngineProcess(t *testing.T) {

}
