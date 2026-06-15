pragma Singleton
pragma ComponentBehavior: Bound

import Quickshell
import qs.Services

Singleton {
    id: root

    function connectToNetwork(network, options) {
        if (!network)
            return;

        const actionOptions = options || {};
        const ssid = network.ssid || "";
        if (!ssid)
            return;

        const connected = actionOptions.connected ?? network.connected ?? (ssid === NetworkService.currentWifiSSID);
        if (connected) {
            if (actionOptions.disconnectWhenConnected ?? false)
                NetworkService.disconnectWifi();
            return;
        }

        if (needsCredentialWindow(network)) {
            PopoutService.showWifiPasswordModal(ssid);
            return;
        }

        NetworkService.connectToWifi(ssid);
    }

    function connectToNetworkFromDetails(ssid, secured, saved, enterprise, connected, options) {
        connectToNetwork({
            ssid: ssid,
            secured: secured,
            saved: saved,
            enterprise: enterprise,
            connected: connected
        }, options);
    }

    function editSavedNetwork(network) {
        if (!network)
            return;

        const ssid = network.ssid || "";
        if (!ssid)
            return;
        if (!canEditCredentials(network))
            return;

        PopoutService.showWifiEditModalForNetwork(network);
    }

    function editSavedNetworkFromDetails(ssid, secured, saved, enterprise) {
        editSavedNetwork({
            ssid: ssid,
            secured: secured,
            saved: saved,
            enterprise: enterprise
        });
    }

    function shouldPromptForCredentials(network) {
        return (network.secured ?? false) && !(network.saved ?? false);
    }

    function needsCredentialWindow(network) {
        return shouldPromptForCredentials(network);
    }

    function canEditCredentials(network) {
        return NetworkService.canEditWifiCredentials(network);
    }

    function canToggleAutoconnect(network) {
        return (network?.saved ?? false) && NetworkService.supportsWifiAutoconnect;
    }

    function canConnectHiddenNetwork() {
        return NetworkService.supportsHiddenWifiConnect;
    }
}
