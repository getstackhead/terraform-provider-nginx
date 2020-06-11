package nginx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/// ProcessMarkers resolves array values into single string replaces
func TestProcessMarkers(t *testing.T) {
	markers := map[string]interface{}{
		"foo":    "bar",
		"fooArr": []string{"bar1", "bar2", "bar3"},
	}
	processedMarkers := ProcessMarkers(markers)

	assert.Equal(t, markers["foo"], processedMarkers["foo"])

	// Assert array was processed
	for i, v := range markers["fooArr"].([]string) {
		assert.Equal(t, v, processedMarkers[fmt.Sprintf("fooArr[%d]", i)])
	}
}
