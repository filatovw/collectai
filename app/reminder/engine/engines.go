package engine

import (
	"fmt"

	"github.com/filatovw/collectai/app/reminder/engine/greedy"
	"github.com/filatovw/collectai/app/reminder/engine/pool"
)

// scheduler engines
const (
	ENGINE_POOL   = "pool"
	ENGINE_GREEDY = "greedy"
)

type Engine interface {
	Init([][]string, string) error
	Process(chan<- struct{}) error
}

func GetEngine(engine string) (Engine, error) {
	switch engine {
	case ENGINE_POOL:
		return &pool.Engine{}, nil
	case ENGINE_GREEDY:
		return &greedy.Engine{}, nil
	}
	return nil, fmt.Errorf(`failed to find engine: %s`, engine)
}
