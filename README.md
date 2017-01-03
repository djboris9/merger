# Merger
This Go package merges different types together.
This can be useful for:

* Merging configurations from different sources
* Deduplicate items before serialization
* Setting field preferences

## Usage

### Merge algorithm
Calling `merger.Merge(a, b)` will merge `a` and `b` together, where `b` has precedence.
So if you call `merger.Merge("Hello", "World")` the output will be `"World"`.

`string`, `int`, `int64`, `complex` and so on will be overwritten by the argument with precendence.
`struct`s and `map`s will be merged together (like FULL OUTER JOIN).
`slice`s and `array`s will be concatenated.

### Example
```Go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/djboris9/merger"
)

func main() {
	A := struct {
		FieldA string
		FieldB string
		FieldC []int
	}{
		"aAaA",
		"bBbB",
		[]int{1, 2},
	}

	B := struct {
		FieldA string
		FieldC []int
	}{
		"NewVal",
		[]int{3, 4},
	}

	// Merge struct A and B together
	V, err := merger.Merge(A, B)
	if err != nil {
		log.Fatal(err)
	}

	// Print it
	ser, _ := json.Marshal(V)
	fmt.Println(string(ser))
	// Output: {"FieldA":"NewVal","FieldB":"bBbB","FieldC":[1,2,3,4]}
}
```
