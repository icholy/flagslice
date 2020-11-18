package flagslice_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/icholy/flagslice"
)

func ExampleVar(t *testing.T) {
	var names []string
	flagslice.Var(&names, "n", "name")
	flag.Parse()
	for _, name := range names {
		fmt.Println(name)
	}
}
