package main

import (
	"flag"
	"fmt"

	"github.com/icholy/flagslice"
)

func main() {
	var names []string
	flag.Var(flagslice.Value(&names), "name", "a name")
	flag.Parse()
	for _, name := range names {
		fmt.Println(name)
	}
}
