package network

import (
	"testing"

	"github.com/AvengeMedia/DankMaterialShell/core/internal/errdefs"
	mock_gonetworkmanager "github.com/AvengeMedia/DankMaterialShell/core/internal/mocks/github.com/Wifx/gonetworkmanager/v2"
	"github.com/Wifx/gonetworkmanager/v2"
	"github.com/stretchr/testify/assert"
)

func TestNetworkManagerBackend_UpdatePrimaryConnection(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	mockNM.EXPECT().GetPropertyActiveConnections().Return([]gonetworkmanager.ActiveConnection{}, nil)
	mockNM.EXPECT().GetPropertyPrimaryConnection().Return(nil, nil)

	err = backend.updatePrimaryConnection()
	assert.NoError(t, err)
}

func TestNetworkManagerBackend_UpdateEthernetState_NoDevice(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.ethernetDevice = nil
	err = backend.updateEthernetState()
	assert.NoError(t, err)
}

func TestNetworkManagerBackend_UpdateWiFiState_NoDevice(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	backend.wifiDevice = nil
	err = backend.updateWiFiState()
	assert.NoError(t, err)
}

func TestNetworkManagerBackend_ClassifyNMStateReason(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	testCases := []struct {
		reason   uint32
		expected string
	}{
		{NmDeviceStateReasonWrongPassword, errdefs.ErrBadCredentials},
		{NmDeviceStateReasonNoSecrets, errdefs.ErrUserCanceled},
		{NmDeviceStateReasonSupplicantTimeout, errdefs.ErrBadCredentials},
		{NmDeviceStateReasonDhcpClientFailed, errdefs.ErrDhcpTimeout},
		{NmDeviceStateReasonNoSsid, errdefs.ErrNoSuchSSID},
		{999, errdefs.ErrConnectionFailed},
	}

	for _, tc := range testCases {
		result := backend.classifyNMStateReason(tc.reason)
		assert.Equal(t, tc.expected, result)
	}
}

func TestShouldForgetFailedWiFiProfile(t *testing.T) {
	testCases := []struct {
		name        string
		reasonCode  string
		preExisting bool
		expected    bool
	}{
		{"new profile bad credentials", errdefs.ErrBadCredentials, false, true},
		{"saved profile bad credentials", errdefs.ErrBadCredentials, true, false},
		{"new profile missing SSID", errdefs.ErrNoSuchSSID, false, true},
		{"saved profile missing SSID", errdefs.ErrNoSuchSSID, true, false},
		{"user cancel always clears", errdefs.ErrUserCanceled, true, true},
		{"dhcp failure keeps new profile", errdefs.ErrDhcpTimeout, false, false},
		{"generic failure keeps new profile", errdefs.ErrConnectionFailed, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, shouldForgetFailedWiFiProfile(tc.reasonCode, tc.preExisting))
		})
	}
}

func TestNetworkManagerBackend_GetDeviceIP_NoConfig(t *testing.T) {
	mockNM := mock_gonetworkmanager.NewMockNetworkManager(t)
	mockDevice := mock_gonetworkmanager.NewMockDevice(t)

	backend, err := NewNetworkManagerBackend(mockNM)
	assert.NoError(t, err)

	mockDevice.EXPECT().GetPropertyIP4Config().Return(nil, nil)

	ip := backend.getDeviceIP(mockDevice)
	assert.Empty(t, ip)
}
