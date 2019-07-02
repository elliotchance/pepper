package pepper

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
)

type Component interface {
	Render() (*template.Template, error)
}

func render(c Component) ([]byte, error) {
	t, err := c.Render()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	data := map[string]interface{}{}
	ty := reflect.TypeOf(c)
	tyElem := ty.Elem()
	val := reflect.ValueOf(c).Elem()

	// Fields
	for i := 0; i < tyElem.NumField(); i++ {
		data[tyElem.Field(i).Name] = val.Field(i).Interface()
	}

	// Methods
	for i := 0; i < ty.NumMethod(); i++ {
		name := ty.Method(i).Name
		data[name] = template.JS(fmt.Sprintf(`send('app.%s')`, name))
	}

	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
