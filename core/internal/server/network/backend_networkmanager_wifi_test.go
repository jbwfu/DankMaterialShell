package network

import (
	"testing"

	mock_gonetworkmanager "github.com/AvengeMedia/DankMaterialShell/core/internal/mocks/github.com/Wifx/gonetworkmanager/v2"
	"github.com/Wifx/gonetworkmanager/v2"
	"github.com/stretchr/testify/assert"
)

func TestNetworkManagerBackend_GetWiFiEnabled(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	mockNM.EXPECT().GetPropertyWirelessEnabled().Return(true, nil)

	enabled, err := backend.GetWiFiEnabled()
	assert.NoError(t, err)
	assert.True(t, enabled)
}

func TestNetworkManagerBackend_SetWiFiEnabled(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	mockNM.EXPECT().SetPropertyWirelessEnabled(true).Return(nil)

	err = backend.SetWiFiEnabled(true)
	assert.NoError(t, err)

	backend.stateMutex.RLock()
	assert.True(t, backend.state.WiFiEnabled)
	backend.stateMutex.RUnlock()
}

func TestNetworkManagerBackend_ScanWiFi_NoDevice(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.wifiDevice = nil
	err = backend.ScanWiFi()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no WiFi device available")
}

func TestNetworkManagerBackend_ScanWiFi_Disabled(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)
	mockDeviceWireless := mock_gonetworkmanager.NewMockDeviceWireless(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.wifiDevice = mockDeviceWireless
	backend.wifiDev = mockDeviceWireless

	backend.stateMutex.Lock()
	backend.state.WiFiEnabled = false
	backend.stateMutex.Unlock()

	err = backend.ScanWiFi()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "WiFi is disabled")
}

func TestNetworkManagerBackend_GetWiFiNetworkDetails_NoDevice(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.wifiDevice = nil
	_, err = backend.GetWiFiNetworkDetails("TestNetwork")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no WiFi device available")
}

func TestNetworkManagerBackend_ConnectWiFi_NoDevice(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.wifiDevice = nil
	req := ConnectionRequest{SSID: "TestNetwork", Password: "password"}
	err = backend.ConnectWiFi(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no WiFi device available")
}

func TestNetworkManagerBackend_ConnectWiFi_AlreadyConnected(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)
	mockDeviceWireless := mock_gonetworkmanager.NewMockDeviceWireless(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.wifiDevice = mockDeviceWireless
	backend.wifiDev = mockDeviceWireless
	backend.wifiDevices = map[string]*wifiDeviceInfo{
		"wlan0": {
			device:    nil,
			wireless:  mockDeviceWireless,
			name:      "wlan0",
			hwAddress: "00:11:22:33:44:55",
		},
	}

	mockDeviceWireless.EXPECT().GetPropertyInterface().Return("wlan0", nil)

	backend.stateMutex.Lock()
	backend.state.WiFiConnected = true
	backend.state.WiFiSSID = "TestNetwork"
	backend.state.WiFiDevice = "wlan0"
	backend.stateMutex.Unlock()

	req := ConnectionRequest{SSID: "TestNetwork", Password: "password"}
	err = backend.ConnectWiFi(req)
	assert.NoError(t, err)
}

func TestNetworkManagerBackend_DisconnectWiFi_NoDevice(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.wifiDevice = nil
	err = backend.DisconnectWiFi()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no WiFi device available")
}

func TestNetworkManagerBackend_IsConnectingTo(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.stateMutex.Lock()
	backend.state.IsConnecting = true
	backend.state.ConnectingSSID = "TestNetwork"
	backend.stateMutex.Unlock()

	assert.True(t, backend.IsConnectingTo("TestNetwork"))
	assert.False(t, backend.IsConnectingTo("OtherNetwork"))
}

func TestNetworkManagerBackend_IsConnectingTo_NotConnecting(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.stateMutex.Lock()
	backend.state.IsConnecting = false
	backend.state.ConnectingSSID = ""
	backend.stateMutex.Unlock()

	assert.False(t, backend.IsConnectingTo("TestNetwork"))
}

func TestApplyNetworkManagerWiFiCredentials_PSK(t *testing.T) {
	password := "new-secret"
	save := true
	settings := map[string]map[string]any{
		"802-11-wireless-security": {
			"key-mgmt":  "wpa-psk",
			"psk-flags": uint32(2),
		},
	}

	err := applyNetworkManagerWiFiCredentials(settings, &WiFiCredentialsConfig{
		Password: &password,
		Save:     &save,
	})

	assert.NoError(t, err)
	assert.Equal(t, password, settings["802-11-wireless-security"]["psk"])
	assert.Equal(t, uint32(0), settings["802-11-wireless-security"]["psk-flags"])
}

func TestApplyNetworkManagerWiFiCredentials_PSKRejectsInvalidPassword(t *testing.T) {
	password := "short"
	settings := map[string]map[string]any{
		"802-11-wireless-security": {
			"key-mgmt":  "wpa-psk",
			"psk":       "old-secret",
			"psk-flags": uint32(0),
		},
	}

	err := applyNetworkManagerWiFiCredentials(settings, &WiFiCredentialsConfig{
		Password: &password,
	})

	assert.ErrorContains(t, err, "WPA password must be 8-63 characters")
	assert.Equal(t, "old-secret", settings["802-11-wireless-security"]["psk"])
}

func TestApplyNetworkManagerWiFiCredentials_PSKAccepts64HexKey(t *testing.T) {
	password := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	settings := map[string]map[string]any{
		"802-11-wireless-security": {
			"key-mgmt":  "wpa-psk",
			"psk-flags": uint32(2),
		},
	}

	err := applyNetworkManagerWiFiCredentials(settings, &WiFiCredentialsConfig{
		Password: &password,
	})

	assert.NoError(t, err)
	assert.Equal(t, password, settings["802-11-wireless-security"]["psk"])
}

func TestApplyNetworkManagerWiFiCredentials_PSKRejects64NonHexPassword(t *testing.T) {
	password := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
	settings := map[string]map[string]any{
		"802-11-wireless-security": {
			"key-mgmt": "wpa-psk",
		},
	}

	err := applyNetworkManagerWiFiCredentials(settings, &WiFiCredentialsConfig{
		Password: &password,
	})

	assert.ErrorContains(t, err, "WPA password must be 8-63 characters")
	assert.NotContains(t, settings["802-11-wireless-security"], "psk")
}

func TestApplyNetworkManagerWiFiCredentials_SAEAllowsShortPassword(t *testing.T) {
	password := "short"
	settings := map[string]map[string]any{
		"802-11-wireless-security": {
			"key-mgmt":  "sae",
			"psk-flags": uint32(2),
		},
	}

	err := applyNetworkManagerWiFiCredentials(settings, &WiFiCredentialsConfig{
		Password: &password,
	})

	assert.NoError(t, err)
	assert.Equal(t, password, settings["802-11-wireless-security"]["psk"])
}

func TestApplyNetworkManagerWiFiCredentials_Enterprise(t *testing.T) {
	username := "alice"
	password := "new-secret"
	anon := "anon"
	domain := "example.com"
	save := true
	settings := map[string]map[string]any{
		"802-11-wireless-security": {
			"key-mgmt": "wpa-eap",
		},
		"802-1x": {
			"eap": []string{"peap"},
		},
	}

	err := applyNetworkManagerWiFiCredentials(settings, &WiFiCredentialsConfig{
		Username:          &username,
		Password:          &password,
		AnonymousIdentity: &anon,
		DomainSuffixMatch: &domain,
		Save:              &save,
	})

	assert.NoError(t, err)
	assert.Equal(t, username, settings["802-1x"]["identity"])
	assert.Equal(t, password, settings["802-1x"]["password"])
	assert.Equal(t, anon, settings["802-1x"]["anonymous-identity"])
	assert.Equal(t, domain, settings["802-1x"]["domain-suffix-match"])
	assert.Equal(t, uint32(0), settings["802-1x"]["password-flags"])
}

func TestNetworkManagerBackend_UpdateWiFiNetworks_NoDevice(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.wifiDevice = nil
	_, err = backend.updateWiFiNetworks()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no WiFi device available")
}

func TestNetworkManagerBackend_UpdateSavedWiFiNetworksPreservesVisibleSavedNetworks(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)
	mockSettings := mock_gonetworkmanager.NewMockSettings(t)
	mockConn := mock_gonetworkmanager.NewMockConnection(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)
	backend.settings = mockSettings

	backend.stateMutex.Lock()
	backend.state.WiFiNetworks = []WiFiNetwork{
		{
			SSID:   "Home",
			Signal: 76,
		},
	}
	backend.stateMutex.Unlock()

	settings := gonetworkmanager.ConnectionSettings{
		"connection": {
			"type":        "802-11-wireless",
			"autoconnect": true,
		},
		"802-11-wireless": {
			"ssid": []byte("Home"),
		},
		"802-11-wireless-security": {},
	}
	mockSettings.EXPECT().ListConnections().Return([]gonetworkmanager.Connection{mockConn}, nil)
	mockConn.EXPECT().GetSettings().Return(settings, nil)

	err = backend.updateSavedWiFiNetworks()
	assert.NoError(t, err)

	backend.stateMutex.RLock()
	savedNetworks := append([]WiFiNetwork(nil), backend.state.SavedWiFiNetworks...)
	wifiNetworks := append([]WiFiNetwork(nil), backend.state.WiFiNetworks...)
	backend.stateMutex.RUnlock()

	assert.Len(t, wifiNetworks, 1)
	assert.True(t, wifiNetworks[0].Saved)
	assert.Len(t, savedNetworks, 1)
	assert.Equal(t, "Home", savedNetworks[0].SSID)
	assert.True(t, savedNetworks[0].Saved)
	assert.False(t, savedNetworks[0].OutOfRange)
	assert.Equal(t, uint8(76), savedNetworks[0].Signal)
}

func TestNetworkManagerBackend_FindConnection_NoSettings(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.settings = nil
	_, err = backend.findConnection("NonExistentNetwork")
	assert.Error(t, err)
}

func TestNetworkManagerBackend_CreateAndConnectWiFi_NoDevice(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.wifiDevice = nil
	backend.wifiDev = nil
	req := ConnectionRequest{SSID: "TestNetwork", Password: "password"}
	err = backend.createAndConnectWiFi(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no WiFi device available")
}
