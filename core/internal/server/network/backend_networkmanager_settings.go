package network

func sanitizeNetworkManagerLegacyIPSettingsForUpdate(settings map[string]map[string]any) {
	for _, sectionName := range []string{"ipv4", "ipv6"} {
		section, ok := settings[sectionName]
		if !ok {
			continue
		}
		delete(section, "addresses")
		delete(section, "routes")
		delete(section, "dns")
	}
}
