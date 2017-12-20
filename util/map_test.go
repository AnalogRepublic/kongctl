package util

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestDiffMaps(t *testing.T) {
	var exists bool

	a := map[interface{}]interface{}{
		"testing": "hello1",
		"another": "world",
	}

	b := map[interface{}]interface{}{
		"testing":   "hello2",
		"different": "world",
		"look":      "old",
	}

	result := DiffMapsKeys(a, b)

	_, exists = result.Additions["another"]
	assert.Equal(t, exists, true, "Should see the 'another' key in the additions")

	_, exists = result.Deletions["different"]
	assert.Equal(t, exists, true, "Should see the 'different' key in the deletions")

	_, exists = result.Deletions["look"]
	assert.Equal(t, exists, true, "Should see the 'look' key in the deletions")

	_, exists = result.Updates["testing"]
	assert.Equal(t, exists, true, "Should see the 'testing' key in the updates")
}
