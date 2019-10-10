package flatmap

import "fmt"

// Inflate takes a flattened map (result of topmap.Flatten function) and returns
// the map to the original nested form (inflates it)
func Inflate(m map[string]MapEntry) (result map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("topmap: %v", r)
				return
			}
			err = fmt.Errorf("%v [recovered]", err)
		}
	}()

	result = make(map[string]interface{})
	for _, v := range m {
		if len(v.OrderedKey) == 0 {
			result = map[string]interface{}{}
			err = fmt.Errorf("no key to insert provided for value %v", v.Value)
			return
		}
		result = upsert(result, v.OrderedKey, v.Value)
	}
	return
}

func upsert(m interface{}, orderedKeys []string, value interface{}) map[string]interface{} {
	switch cast := m.(type) {
	case map[string]interface{}:
		if subMap, ok := cast[orderedKeys[0]]; ok {
			subres := upsert(subMap, orderedKeys[1:], value)
			cast[orderedKeys[0]] = subres
			return cast
		}
		cast[orderedKeys[0]] = overwrite(orderedKeys[1:], value)
		return cast
	default:
		panic(fmt.Errorf("you should not get here - https://xkcd.com/2200/"))
	}
}

func overwrite(orderedKeys []string, value interface{}) interface{} {
	if len(orderedKeys) == 0 {
		return value
	}
	for i := len(orderedKeys) - 1; i > 0; i-- {
		value = map[string]interface{}{orderedKeys[i]: value}
	}
	return map[string]interface{}{orderedKeys[0]: value}
}
