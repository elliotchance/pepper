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
			if payload["method"] != "app.Refresh" {
				name := strings.Split(payload["method"], ".")[1]
				reflect.ValueOf(app).MethodByName(name).Call(nil)
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
