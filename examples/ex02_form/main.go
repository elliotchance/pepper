package main

import (
	"github.com/elliotchance/pepper"
	"strconv"
)

type People struct {
	Names []string
	Name  string
}

func (c *People) Render() (string, error) {
	return `
		<table>
			{{ range $i, $name := .Names }}
				<tr><td>
					{{ $name }}
					<button key="{{ $i }}" @click="Delete">Delete</button>
				</td></tr>
			{{ end }}
		</table>
		Add name: <input type="text" @value="Name">
		<button @click="Add">Add</button>
	`, nil
}

func (c *People) Delete(key string) {
	index, _ := strconv.Atoi(key)
	c.Names = append(c.Names[:index], c.Names[index+1:]...)
}

func (c *People) Add() {
	c.Names = append(c.Names, c.Name)
	c.Name = ""
}

func main() {
	panic(pepper.StartServer(func() pepper.Component {
		return &People{
			Names: []string{"Jack", "Jill"},
		}
	}))
}
