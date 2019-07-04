package pepper

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type NewConnectionFunc func() Component

var upgrader = websocket.Upgrader{} // use default options

func newConnection(newConnection NewConnectionFunc) func(w http.ResponseWriter, r *http.Request) {
	app := newConnection()

	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)

			var payload map[string]string
			err = json.Unmarshal(message, &payload)
			if err != nil {
				log.Println("read:", err)
				break
			}

			var response string
			parts := strings.Split(payload["method"], ".")

			if parts[0] == "app" && parts[1] == "Refresh" {
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

			response, err = render(app)
			if err != nil {
				log.Println("write:", err)
				break
			}

			err = c.WriteMessage(mt, []byte(response))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}
