package nginx

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/// TestMarkers tests marker replacements for string and splitted strings
func TestMarkers(t *testing.T) {
	markers := map[string]interface{}{
		"foo":    "bar",
		"fooArr": "bar1,bar2,bar3",
	}
	markersSplit := map[string]interface{}{
		"fooArr": ",",
	}
	processedMarkers := ProcessMarkers(markers, markersSplit)

	assert.Equal(t, markers["foo"], processedMarkers["foo"])

	// Assert array was processed
	for i, v := range strings.Split(markers["fooArr"].(string), markersSplit["fooArr"].(string)) {
		assert.Equal(t, v, processedMarkers[fmt.Sprintf("fooArr[%d]", i)])
	}

	text := ReplaceMarkers("Marker is {# foo #} with {# fooArr[1] #}", processedMarkers)
	assert.Equal(t, text, "Marker is bar with bar2")
	text = ReplaceMarkers("Marker is {~ foo ~} with {~ fooArr[1] ~}", processedMarkers)
	assert.Equal(t, text, "Marker is bar with bar2")
	text = ReplaceMarkers("Marker is {* foo *} with {* fooArr[1] *}", processedMarkers)
	assert.Equal(t, text, "Marker is bar with bar2")
	text = ReplaceMarkers("Marker is {# foo #} with {# fooArr[1] #}", processedMarkers)
	assert.Equal(t, text, "Marker is bar with bar2")
}
