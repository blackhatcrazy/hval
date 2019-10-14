package main

import (
	"bytes"
	"fmt"
	"hval/internal/render"
	"hval/internal/yamlinput"
	"hval/pkg/flags"
	"hval/pkg/flatmap"
	"log"

	"gopkg.in/yaml.v2"
)

func aggregate(files map[string]map[string]interface{}) (
	map[string]flatmap.MapEntry,
	error,
) {
	aggregatedValues := map[string]flatmap.MapEntry{}
	for _, content := range files {
		fm, err := flatmap.Flatten(content)
		if err != nil {
			return map[string]flatmap.MapEntry{}, err
		}
		for k, v := range fm {
			aggregatedValues[k] = v
		}
	}
	return aggregatedValues, nil
}
func main() {

	args, err := flags.Parse()
	check(err)
	fmt.Println(args)
	debug := new(bytes.Buffer)
	input := yamlinput.New(debug, false)
	files, err := input.Load(args.Files)
	aggMap, err := aggregate(files)
	check(err)
	fmt.Println(files)
	fmt.Println(aggMap)
	infl, err := flatmap.Inflate(aggMap)
	check(err)
	inflBytes, err := yaml.Marshal(&infl)
	tmpl, err := input.Desanitize(inflBytes)
	r, err := render.NewTemplate(tmpl, infl, aggMap)
	check(err)
	output := new(bytes.Buffer)
	check(r.Render(output))
	fmt.Println(output.String())

}

func check(err error) {
	// TODO add proper error and log handling
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
