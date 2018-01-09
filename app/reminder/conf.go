package reminder

import (
	"errors"
	"flag"
	"os"

	"github.com/filatovw/collectai/app/reminder/engine"
)

type Conf struct {
	CommserviceHost string
	SchedulePath    string
	Engine          string
}

func (c Conf) isReady() bool {
	hasEngine := false
	if c.Engine == engine.ENGINE_GREEDY || c.Engine == engine.ENGINE_POOL {
		hasEngine = true
	}
	return (c.CommserviceHost != "" &&
		c.SchedulePath != "" &&
		hasEngine)
}

// ParseCLI update from command line arguments
func (c *Conf) ParseCLI() {
	conf := &Conf{}
	flag.StringVar(&conf.CommserviceHost, "commservice-host", "", "address of commservice [host]:[port].")
	flag.StringVar(&conf.SchedulePath, "schedule-path", "", "path to schedule file CSV formatted")
	flag.StringVar(&conf.Engine, "engine", "", "scheduler engine. Greedy(default) or Pool")
	flag.Parse()

	if conf.CommserviceHost != "" {
		c.CommserviceHost = conf.CommserviceHost
	}
	if conf.SchedulePath != "" {
		c.SchedulePath = conf.SchedulePath
	}
	if conf.Engine != "" {
		c.Engine = conf.Engine
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

	v = os.Getenv("ENGINE")
	if v != "" {
		c.Engine = v
	}
}

// NewConf creates Conf, that filled from Environment variables and CLI
func NewConf() (*Conf, error) {
	conf := &Conf{
		Engine: engine.ENGINE_GREEDY,
	}
	conf.ParseEnv()
	conf.ParseCLI()
	if !conf.isReady() {
		flag.PrintDefaults()
		return nil, errors.New(`failed to configure`)
	}
	return conf, nil
}
