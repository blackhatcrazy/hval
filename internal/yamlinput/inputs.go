package yamlinput

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func (i *input) Load(filePaths []string) (
	map[string]map[string]interface{},
	error,
) {
	result := map[string]map[string]interface{}{}
	for _, path := range filePaths {
		f, err := ioutil.ReadFile(path)
		if err != nil {
			return map[string]map[string]interface{}{}, err
		}
		fSan, err := i.Sanitize(f)
		if err != nil {
			return map[string]map[string]interface{}{}, err
		}
		fmt.Println(string(fSan))
		y := map[string]interface{}{}
		if err := yaml.Unmarshal(fSan, y); err != nil {
			return map[string]map[string]interface{}{}, err
		}
		fmt.Println(y)
		result[path] = y
		fmt.Println(result)
	}
	return result, nil
}
