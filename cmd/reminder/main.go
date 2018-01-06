package main

import (
	"log"
	"os"

	"github.com/filatovw/collectai/app/reminder"
)

func main() {
	// read configuration
	conf, err := reminder.NewConf()
	if err != nil {
		log.Fatal(err)
	}
	// create app
	app, err := reminder.NewApp(conf)
	if err != nil {
		log.Fatalf(`failed to start application: %s`, err)
	}
	// run app
	os.Exit(app.Run())
}
