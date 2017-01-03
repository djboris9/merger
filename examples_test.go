package merger

import (
	"encoding/json"
	"fmt"
	"log"
)

func ExampleMerge() {
	// Create struct A
	A := struct {
		FieldA string
		FieldB string
		FieldC []int
	}{
		"aAaA",
		"bBbB",
		[]int{1, 2},
	}

	// Create struct B
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

	// Serialize
	ser, _ := json.Marshal(V)
	fmt.Println(string(ser))
	// Output: {"FieldA":"NewVal","FieldB":"bBbB","FieldC":[1,2,3,4]}
}

func ExampleMerge_nested() {
	// Create two maps with nested types
	A := map[string]interface{}{
		"a": []int{1, 2},
		"b": "aAaAaA",
		"c": map[int]int{1: 10, 2: 20},
	}

	B := map[string]interface{}{
		"a": []int{3, 4},
		"b": "bBbBbB",
		"c": map[int]int{1: 100, 3: 30},
		"d": "Only in B",
	}

	// Merge map A and B together, including their values
	V, err := Merge(A, B)
	if err != nil {
		log.Fatal(err)
	}

	// Serialize
	ser, _ := json.Marshal(V)
	fmt.Println(string(ser))
	// Output: {"a":[1,2,3,4],"b":"bBbBbB","c":{"1":100,"2":20,"3":30},"d":"Only in B"}
}
