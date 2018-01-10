package greedy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	CTYPE_JSON = "application/json"
)

// Engine implements Engine interface
type Engine struct {
	data []User // user data with offsets
	host string // connservice host
}

// Meta represents all we know about the notified person
type Meta struct {
	Email string `json:"email"`
	Text  string `josn:"text"`
}

type messagesResponse struct {
	Paid bool `json:"paid"`
}

type User struct {
	offsets []time.Duration
	meta    *Meta
}

// Init reads initial CSV data and parse offsets
func (e *Engine) Init(input [][]string, host string) error {
	if len(input) == 0 {
		return errors.New("empty input")
	}
	if strings.TrimSpace(host) == "" {
		return errors.New("empty host is not allowed")
	}
	e.host = host
	data := []User{}
	for i, item := range input {
		timestamps := strings.Split(item[2], "-")
		offsets := []time.Duration{}
		if len(timestamps) > 0 {
			for _, ts := range timestamps {
				offset, err := time.ParseDuration(ts)
				if err != nil {
					return fmt.Errorf(`failed to parse time offset: %s`, err)
				}
				offsets = append(offsets, offset)
			}
		} else {
			return fmt.Errorf(`failed to get time offsets; row: %d, data: %v`, i, item[2])
		}
		user := User{
			meta: &Meta{
				Email: item[0],
				Text:  item[1],
			},
			offsets: offsets,
		}
		data = append(data, user)
	}
	e.data = data
	return nil
}

// Process starts goroutines. One for each time offset.
func (e Engine) Process() error {
	wg := &sync.WaitGroup{}
	endpoint := fmt.Sprintf("http://%s/messages", e.host)
	log.Printf(`endpoint is: %s`, endpoint)
	for _, data := range e.data {
		cancel := make(chan struct{})
		body, err := json.Marshal(data.meta)
		if err != nil {
			return fmt.Errorf(`failed to message: %s`, err)
		}
		for _, ts := range data.offsets {
			go e.notify(cancel, wg, ts, body, endpoint)
		}
	}
	wg.Wait()
	return nil
}

func (e Engine) notify(cancel chan struct{}, wg *sync.WaitGroup, ts time.Duration, body []byte, endpoint string) {
	wg.Add(1)
	defer wg.Done()
	timer := time.NewTimer(ts)
	select {
	case _, ok := <-timer.C:
		if ok {
			log.Printf("timer fires event: %s", body)
			buf := bytes.NewBuffer(body)
			response, err := http.Post(endpoint, CTYPE_JSON, buf)
			if err != nil {
				log.Printf(`failed to perform request: %s`, err)
				return
			}
			if response.StatusCode != http.StatusCreated {
				log.Printf(`commservice responded with an unexpected status code: %d; status: %s`, response.StatusCode, response.Status)
				return
			}
			defer response.Body.Close()
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Printf(`failed to read response body: %s`, err)
				return
			}
			b := &messagesResponse{}
			json.Unmarshal(responseBody, b)
			log.Printf("response body: %s", responseBody)
			if b.Paid == true {
				close(cancel)
			}
		}
	case <-cancel:
		timer.Stop()
		log.Printf(`notification cancelled %s`, body)
	}
}
