package pepper

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Connection struct {
	LastSeen      time.Time
	RootComponent Component
	Connection    *websocket.Conn
}

type NewConnectionFunc func(*Connection) Component

var upgrader = websocket.Upgrader{} // use default options

func websocketHandler(newConnectionFunc NewConnectionFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()

		connection := newConnection(c, r.URL.Path, newConnectionFunc)
		err = connection.start()
		if err != nil {
			log.Print("start:", err)
			return
		}
	}
}

var connections = map[string]*Connection{}

func newConnection(c *websocket.Conn, cid string, newConnection NewConnectionFunc) *Connection {
	if connection, isExisting := connections[cid]; isExisting {
		log.Println("client reconnected:", cid)
		connection.Connection = c
		connection.LastSeen = time.Now()

		return connection
	}

	log.Println("new client:", cid)

	connection := &Connection{
		Connection: c,
		LastSeen:   time.Now(),
	}

	connection.RootComponent = newConnection(connection)
	connections[cid] = connection

	return connection
}

func (conn *Connection) start() error {
	for {
		_, message, err := conn.Connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		var payload map[string]string
		err = json.Unmarshal(message, &payload)
		if err != nil {
			log.Println("unmarshal:", err)
			break
		}

		var response string
		parts := strings.Split(payload["method"], ".")

		// Heartbeat
		if parts[0] == "heartbeat" {
			conn.LastSeen = time.Now()
		} else if parts[0] == "app" && parts[1] == "Refresh" {
			// Do nothing, fall through to rerender.
		} else if parts[1] == "SetAttribute" {
			component := getComponentByID(parts[0])
			reflect.ValueOf(component).
				Elem().
				FieldByName(payload["key"]).
				SetString(payload["value"])
		} else {
			component := getComponentByID(parts[0])
			method := reflect.ValueOf(component).MethodByName(parts[1])

			var params []reflect.Value

			if method.Type().NumIn() > 0 {
				params = append(params, reflect.ValueOf(payload["key"]))
			}

			if method.Type().NumIn() > 1 {
				params = append(params, reflect.ValueOf(payload["value"]))
			}

			method.Call(params)
		}

		response, err = Render(conn.RootComponent)
		if err != nil {
			log.Println("render:", err)
			break
		}

		err = conn.Connection.WriteMessage(websocket.TextMessage, []byte(response))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

	return nil
}

func (conn *Connection) Update() {
	response, err := Render(conn.RootComponent)
	if err != nil {
		log.Println("render:", err)
		return
	}

	err = conn.Connection.WriteMessage(websocket.TextMessage, []byte(response))
	if err != nil {
		log.Println("write:", err)
		return
	}
}
