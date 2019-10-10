package flatmap_test

import (
	"fmt"
	"hval/pkg/flatmap"
	"reflect"
	"testing"
)

var testInflate = []struct {
	in  map[string]flatmap.MapEntry
	res map[string]interface{}
	err error
}{
	{
		in: map[string]flatmap.MapEntry{
			"fail": flatmap.MapEntry{
				OrderedKey: []string{},
				Value:      "val",
			}},
		res: map[string]interface{}{},
		err: fmt.Errorf("no key to insert provided for value %v", "val"),
	},
	{
		in: map[string]flatmap.MapEntry{
			"a.b.c": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "c"},
				Value:      "val",
			}},
		res: map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "val",
				}}},
		err: nil,
	},
	{
		in: map[string]flatmap.MapEntry{
			"a.b.c": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "c"},
				Value:      "val",
			},
			"a.b.d": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "d"},
				Value:      []string{"res1", "res2"},
			}},
		res: map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "val",
					"d": []string{"res1", "res2"},
				}}},
		err: nil,
	},
	{
		in: map[string]flatmap.MapEntry{
			"a.b.c": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "c"},
				Value:      "val",
			},
			"a.b.d": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "d"},
				Value:      []string{"res1", "res2"},
			},
			"x.b.d": flatmap.MapEntry{
				OrderedKey: []string{"x", "b", "d"},
				Value:      map[string]interface{}{"hello": "map"},
			},
		},
		res: map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "val",
					"d": []string{"res1", "res2"},
				}},
			"x": map[string]interface{}{
				"b": map[string]interface{}{
					"d": map[string]interface{}{"hello": "map"},
				}},
		},
		err: nil,
	},
}

func TestInflate(t *testing.T) {

	for _, test := range testInflate {
		res, err := flatmap.Inflate(test.in)
		if err != test.err {
			if err.Error() != test.err.Error() {
				t.Errorf("expected error does not match.\nResult: %+v \nExpect: %+v", err, test.err)
			}
		}
		if !reflect.DeepEqual(res, test.res) {
			t.Errorf("maps are not equal.\nResult: %+v \nExpect: %+v", res, test.res)
		}

	}
}
