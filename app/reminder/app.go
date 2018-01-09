package reminder

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/filatovw/collectai/app/reminder/engine"
)

type App struct {
	Conf      *Conf
	schedule  [][]string
	scheduler engine.Engine
}

func (a *App) Run() (exitcode int) {
	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
			exitcode = 2
		}
	}()
	// TODO: init ctrl+v cancellation

	// read schedule
	err := a.readCSV()
	if err != nil {
		log.Printf(`failed to read CSV: %s`, err)
		return 1
	}
	if err := a.scheduler.Init(a.schedule, a.Conf.CommserviceHost); err != nil {
		log.Printf(`failed to Init engine: %s`, err)
	}
	cancel := make(chan struct{})
	a.scheduler.Process(cancel)

	return 0
}

func (a *App) readCSV() error {
	csvfile, err := os.Open(a.Conf.SchedulePath)
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
		return err
	}
	if len(all) > 1 {
		a.schedule = all[1:]
		return nil
	}
	return fmt.Errorf(`empty CSV file: %s`, a.Conf.SchedulePath)
}

func NewApp(conf *Conf) (*App, error) {
	engine, err := engine.GetEngine(conf.Engine)
	if err != nil {
		return nil, fmt.Errorf(`failed to create Application: %s`, err)
	}
	return &App{
		Conf:      conf,
		scheduler: engine,
	}, nil
}
