package outputstdout

import (
	"github.com/stretchr/testify/assert"
	"github.com/tsaikd/gogstash/config"
	"os"
	"testing"
	"time"
)

func Test_main(t *testing.T) {

	assert := assert.New(t)

	conftest, err := config.LoadConfig("config_test.json")
	assert.NoError(err)

	outputs := conftest.Output()

	assert.Len(outputs, 1)
	if len(outputs) > 0 {
		output := outputs[0].(*OutputConfig)

		assert.IsType(&OutputConfig{}, output)
		assert.Equal("file", output.GetType())

		err = output.Event(config.LogEvent{
			Timestamp: time.Now(),
			Message:   "output file test message",
		})
		assert.NoError(err)
		_, err := os.Open(output.Path)

		assert.NoError(err)
	}
}
