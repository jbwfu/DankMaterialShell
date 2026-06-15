import QtQuick
import QtQuick.Controls
import qs.Common
import qs.Services
import qs.Widgets

Menu {
    id: root

    width: 170
    closePolicy: Popup.CloseOnEscape | Popup.CloseOnPressOutsideParent

    property string currentSSID: ""
    property bool currentSecured: false
    property bool currentEnterprise: false
    property bool currentConnected: false
    property bool currentSaved: false
    property int currentSignal: 0
    property bool currentAutoconnect: false
    property bool currentOutOfRange: false
    property bool showNetworkInfoAction: false

    readonly property var currentNetwork: ({
        ssid: currentSSID,
        secured: currentSecured,
        enterprise: currentEnterprise,
        connected: currentConnected,
        saved: currentSaved,
        signal: currentSignal,
        autoconnect: currentAutoconnect,
        outOfRange: currentOutOfRange
    })
    readonly property bool showSavedOptions: currentSaved || currentConnected
    readonly property bool canEdit: WifiConnectionActions.canEditCredentials(currentNetwork)
    readonly property bool canToggleAutoconnect: WifiConnectionActions.canToggleAutoconnect(currentNetwork)

    signal networkInfoRequested(string ssid)
    signal forgetRequested(string ssid)

    background: Rectangle {
        color: Theme.withAlpha(Theme.surfaceContainer, Theme.popupTransparency)
        radius: Theme.cornerRadius
        border.width: 0
        border.color: Qt.rgba(Theme.outline.r, Theme.outline.g, Theme.outline.b, 0.12)
    }

    MenuItem {
        text: root.currentConnected ? I18n.tr("Disconnect") : I18n.tr("Connect")
        height: root.currentOutOfRange ? 0 : 32
        visible: !root.currentOutOfRange

        contentItem: StyledText {
            text: parent.text
            font.pixelSize: Theme.fontSizeSmall
            color: Theme.surfaceText
            leftPadding: Theme.spacingS
            verticalAlignment: Text.AlignVCenter
        }

        background: Rectangle {
            color: parent.hovered ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.08) : "transparent"
            radius: Theme.cornerRadius / 2
        }

        onTriggered: {
            WifiConnectionActions.connectToNetworkFromDetails(root.currentSSID, root.currentSecured, root.currentSaved, root.currentEnterprise, root.currentConnected, {
                disconnectWhenConnected: true
            });
        }
    }

    MenuItem {
        text: I18n.tr("Network Info")
        height: root.showNetworkInfoAction ? 32 : 0
        visible: root.showNetworkInfoAction

        contentItem: StyledText {
            text: parent.text
            font.pixelSize: Theme.fontSizeSmall
            color: Theme.surfaceText
            leftPadding: Theme.spacingS
            verticalAlignment: Text.AlignVCenter
        }

        background: Rectangle {
            color: parent.hovered ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.08) : "transparent"
            radius: Theme.cornerRadius / 2
        }

        onTriggered: root.networkInfoRequested(root.currentSSID)
    }

    MenuItem {
        text: I18n.tr("Edit")
        height: root.canEdit ? 32 : 0
        visible: root.canEdit

        contentItem: StyledText {
            text: parent.text
            font.pixelSize: Theme.fontSizeSmall
            color: Theme.surfaceText
            leftPadding: Theme.spacingS
            verticalAlignment: Text.AlignVCenter
        }

        background: Rectangle {
            color: parent.hovered ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.08) : "transparent"
            radius: Theme.cornerRadius / 2
        }

        onTriggered: WifiConnectionActions.editSavedNetworkFromDetails(root.currentSSID, root.currentSecured, root.currentSaved, root.currentEnterprise)
    }

    MenuItem {
        text: root.currentAutoconnect ? I18n.tr("Disable Autoconnect") : I18n.tr("Enable Autoconnect")
        height: root.canToggleAutoconnect ? 32 : 0
        visible: root.canToggleAutoconnect

        contentItem: StyledText {
            text: parent.text
            font.pixelSize: Theme.fontSizeSmall
            color: Theme.surfaceText
            leftPadding: Theme.spacingS
            verticalAlignment: Text.AlignVCenter
        }

        background: Rectangle {
            color: parent.hovered ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.08) : "transparent"
            radius: Theme.cornerRadius / 2
        }

        onTriggered: NetworkService.setWifiAutoconnect(root.currentSSID, !root.currentAutoconnect)
    }

    MenuItem {
        text: I18n.tr("Forget Network")
        height: root.showSavedOptions ? 32 : 0
        visible: root.showSavedOptions

        contentItem: StyledText {
            text: parent.text
            font.pixelSize: Theme.fontSizeSmall
            color: Theme.error
            leftPadding: Theme.spacingS
            verticalAlignment: Text.AlignVCenter
        }

        background: Rectangle {
            color: parent.hovered ? Qt.rgba(Theme.error.r, Theme.error.g, Theme.error.b, 0.08) : "transparent"
            radius: Theme.cornerRadius / 2
        }

        onTriggered: root.forgetRequested(root.currentSSID)
    }
}
