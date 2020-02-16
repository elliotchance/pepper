package pepper

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Server struct {
	// ConnectionRetentionTime is the maximum time allowed for a connection to
	// exist in memory after the client has disconnected. A duration greater
	// than zero ensured that a client can reconnect without losing any state.
	// However, the trade off is the extra memory consumed by holding too many
	// truly gone clients in memory. The default is 1 minute.
	//
	// ReconnectInterval should be less than ConnectionRetentionTime.
	ConnectionRetentionTime time.Duration

	// ConnectionRetentionInterval is the amount of time between each operation
	// to locate and evict old connections that are beyond
	// ConnectionRetentionTime. Lowing this value will increase the amount of
	// CPU, but may also decrease the amount of memory if many connections are
	// short-lived. The default is 1 minute.
	ConnectionRetentionInterval time.Duration

	// HeartbeatInterval is how often the client will notify the server that the
	// connection is still active. It is important that this value be lower than
	// ConnectionRetentionTime. The default is 1 second.
	HeartbeatInterval time.Duration

	// OfflineAction controls the behavior of the client when it loses
	// connection with the server. See constants for explanation.
	OfflineAction OfflineAction

	// ReconnectInterval configures how long the client should wait before
	// trying to reconnect to the server. The default is 1 second.
	//
	// Also see ConnectionRetentionTime.
	ReconnectInterval time.Duration
}

func NewServer() *Server {
	return &Server{
		ConnectionRetentionTime:     time.Minute,
		ConnectionRetentionInterval: time.Minute,
		HeartbeatInterval:           time.Second,
		OfflineAction:               OfflineActionDisablePage,
		ReconnectInterval:           time.Second,
	}
}

// Start will start the application. Each client that connects will call
// newConnectionFn.
func (server *Server) Start(newConnectionFn NewConnectionFunc) error {
	go server.startConnectionRetentionGarbageCollector()

	http.HandleFunc("/ws/", websocketHandler(newConnectionFn))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cid := rand.Int()
		err := homeTemplate.Execute(w, map[string]interface{}{
			"ws":                fmt.Sprintf("ws://%s/ws/%d", r.Host, cid),
			"isConnected":       template.JS(getIsConnectedJavascript(server.OfflineAction)),
			"reconnectInterval": server.ReconnectInterval.Milliseconds(),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}
	})

	return http.ListenAndServe("localhost:8080", nil)
}

func (server *Server) startConnectionRetentionGarbageCollector() {
	for {
		time.Sleep(server.ConnectionRetentionInterval)

		evictConnectionsOlderThan := time.Now().Add(-server.ConnectionRetentionTime)

		for cid, connection := range connections {
			if connection.LastSeen.Before(evictConnectionsOlderThan) {
				log.Println("evicting connection:", cid)
				delete(connections, cid)
			}
		}
	}
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
var ws, activeKey, heartbeat;

function setIsConnected(isConnected) {
	{{ .isConnected }}
}

function openConnection() {
	ws = new WebSocket({{ .ws }});
	ws.onopen = function(evt) {
		console.log("WebSocketOpened");
		send("app.Refresh");
		setIsConnected(true);
		setInterval(function () {
			ws.send(JSON.stringify({method: 'heartbeat'}));
		}, {{ .heartbeatInterval }});
	}
	ws.onclose = function(evt) {
		console.log("WebSocketClosed");
		ws = null;
		setIsConnected(false);
		setTimeout(openConnection, {{ .reconnectInterval }});
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

function sendSetAttribute(component, name, value) {
	var payload = {
		method: component + ".SetAttribute",
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
<style>
.disconnected {
    display: flex;
    position: fixed; 
    top: 0; bottom: 0; left: 0; right: 0;
    width: 100%;
    height: 100%;
    background-color: black;
    opacity: 0.5;
	color: white;
	font-size: 32px;
	font-weight: bold;
	align-items: center;
}
</style>
</head>
<body>
<div id="app"></div>
<div class="disconnected" id="disconnectedoverlay"><div style="text-align: center; width: 100%">Disconnected</div></div>
</body>
</html>
`))
