package render

import (
	"bytes"
	"fmt"
	"hval/pkg/flatmap"
	"io"
	"text/template"
)

type tmpl struct {
	err          error
	dataCurrent  []byte
	dataPrevious []byte
	input        interface{}
	iteration    int
}

func NewTemplate(
	template []byte,
	input interface{},
	flatInput map[string]flatmap.MapEntry,
) (tmpl, error) {

	if cyclicReferences(flatInput) {
		return tmpl{}, fmt.Errorf("cyclic dependency in aggregated values")
	}

	return tmpl{
		err:          nil,
		dataCurrent:  template,
		dataPrevious: []byte{},
		input:        input,
		iteration:    0,
	}, nil
}

func cyclicReferences(input map[string]flatmap.MapEntry) bool {
	return false
}

func (t *tmpl) Render(w io.Writer) error {
	// TODO: check for closed loops here?
	for t.hasChanged() {
		t.render()
		if t.err != nil {
			return t.err
		}
	}
	_, err := w.Write(t.dataCurrent)
	return err
}

func (t *tmpl) hasChanged() bool {
	if n := bytes.Compare(t.dataPrevious, t.dataCurrent); n == 0 {
		return false
	}
	return true
}

func (t *tmpl) render() {
	if t.err != nil {
		return
	}
	t.dataPrevious = make([]byte, len(t.dataCurrent))
	copy(t.dataPrevious, t.dataCurrent)

	tt, err := template.New(
		fmt.Sprintf("tmp_%d", t.iteration),
	).Parse(string(t.dataPrevious))
	if err != nil {
		t.err = err
		return
	}
	res := new(bytes.Buffer)

	if err := tt.Execute(res, t.input); err != nil {
		t.err = err
		return
	}
	t.iteration++
	t.dataCurrent = res.Bytes()
	return
}
