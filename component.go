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

func render(c Component) ([]byte, error) {
	templateData, err := c.Render()
	if err != nil {
		return nil, err
	}

	templateData = "{{ $template := . }}" + replaceAllStringSubmatchFunc(
		regexp.MustCompile(`@(\w+)="(.*?)"`),
		templateData,
		func(strings []string) string {
			event, attribute := strings[1], strings[2]
			if event == "value" {
				return fmt.Sprintf(
					`activekey="%s" value="{{ $template.%s }}" onkeyup="sendSetAttribute('%s', this.value)"`,
					attribute, attribute, attribute)
			}

			return fmt.Sprintf(`on%s="send('app.%s', this)"`, event, attribute)
		})

	t, err := template.New("").Parse(templateData)
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
		data[name] = template.JS(fmt.Sprintf(`send('app.%s', this)`, name))
	}

	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
