package reminder

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/filatovw/collectai/app/reminder/engine"
)

// App application object
type App struct {
	conf      *Conf
	schedule  [][]string
	scheduler engine.Engine
}

// Run start application
func (a *App) Run() (exitcode int) {
	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
			exitcode = 2
		}
	}()

	// read schedule
	err := a.readCSV()
	if err != nil {
		log.Printf(`failed to read CSV: %s`, err)
		return 1
	}
	// init scheduler
	if err = a.scheduler.Init(a.schedule, a.conf.CommserviceHost); err != nil {
		log.Printf(`failed to Init engine: %s`, err)
		return 1
	}
	// process notification
	if err = a.scheduler.Process(); err != nil {
		log.Printf(`failed to Init engine: %s`, err)
		return 1
	}
	return 0
}

// readCSV read csv file
func (a *App) readCSV() error {
	csvfile, err := os.Open(a.conf.SchedulePath)
	if err != nil {
		log.Printf(`failed to open file with the schedule: %s`, err)
	}
	defer csvfile.Close()

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
	return fmt.Errorf(`empty CSV file: %s`, a.conf.SchedulePath)
}

// NewApp configures application
func NewApp(conf *Conf) (*App, error) {
	engine, err := engine.GetEngine(conf.Engine)
	if err != nil {
		return nil, fmt.Errorf(`failed to create Application: %s`, err)
	}
	return &App{
		conf:      conf,
		scheduler: engine,
	}, nil
}
