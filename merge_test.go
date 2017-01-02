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
			t.Error(err)
			continue
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
	var A [2]int = [2]int{1, 2}
	var B [3]int = [3]int{3, 4, 5}
	var Exp [5]int = [5]int{1, 2, 3, 4, 5}
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
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(res, Exp) {
		t.Errorf("Merge(%v, %v) => %v, want %v", A, B, res, Exp)
	} else {
		t.Logf("Merge(%v, %v) => %v", A, B, res)
	}
}

func TestTrace(t *testing.T) {
	A := []int{1, 3, 4}
	B := []int{4, 5, 6}
	Exp := "bbb"
	res, err := MergeTraced(A, B)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(res, Exp) {
		t.Errorf("Merge(%v, %v) => %v, want %v", A, B, res, Exp)
	} else {
		t.Logf("Merge(%v, %v) => %v", A, B, res)
	}
}
