package reminder

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type App struct {
	Conf *Conf
}

func (a *App) Run() int {
	// read schedule
	v, err := readCSV(a.Conf.SchedulePath)
	if err != nil {
		log.Printf(`failed to read CSV: %s`, err)
	}
	a.setupScheduler(v)
	return 0
}

type Person struct {
	Email string `json:"email"`
	Text  string `json:"text"`
}

func (a App) setupScheduler(data [][]string) {
	wg := sync.WaitGroup{}
	for _, item := range data {
		timestamps := strings.Split(item[2], "-")
		stop := make(chan struct{})
		person := &Person{
			Email: item[0],
			Text:  item[1],
		}
		for _, ts := range timestamps {
			dur, err := time.ParseDuration(ts)
			if err != nil {
				log.Printf("%v %v", err, ts)
			}
			timer := time.NewTimer(dur)

			wg.Add(1)
			go func(email, message string, timer *time.Timer, stop chan struct{}) {
				defer wg.Done()
				select {
				case <-timer.C:
					log.Printf("notified person:%s;  message: %s", email, message)
					body, err := json.Marshal(person)
					if err != nil {
						log.Printf(`failed to serialize: %s`, person)
						return
					}
					resp, err := http.Post("http://"+a.Conf.CommserviceHost+"/messages", "application/json", bytes.NewBuffer(body))
					if err != nil {
						log.Printf(`failed to read response: %s`, err)
						return
					}
					defer resp.Body.Close()
					body, err = ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Printf(`failed to read response body: %s`, err)
					}
					var vv interface{}
					json.Unmarshal(body, vv)
					log.Printf("body :%v", vv)

				case <-stop:
					log.Printf("stopped person: %s", email)
					timer.Stop()
				}
			}(item[0], item[1], timer, stop)
		}
	}
	wg.Wait()
}

func readCSV(path string) ([][]string, error) {
	csvfile, err := os.Open(path)
	if err != nil {
		log.Printf(`failed to open file with the schedule: %s`, err)
	}
	r := csv.NewReader(csvfile)
	r.Comma = ','
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	r.Comment = '#'
	all, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(all) > 1 {
		return all[1:], nil
	}
	return nil, fmt.Errorf(`empty CSV file: %s`, path)
}

func NewApp(conf *Conf) (*App, error) {
	return &App{
		Conf: conf,
	}, nil
}
