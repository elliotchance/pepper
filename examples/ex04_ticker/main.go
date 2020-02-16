package main

import (
	"github.com/elliotchance/pepper"
	"time"
)

type Clock struct{}

func (c *Clock) Render() (string, error) {
	return `
		The time now is {{ call .Now }}.
	`, nil
}

func (c *Clock) Now() string {
	return time.Now().Format(time.RFC1123)
}

func main() {
	panic(pepper.NewServer().Start(func(conn *pepper.Connection) pepper.Component {
		go func() {
			for range time.NewTicker(time.Second).C {
				conn.Update()
			}
		}()

		return &Clock{}
	}))
}
