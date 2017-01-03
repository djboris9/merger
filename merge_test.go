package merger

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

// TODO: Make negative test cases and check for errors

type Case struct {
	A   interface{}
	B   interface{}
	Exp interface{}
}

func TestMain(m *testing.M) {
	// Read tabletest content
	file, err := os.Open("testcases.list")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var counter int
	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "#") {
			continue
		}

		var data interface{}
		err := json.Unmarshal([]byte(t), &data)
		if err != nil {
			log.Fatal(err)
		}

		switch counter % 3 {
		case 0:
			testcases = append(testcases, Case{})
			testcases[len(testcases)-1].A = data
		case 1:
			testcases[len(testcases)-1].B = data
		case 2:
			testcases[len(testcases)-1].Exp = data
		}
		counter++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

/*
 * Table tests
 */
var testcases []Case

func TestTabletest(t *testing.T) {
	for _, c := range testcases {
		res, err := Merge(c.A, c.B)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(res, c.Exp) {
			t.Errorf("Merge(%v, %v) => %v, want %v", c.A, c.B, res, c.Exp)
		} else {
			t.Logf("Merge(%v, %v) => %v", c.A, c.B, res)
		}
	}
}

/*
 * Merge tests
 */
func TestArray(t *testing.T) {
	A := [2]int{1, 2}
	B := [3]int{3, 4, 5}
	Exp := [5]int{1, 2, 3, 4, 5}
	res, err := Merge(A, B)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, Exp) {
		t.Errorf("Merge(%v, %v) => %v, want %v", A, B, res, Exp)
	} else {
		t.Logf("Merge(%v, %v) => %v", A, B, res)
	}
}

func TestStructReplace2(t *testing.T) {
	A := struct {
		A int
		B int
		C int
	}{1, 2, 3}
	B := struct {
		A int
		C int
	}{2, 4}
	Exp := struct {
		A int
		B int
		C int
	}{2, 2, 4}
	res, err := Merge(A, B)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, Exp) {
		t.Errorf("Merge(%v, %v) => %v, want %v", A, B, res, Exp)
	} else {
		t.Logf("Merge(%v, %v) => %v", A, B, res)
	}
}

func TestStructReplace(t *testing.T) {
	A := struct {
		A int
		B int
	}{1, 2}
	B := struct {
		A int
		B int
		C int
	}{1, 4, 5}
	Exp := struct {
		A int
		B int
		C int
	}{1, 4, 5}
	res, err := Merge(A, B)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(res, Exp) {
		t.Errorf("Merge(%v, %v) => %v, want %v", A, B, res, Exp)
	} else {
		t.Logf("Merge(%v, %v) => %v", A, B, res)
	}
}

// Negative tests
func TestDiffSliceTypes(t *testing.T) {
	A := []int{1, 2}
	B := []string{"a", "b"}
	res, err := Merge(A, B)
	t.Logf("Res: %v, err: %s", res, err)

	if err == nil {
		t.Fatalf("Expected an error, but got %v", res)
	}

	if e2, ok := err.(*MergeError); !ok {
		t.Error("Expected an MergeError")
	} else if e2.Type() != ErrDiffSliceTypes {
		t.Errorf("Expected an ErrDiffSliceTypes but got %v (%v)", e2.Type(), e2)
	}
}

func TestDiffArrayTypes(t *testing.T) {
	A := [2]int{1, 2}
	B := [2]string{"a", "b"}
	res, err := Merge(A, B)
	t.Logf("Res: %v, err: %s", res, err)

	if err == nil {
		t.Fatalf("Expected an error, but got %v", res)
	}

	if e2, ok := err.(*MergeError); !ok {
		t.Error("Expected an MergeError")
	} else if e2.Type() != ErrDiffArrayTypes {
		t.Errorf("Expected an ErrDiffArrayTypes but got %v (%v)", e2.Type(), e2)
	}
}

func TestDiffMapKeyTypes(t *testing.T) {
	A := map[int]int{1: 1, 2: 2}
	B := map[string]int{"1": 1, "2": 2}
	res, err := Merge(A, B)
	t.Logf("Res: %v, err: %s", res, err)

	if err == nil {
		t.Fatalf("Expected an error, but got %v", res)
	}

	if e2, ok := err.(*MergeError); !ok {
		t.Error("Expected an MergeError")
	} else if e2.Type() != ErrDiffMapKeyTypes {
		t.Errorf("Expected an ErrDiffMapKeyTypes but got %v (%v)", e2.Type(), e2)
	}
}

func TestDiffMapValueTypes(t *testing.T) {
	A := map[int]int{1: 1, 2: 2}
	B := map[int]string{1: "1", 2: "2"}
	res, err := Merge(A, B)
	t.Logf("Res: %v, err: %s", res, err)

	if err == nil {
		t.Fatalf("Expected an error, but got %v", res)
	}

	if e2, ok := err.(*MergeError); !ok {
		t.Error("Expected an MergeError")
	} else if e2.Type() != ErrDiffMapValueTypes {
		t.Errorf("Expected an ErrDiffMapValueTypes but got %v (%v)", e2.Type(), e2)
	}
}

func TestSliceArrayMerge(t *testing.T) {
	A := []int{1, 2}
	B := [2]int{3, 4}
	res, err := Merge(A, B)
	t.Logf("Res: %v, err: %s", res, err)

	if err == nil {
		t.Fatalf("Expected an error, but got %v", res)
	}

	if e2, ok := err.(*MergeError); !ok {
		t.Error("Expected an MergeError")
	} else if e2.Type() != ErrDiffKind {
		t.Errorf("Expected an ErrDiffKind but got %v (%v)", e2.Type(), e2)
	}
}

func TestPointerMerge(t *testing.T) {
	A := "a"
	B := &A
	res, err := Merge(A, B)
	t.Logf("Res: %v, err: %s", res, err)

	if err == nil {
		t.Fatalf("Expected an error, but got %v", res)
	}

	if e2, ok := err.(*MergeError); !ok {
		t.Error("Expected an MergeError")
	} else if e2.Type() != ErrMergeUnsupported {
		t.Errorf("Expected an ErrMergeUnsupported but got %v (%v)", e2.Type(), e2)
	}
}

func TestDeepMapError(t *testing.T) {
	A := map[int]interface{}{1: "a"}
	B := map[int]interface{}{1: []int{}}
	res, err := Merge(A, B)
	t.Logf("Res: %v, err: %s", res, err)

	if err == nil {
		t.Fatalf("Expected an error, but got %v", res)
	}

	if e2, ok := err.(*MergeError); !ok {
		t.Error("Expected an MergeError")
	} else if e2.Type() != ErrMergeUnsupported {
		t.Errorf("Expected an ErrMergeUnsupported but got %v (%v)", e2.Type(), e2)
	}
}

func TestDeepStructError(t *testing.T) {
	A := struct {
		F string
	}{"a"}
	B := struct {
		F []int
	}{[]int{2}}
	res, err := Merge(A, B)
	t.Logf("Res: %v, err: %s", res, err)

	if err == nil {
		t.Fatalf("Expected an error, but got %v", res)
	}

	if e2, ok := err.(*MergeError); !ok {
		t.Error("Expected an MergeError")
	} else if e2.Type() != ErrMergeUnsupported {
		t.Errorf("Expected an ErrMergeUnsupported but got %v (%v)", e2.Type(), e2)
	}
}

func TestStructInvalidFields(t *testing.T) {
	// Expect ErrInvalidFields
	//
	// sf := []reflect.StructField{reflect.StructField{
	// 	Name: "F",
	// 	Type: reflect.TypeOf(nil),
	// }}
	// sT := reflect.StructOf(sf)
	// A := reflect.New(sT).Interface()
	// B := reflect.New(sT).Interface()

	// TODO: Research how to do this
	t.Skip("I have no idea how to construct a struct with an invalid field to provoke ErrInvalidFields")
}
