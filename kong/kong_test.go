package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKong(t *testing.T) {
	kong, err := NewKong("http://127.0.0.1:8001")

	assert.Nil(t, err)
	assert.NotNil(t, kong)
}

func TestPingFail(t *testing.T) {
	_, err := NewKong("http://asdmhasduaishjdoasndasgdioasjdghvjk")
	assert.NotNil(t, err)
}

func TestPluginsList(t *testing.T) {
	kong, err := NewKong("http://127.0.0.1:8001")

	plugins, err := kong.Plugins().List()

	assert.Nil(t, err)
	assert.NotNil(t, plugins)
}
