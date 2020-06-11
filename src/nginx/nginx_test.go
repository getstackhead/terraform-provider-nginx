package nginx

import (
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

	text := ReplaceMarkers("Marker is {# foo #} with {# fooArr[1] #}", processedMarkers)
	assert.Equal(t, "Marker is bar with bar2", text)
	text = ReplaceMarkers("Marker is {~ foo ~} with {~ fooArr[1] ~}", processedMarkers)
	assert.Equal(t, "Marker is bar with bar2", text)
	text = ReplaceMarkers("Marker is {* foo *} with {* fooArr[1] *}", processedMarkers)
	assert.Equal(t, "Marker is bar with bar2", text)
	text = ReplaceMarkers("Marker is {# foo #} with {# fooArr[99] #}", processedMarkers)
	assert.Equal(t, "Marker is bar with ", text)
}
