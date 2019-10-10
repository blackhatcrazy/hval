package main

import (
	"fmt"
	"hval/internal/aggregate"
	"hval/pkg/flags"
	"hval/pkg/flatmap"
	"log"

	"gopkg.in/yaml.v2"
)

func main() {

	args, err := flags.Parse()
	check(err)
	fmt.Println(args)
	aggMap, err := aggregate.Inputs(args.Files)
	check(err)
	infl, err := flatmap.Inflate(aggMap)
	check(err)
	d, err := yaml.Marshal(&infl)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))
	fmt.Println()

}

func check(err error) {
	// TODO add proper error and log handling
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
