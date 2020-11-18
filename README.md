# flagslice

[![PkgGoDev](https://pkg.go.dev/badge/github.com/icholy/flagslice)](https://pkg.go.dev/github.com/icholy/flagslice)

> This package provides a reflect based solution for reading flags into a slice.

## Supported Types:

* `bool`
* `int`
* `int64`
* `uint`
* `uint64`
* `float64`
* `string`
* `time.Duration`
* anything implementing `flag.Value`

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
