package pepper

import "fmt"

type OfflineAction int

const (
	// OfflineActionDisablePage is the default option. It will disable the page
	// and cover the screen with a semi-transparent sheet containing a
	// "Disconnected" message. Using this mode, not even text will be able to be
	// copied.
	OfflineActionDisablePage = iota

	// OfflineActionDisableForms will disable any form elements on the page.
	// However, the rest of the page will remain viewable/scrollable.
	OfflineActionDisableForms

	// OfflineActionDoNothing is not recommended. It will not give any
	// indication that the server has gone offline. Interacting with the page
	// (such as form elements) will partially to work but buttons will do
	// nothing, and any state of the page will be lost when the server
	// reconnects.
	//
	// This mode may only be useful when showing static data (non-interactive
	// apps) over a flaky connection.
	OfflineActionDoNothing
)

func getIsConnectedJavascript(action OfflineAction) (js string) {
	switch action {
	case OfflineActionDisablePage:
		return `
		var display = isConnected ? "none" : "flex";
		document.getElementById("disconnectedoverlay").style.display = display;
	`

	case OfflineActionDisableForms:
		return `
		document.getElementById("disconnectedoverlay").style.display = "none";
		if (!isConnected) {
			["input", "button", "textarea", "select"].forEach((selector) => {
				var matches = document.querySelectorAll(selector);
				for (var i = 0; i < matches.length; i++) {
					matches[i].disabled = true;
				}
			})
		}
	`

	case OfflineActionDoNothing:
		return ""
	}

	panic(fmt.Sprintf("Invalid OfflineAction: %v", action))
}
