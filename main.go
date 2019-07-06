package pepper

import (
	"html/template"
	"net/http"
)

func StartServer(newConnectionFn NewConnectionFunc) error {
	http.HandleFunc("/ws", websocketHandler(newConnectionFn))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
	})

	return http.ListenAndServe("localhost:8080", nil)
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
var ws, activeKey;

function openConnection() {
	ws = new WebSocket("{{.}}");
	ws.onopen = function(evt) {
		console.log("WebSocketOpened");
		send("app.Refresh");
	}
	ws.onclose = function(evt) {
		console.log("WebSocketClosed");
		ws = null;
	}
	ws.onmessage = function(evt) {
		console.log("WebSocketReceived", evt.data);
		document.getElementById("app").innerHTML = evt.data;
		if (activeKey) {
			const el = document.querySelectorAll('[activekey=' + activeKey + ']')[0];
			el.focus();
			el.setSelectionRange(20, 20);
		}
	}
	ws.onerror = function(evt) {
		console.log("WebSocketError", evt.data);
	}
}

function send(method, self) {
	var payload = {
		method: method,
		key: self ? (self.attributes.key ? self.attributes.key.value : undefined) : undefined,
		value: self ? self.value : undefined,
	};
	console.log("WebSocketSending", payload);
	ws.send(JSON.stringify(payload));
}

function sendSetAttribute(name, value) {
	var payload = {
		method: "app.SetAttribute",
		key: name,
		value: value,
	};
	activeKey = name;
	console.log("WebSocketSending", payload);
	ws.send(JSON.stringify(payload));
}

window.addEventListener("load", function(evt) {
	openConnection();
});
</script>
</head>
<body>
<div id="app"></div>
</body>
</html>
`))
