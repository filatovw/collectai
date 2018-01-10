package reminder_test

import (
	"os"
	"strings"
	"testing"

	"github.com/filatovw/collectai/app/reminder"
	"github.com/stretchr/testify/assert"
)

var data = []struct {
	Key   string
	Value string
}{
	{Key: "commservice-host", Value: "somehost"},
	{Key: "schedule-path", Value: "somepath"},
	{Key: "engine", Value: "schedulerengine"},
}

func TestConfCLI(t *testing.T) {
	conf := &reminder.Conf{}

	os.Args = []string{
		"", "-commservice-host=somehost", "-schedule-path=somepath", "-engine=schedulerengine",
	}
	conf.ParseCLI()

	assert.Equal(t, data[0].Value, conf.CommserviceHost)
	assert.Equal(t, data[1].Value, conf.SchedulePath)
	assert.Equal(t, data[2].Value, conf.Engine)
}

func TestConfEnv(t *testing.T) {
	conf := &reminder.Conf{}

	for _, item := range data {
		key := strings.Replace(item.Key, "-", "_", -1)
		key = strings.ToUpper(key)
		os.Setenv(key, item.Value)
	}
	conf.ParseEnv()

	assert.Equal(t, data[0].Value, conf.CommserviceHost)
	assert.Equal(t, data[1].Value, conf.SchedulePath)
	assert.Equal(t, data[2].Value, conf.Engine)
}
