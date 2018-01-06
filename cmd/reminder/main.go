package main

import (
	"log"
	"os"

	"github.com/filatovw/collectai/app/reminder"
)

func main() {
	app, err := reminder.NewApp()
	if err != nil {
		log.Fatalf(`failed to start application`, err)
	}
	os.Exit(app.Run())
}
