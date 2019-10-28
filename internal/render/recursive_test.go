package render

import (
	"fmt"
	"log"
	"testing"

	"gopkg.in/yaml.v2"
)

var testTmpl = `
k2: world
t1: {{ printf "hello-%s" .k2 }}
t2: {{ .t1 }}{{ .t1 }}
t3: {{ .t2 }} bla
v4:
  v41: test
t4: {{ .v4 }}
t5:
	t51: {{ .t4 }}
t6: {{ .t5.t51 }}
normal: value
`

type myMap map[string]interface{}

func (m myMap) String() string {
	res := "{"
	for k, v := range m {
		res = fmt.Sprintf("%s%s: %s, ", res, k, v)
	}
	return fmt.Sprintf("%s}", res[:len(res)-2])
}

var testInput = map[string]interface{}{
	"k2":     "world",
	"t1":     `{{ printf "hello-%s" .k2 }}`,
	"t2":     `{{ .t1 }}{{ .t1 }}`,
	"t3":     `{{ .t2 }} bla`,
	"v4":     myMap{"v41": "test"},
	"t4":     `{{ .v4 }}`,
	"t5":     myMap{"t51": `{{ .t4 }}`},
	"normal": "value",
}

func TestRender(t *testing.T) {
	templ := tmpl{}
	templ.dataCurrent = []byte(testTmpl)
	newMap := myMap(testInput)
	templ.input = newMap
	fmt.Println(newMap)
	for templ.hasChanged() {
		templ.render()
		fmt.Println("it ", templ.iteration)
		fmt.Println(templ.hasChanged())
		fmt.Printf("err %v\n", templ.err)
		fmt.Println("data", string(templ.dataCurrent))
	}
	var tr map[string]interface{}

	err := yaml.Unmarshal(templ.dataCurrent, &tr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", tr)

	t.FailNow()

}
