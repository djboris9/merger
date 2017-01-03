package merger

import (
	"encoding/json"
	"fmt"
	"log"
)

func ExampleMerge() {
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
	V, err := Merge(A, B)
	if err != nil {
		log.Fatal(err)
	}

	// Print it
	ser, _ := json.Marshal(V)
	fmt.Println(string(ser))
	// Output: {"FieldA":"NewVal","FieldB":"bBbB","FieldC":[1,2,3,4]}
}
