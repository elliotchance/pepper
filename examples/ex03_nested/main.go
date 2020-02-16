package main

import "github.com/elliotchance/pepper"

type Counter struct {
	Number int
}

func (c *Counter) Render() (string, error) {
	return `
		Counter: {{ .Number }}
		<button @click="AddOne">+</button>
	`, nil
}

func (c *Counter) AddOne() {
	c.Number++
}

type Counters struct {
	Counters []*Counter
}

func (c *Counters) Render() (string, error) {
	return `
		<table>
			{{ range .Counters }}
				<tr><td>
					{{ render . }}
				</td></tr>
			{{ end }}
			<tr><td>
				Total: {{ call .Total }}
			</td></tr>
		</table>
	`, nil
}

func (c *Counters) Total() int {
	total := 0
	for _, counter := range c.Counters {
		total += counter.Number
	}

	return total
}

func main() {
	panic(pepper.NewServer().Start(func(_ *pepper.Connection) pepper.Component {
		return &Counters{
			Counters: []*Counter{
				{}, {}, {},
			},
		}
	}))
}
