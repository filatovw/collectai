package engine

import (
	"fmt"

	"github.com/filatovw/collectai/app/reminder/engine/greedy"
	"github.com/filatovw/collectai/app/reminder/engine/pool"
)

// scheduler engines
const (
	// `pool` engine creates a pool of goroutines. Helps to utilize resources much more efficiently.
	ENGINE_POOL = "pool"
	// `greedy` engine creates number of goroutines equals to the number of offsets from the `schedule` column.
	ENGINE_GREEDY = "greedy"
)

// Engine interface for a scheduler
type Engine interface {
	Init([][]string, string) error
	Process() error
}

// GetEngine gives scheduler engine
func GetEngine(engine string) (Engine, error) {
	switch engine {
	case ENGINE_POOL:
		return &pool.Engine{}, nil
	case ENGINE_GREEDY:
		return &greedy.Engine{}, nil
	}
	return nil, fmt.Errorf(`failed to find engine: %s`, engine)
}
