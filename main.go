package pepper

import (
	"html/template"
	"net/http"
)

func StartServer(newConnectionFn NewConnectionFunc) error {
	http.HandleFunc("/ws", newConnection(newConnectionFn))
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
var ws;

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
	}
	ws.onerror = function(evt) {
		console.log("WebSocketError", evt.data);
	}
}

function send(method) {
	var payload = {method: method};
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
