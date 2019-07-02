package main

import (
	"github.com/elliotchance/pepper"
	"html/template"
)

type Counter struct {
	Number int
}

func (c *Counter) Render() (*template.Template, error) {
	return template.New("").Parse(`
		Counter: {{ .Number }}<br/>
		<button onclick="{{ .AddOne }}">+</button>`)
}

func (c *Counter) AddOne() {
	c.Number++
}

func main() {
	panic(pepper.StartServer(func() pepper.Component {
		return &Counter{}
	}))
}
