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

			var response []byte
			switch payload["method"] {
			case "app.Refresh":
			// Do nothing, fall through to rerender.

			case "app.SetAttribute":
				reflect.ValueOf(app).
					Elem().
					FieldByName(payload["key"]).
					SetString(payload["value"])

			default:
				name := strings.Split(payload["method"], ".")[1]
				method := reflect.ValueOf(app).MethodByName(name)

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

			err = c.WriteMessage(mt, response)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}
