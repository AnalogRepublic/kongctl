package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKong(t *testing.T) {
	kong, err := NewKong("http://testing.kongctl.io:8001")

	assert.Nil(t, err)
	assert.NotNil(t, kong)
}

func TestPluginsList(t *testing.T) {
	kong, err := NewKong("http://testing.kongctl.io:8001")

	plugins, err := kong.Plugins().List()

	assert.Nil(t, err)
	assert.NotNil(t, plugins)
}
