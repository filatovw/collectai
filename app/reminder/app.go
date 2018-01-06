package reminder

import (
	"encoding/csv"
	"log"
	"os"
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
	log.Printf(`%#v`, v)
	return 0
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

	var rows [][]string

	for {
		row, err := r.Read()
		if err != nil {
			break
		}
		rows = append(rows, row)
	}
	return rows, err
}

func NewApp(conf *Conf) (*App, error) {
	return &App{
		Conf: conf,
	}, nil
}
