# Merger
This Go package merges different types together. Nested types will be merged too.
This can be useful for:

* Merging configurations from different sources
* Deduplicate items before serialization
* Setting field preferences

It works with Go 1.7 and greater

## Usage
Documentation with examples is available on https://godoc.org/github.com/djboris9/merger
This readme is only an introduction.

### Merge algorithm
Let A be the first arbitrary value and B be the second arbitrary value with precendence.
boolean, numeric and string types will be overwritten by the argument with precendence (B).
slice and array types will be concatenated (A ∥ B).
struct and map types will be merged together giving a union of all fields, where the values of them are merged too (A ∪ B)

### Example
Full usage example:
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

Merging maps:
```Go
A := map[string]int{
	"x": 1,
	"y": 2,
}
B := map[string]int{
	"a": 1,
	"b": 2,
}
V, _ := merger.Merge(A, B)
// V: map[string]int{"a": 1, "b": 2, "x": 1, "y": 2}
```

Merging slices:
```Go
A := []int{1, 2, 3}
B := []int{4, 5, 6}
V, _ := merger.Merge(A, B)
// V: []int{1, 2, 3, 4, 5}
```

Other examples are in the godoc.
