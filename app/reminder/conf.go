package reminder

import (
	"errors"
	"flag"
	"os"
)

type Conf struct {
	CommserviceHost string
	SchedulePath    string
}

func (c Conf) isReady() bool {
	return (c.CommserviceHost != "" &&
		c.SchedulePath != "")
}

// ParseCLI update from command line arguments
func (c *Conf) ParseCLI() {
	conf := &Conf{}
	flag.StringVar(&conf.CommserviceHost, "commservice-host", "", "address of commservice [host]:[port].")
	flag.StringVar(&conf.SchedulePath, "schedule-path", "", "path to schedule file CSV formatted")
	flag.Parse()

	if conf.CommserviceHost != "" {
		c.CommserviceHost = conf.CommserviceHost
	}
	if conf.SchedulePath != "" {
		c.SchedulePath = conf.SchedulePath
	}
}

// ParseEnv update from Environment variables
func (c *Conf) ParseEnv() {
	v := os.Getenv("COMMSERVICE_HOST")
	if v != "" {
		c.CommserviceHost = v
	}

	v = os.Getenv("SCHEDULE_PATH")
	if v != "" {
		c.SchedulePath = v
	}
}

// NewConf creates Conf, that filled from Environment variables and CLI
func NewConf() (*Conf, error) {
	conf := &Conf{}
	conf.ParseEnv()
	conf.ParseCLI()
	if !conf.isReady() {
		flag.PrintDefaults()
		return nil, errors.New(`failed to configure`)
	}
	return conf, nil
}
