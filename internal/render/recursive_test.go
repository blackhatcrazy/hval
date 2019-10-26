package render

import (
	"fmt"
	"testing"
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
normal: value
`

var testInput = map[string]interface{}{
	"k2":     "world",
	"t1":     `{{ printf "hello-%s" .k2 }}`,
	"t2":     `{{ .t1 }}{{ .t1 }}`,
	"t3":     `{{ .t2 }} bla`,
	"v4":     map[string]interface{}{"v41": "test"},
	"t4":     `{{ .v4 }}`,
	"t5":     map[string]interface{}{"t51": `{{ .t4 }}`},
	"normal": "value",
}

func TestRender(t *testing.T) {
	templ := tmpl{}
	templ.dataCurrent = []byte(testTmpl)
	templ.input = testInput
	for templ.hasChanged() {
		templ.render()
		fmt.Println("it ", templ.iteration)
		fmt.Println(templ.hasChanged())
		fmt.Printf("err %v\n", templ.err)
		fmt.Println("data", string(templ.dataCurrent))
	}

	t.FailNow()

}
