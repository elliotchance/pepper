package main

import "github.com/elliotchance/pepper"

type Counter struct {
	Number int
}

func (c *Counter) Render() (string, error) {
	return `
		Counter: {{ .Number }}<br/>
		<button @click="AddOne">+</button>
	`, nil
}

func (c *Counter) AddOne() {
	c.Number++
}

func main() {
	panic(pepper.StartServer(func() pepper.Component {
		return &Counter{}
	}))
}
