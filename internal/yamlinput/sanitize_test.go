package yamlinput_test

import (
	"bytes"
	"fmt"
	"hval/internal/yamlinput"
	"strings"
	"testing"
)

var testValues = []struct {
	deSanitized string
	sanitized   string
	errSan      error
	errDesan    error
}{
	{
		deSanitized: `t1: {{ printf "%s" .k2 }}
t2: {{ .t1 }}{{ .t1 }}
t3: {{ .t1 }}{{ .t1 }} additional text
normal: value
another: normal value`,
		sanitized: `t1: ` + "`" + `{{ printf "%s" .k2 }}` + "`" + `
t2: ` + "`" + `{{ .t1 }}{{ .t1 }}` + "`" + `
t3: ` + "`" + `{{ .t1 }}{{ .t1 }} additional text` + "`" + `
normal: value
another: normal value`,
		errSan:   nil,
		errDesan: nil,
	},
	{
		deSanitized: `{{ .illegal }}: {{ printf "%s" .k2 }}
t2: {{ .t1 }}{{ .t1 }}
t3: {{ .t1 }}{{ .t1 }} additional text
normal: value`,
		sanitized: `{{ .illegal }}: ` + "`" + `{{ printf "%s" .k2 }}` + "`" + `
t2: ` + "`" + `{{ .t1 }}{{ .t1 }}` + "`" + `
t3: ` + "`" + `{{ .t1 }}{{ .t1 }} additional text` + "`" + `
normal: value`,
		errSan: fmt.Errorf("illegal key found in line \"%s\"",
			`{{ .illegal }}: {{ printf "%s" .k2 }}`,
		),
		errDesan: nil,
	},
}

func TestSanitize(t *testing.T) {
	for _, test := range testValues {
		debugSan := new(bytes.Buffer)
		v := yamlinput.New(debugSan, true)
		debugDesan := new(bytes.Buffer)
		vOut := yamlinput.New(debugDesan, true)
		res, err := v.Sanitize([]byte(test.deSanitized))
		if err != test.errSan {
			if err.Error() != test.errSan.Error() {
				t.Errorf("expected error does not match.\nResult: %+v \nExpect: %+v",
					err, test.errSan)
			}
		} else {
			fmt.Println(res)
			fmt.Println([]byte(test.sanitized))
			if n := strings.Compare(string(res), test.sanitized); n != 0 {
				t.Errorf("sanitization result does not match.\nResult: \"%+v\" \nExpect: \"%+v\"",
					string(res), test.sanitized)
			}
			if n := strings.Compare(string(res), debugSan.String()); n != 0 {
				t.Errorf("sanitization debugin output does not match result.\nResult: \"%+v\" \nExpect: \"%+v\"",
					debugSan.String(), string(res))
			}
		}

		resDesan, err := vOut.Desanitize([]byte(test.sanitized))
		if err != test.errDesan {
			if err.Error() != test.errDesan.Error() {
				t.Errorf("expected error does not match.\nResult: %+v \nExpect: %+v",
					err, test.errSan)
			}
		} else {
			if n := strings.Compare(string(resDesan), test.deSanitized); n != 0 {
				t.Errorf("unsanitization result does not match.\nResult: \"%+v\" \nExpect: \"%+v\"",
					string(resDesan), test.deSanitized)
			}
			if n := strings.Compare(string(resDesan), debugDesan.String()); n != 0 {
				t.Errorf("sanitization debugin output does not match result.\nResult: \"%+v\" \nExpect: \"%+v\"",
					debugDesan.String(), string(resDesan))
			}
		}
	}
}
