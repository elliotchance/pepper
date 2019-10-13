package pepper

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"regexp"
)

type Component interface {
	Render() (string, error)
}

func replaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0
	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		var groups []string
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}
		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}
	return result + str[lastIndex:]
}

var componentIDs = map[Component]string{}

func getComponentByID(id string) Component {
	for c, cID := range componentIDs {
		if cID == id {
			return c
		}
	}

	return nil
}

func getComponentID(c Component) string {
	if _, ok := componentIDs[c]; !ok {
		componentIDs[c] = fmt.Sprintf("%p", c)[2:]
	}

	return componentIDs[c]
}

func Render(c Component) (string, error) {
	templateData, err := c.Render()
	if err != nil {
		return "", err
	}

	templateData = "{{ $template := . }}" + replaceAllStringSubmatchFunc(
		regexp.MustCompile(`@(\w+)="(.*?)"`),
		templateData,
		func(strings []string) string {
			event, attribute := strings[1], strings[2]
			if event == "value" {
				return fmt.Sprintf(
					`activekey="%s" value="{{ $template.%s }}" onkeyup="sendSetAttribute('%s', '%s', this.value)"`,
					attribute, attribute, getComponentID(c), attribute)
			}

			return fmt.Sprintf(`on%s="send('%s.%s', this)"`,
				event, getComponentID(c), attribute)
		})

	t, err := template.New("").Funcs(map[string]interface{}{
		"render": func(c Component) (template.HTML, error) {
			data, err := Render(c)

			return template.HTML(data), err
		},
	}).Parse(templateData)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	data := map[string]interface{}{}
	ty := reflect.TypeOf(c)
	tyElem := ty.Elem()
	val := reflect.ValueOf(c)
	valElem := val.Elem()

	// Fields
	for i := 0; i < tyElem.NumField(); i++ {
		data[tyElem.Field(i).Name] = valElem.Field(i).Interface()
	}

	// Methods
	for i := 0; i < ty.NumMethod(); i++ {
		name := ty.Method(i).Name
		data[name] = val.Method(i).Interface()
	}

	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}

	div := fmt.Sprintf(`<div id="component-%s">%s</div>`,
		getComponentID(c), buf.String())

	return div, nil
}
