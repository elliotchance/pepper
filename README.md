# 🌶️ pepper

Create reactive frontends without ever writing frontend code.

   * [How Does It Work?](#how-does-it-work)
   * [What Should/Shouldn't I Use It For?](#what-shouldshouldnt-i-use-it-for)
   * [Examples](#examples)
      * [#1: A Simple Counter](#1-a-simple-counter)
      * [#2: Forms](#2-forms)
      * [#3: Nested Components](#3-nested-components)

# How Does It Work?

pepper runs a HTTP server that returns a tiny empty HTML page with just a few
lines of inline javascript. Immediately after the initial page loads it will
connect through a websocket.

All event triggered on the browser will be sent through the websocket where
state changes and rerendering occurs. The result is passed back to the websocket
to update the UI.

At the moment it returns the whole rendered component. However, this could be
optimized in the future to only return the differences.

# What Should/Shouldn't I Use It For?

pepper requires a constant connection to the server (for the websocket) so it
wouldn't work for anything that must function offline, or used in cases where
the internet is flaky.

I imagine some good use cases for pepper would be:

1. Showing real time data. Streaming metrics, graphs, logs, dashboards, etc.
2. Apps that rely on a persistent connection. Such as chat clients, timed
interactive exams, etc.
3. Apps that would benefit from persistent state. The entire state can be saved
or restored into a serialized format like JSON. Great for forms or surveys with
many questions/steps.
4. Prototyping a frontend app. It's super easy to get up and running and iterate
changes without setting up a complex environment, build tools and dependencies.

# Examples

## #1: A Simple Counter

```go
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

func main() {
	panic(pepper.StartServer(func() pepper.Component {
		return &Counter{}
	}))
}
```

- The `Render` method returns a `html/template` syntax, or an error.
- `@click` will trigger `AddOne` to be called when the button is clicked.
- Any event triggered from the browser will cause the component to rerender
automatically.

Try it now:

```bash
go get -u github.com/elliotchance/pepper/examples/ex01_counter
ex01_counter
```

Then open: [http://localhost:8080/](http://localhost:8080/)

## #2: Forms

```go
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
```

- Any `html/template` syntax will work, including loops with `{{ range }}`.
- `@value` will cause the `Name` property to be bound with the text box in both
directions.
- Since there are multiple "Delete" buttons (one for each person), you should
specify a `key`. The `key` is passed as the first argument to the `Delete`
function.

Try it now:

```bash
go get -u github.com/elliotchance/pepper/examples/ex02_form
ex02_form
```

Then open: [http://localhost:8080/](http://localhost:8080/)


## #3: Nested Components

```go
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
	panic(pepper.StartServer(func() pepper.Component {
		return &Counters{
			Counters: []*Counter{
				{}, {}, {},
			},
		}
	}))
}
```

- This example uses three `Counter` components (from Example #1) and includes a
live total.
- Components can be nested with the `render` function. The nested components do
not need to be modified in any way.
- Invoke methods with the `call` function.

Try it now:

```bash
go get -u github.com/elliotchance/pepper/examples/ex03_nested
ex03_nested
```

Then open: [http://localhost:8080/](http://localhost:8080/)
