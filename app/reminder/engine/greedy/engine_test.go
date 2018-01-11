package greedy_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
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
	}
	check([][]string{[]string{"me@mail.com", "message to me", "0s"}}, "localhost",
		[]greedy.User{greedy.User{Offsets: []time.Duration{0}, Meta: greedy.Meta{Email: "me@mail.com", Text: "message to me"}}})

	check([][]string{[]string{"me@mail.com", "message to me", "0s-13s"}}, "localhost",
		[]greedy.User{greedy.User{Offsets: []time.Duration{0, 13000000000}, Meta: greedy.Meta{Email: "me@mail.com", Text: "message to me"}}})

	check([][]string{[]string{"me@mail.com", "message to me", "-0s-13s"}}, "localhost",
		[]greedy.User{greedy.User{Offsets: []time.Duration{0, 13000000000}, Meta: greedy.Meta{Email: "me@mail.com", Text: "message to me"}}})
}

func TestEngineProcess(t *testing.T) {
	check := func(input [][]string, expected int32) {
		var (
			counter int32
			step    = int32(1)
		)

		ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&counter, step)
			assert.Equal(t, "/messages", r.URL.EscapedPath())
			assert.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(http.StatusCreated)
			body := []byte("{}")
			w.Write(body)
		}))
		defer ts.Close()
		ts.Start()

		engine := greedy.Engine{}
		engine.Init(input, ts.URL)
		err := engine.Process()
		assert.NoError(t, err)
		assert.Equal(t, expected, atomic.LoadInt32(&counter))
	}

	check([][]string{
		[]string{"me@mail.com", "message to me", "0s"}}, int32(1))

	check([][]string{
		[]string{"me@mail.com", "message to me", "0s-10ms"},
		[]string{"another@mail.com", "another message", "1ms-11ms-12ms"}}, int32(5))

	check([][]string{
		[]string{"me@mail.com", "message to me", "0s"},
		[]string{"another@mail.com", "another message", "1ms"}}, int32(2))
}

func TestEngineProcessWithPaid(t *testing.T) {
	check := func(input [][]string, expected int32) {
		var (
			counter int32
			step    = int32(1)
		)

		ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&counter, step)
			assert.Equal(t, "/messages", r.URL.EscapedPath())
			assert.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(http.StatusCreated)
			body := []byte(`{"paid": true}`)
			w.Write(body)
		}))
		defer ts.Close()
		ts.Start()

		engine := greedy.Engine{}
		engine.Init(input, ts.URL)
		err := engine.Process()
		assert.NoError(t, err)
		assert.Equal(t, expected, atomic.LoadInt32(&counter))
	}

	check([][]string{
		[]string{"me@mail.com", "message to me", "0s"}}, int32(1))

	check([][]string{
		[]string{"me@mail.com", "message to me", "0s-10ms"},
		[]string{"another@mail.com", "another message", "1ms-11ms-12ms"}}, int32(2))

	check([][]string{
		[]string{"me@mail.com", "message to me", "0s"},
		[]string{"another@mail.com", "another message", "1ms"}}, int32(2))
}
