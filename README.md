# flagslice

> This package provides a reflect based solution for reads flags into a slice.

## Supported Types:

* `int`
* `int64`
* `uint`
* `uint64`
* `float64`
* `string`
* `time.Duration`
* `flag.Value`

## Example

``` go
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
```
