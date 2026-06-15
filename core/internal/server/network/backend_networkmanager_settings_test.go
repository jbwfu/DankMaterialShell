package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeNetworkManagerLegacyIPSettingsForUpdate(t *testing.T) {
	settings := map[string]map[string]any{
		"ipv4": {
			"method":       "auto",
			"addresses":    []any{},
			"routes":       []any{},
			"dns":          []any{},
			"address-data": []map[string]any{{"address": "192.0.2.10", "prefix": uint32(24)}},
			"dns-search":   []string{"example.com"},
		},
		"ipv6": {
			"method":        "auto",
			"addresses":     []any{},
			"routes":        []any{},
			"dns":           []any{},
			"addr-gen-mode": int32(1),
			"route-data":    []map[string]any{{"dest": "::", "prefix": uint32(0)}},
		},
	}

	sanitizeNetworkManagerLegacyIPSettingsForUpdate(settings)

	assert.NotContains(t, settings["ipv4"], "addresses")
	assert.NotContains(t, settings["ipv4"], "routes")
	assert.NotContains(t, settings["ipv4"], "dns")
	assert.Equal(t, "auto", settings["ipv4"]["method"])
	assert.Equal(t, []map[string]any{{"address": "192.0.2.10", "prefix": uint32(24)}}, settings["ipv4"]["address-data"])
	assert.Equal(t, []string{"example.com"}, settings["ipv4"]["dns-search"])

	assert.NotContains(t, settings["ipv6"], "addresses")
	assert.NotContains(t, settings["ipv6"], "routes")
	assert.NotContains(t, settings["ipv6"], "dns")
	assert.Equal(t, "auto", settings["ipv6"]["method"])
	assert.Equal(t, int32(1), settings["ipv6"]["addr-gen-mode"])
	assert.Equal(t, []map[string]any{{"dest": "::", "prefix": uint32(0)}}, settings["ipv6"]["route-data"])
}
