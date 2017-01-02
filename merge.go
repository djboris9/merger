package merger

import (
	"errors"
	"fmt"
	"reflect"
)

// Merge merges fields from `b` into `a`
func Merge(a interface{}, b interface{}) (interface{}, error) {
	aKind := reflect.ValueOf(a).Kind()
	bKind := reflect.ValueOf(b).Kind()

	if kindContains(overwriteables, aKind) && kindContains(overwriteables, bKind) {
		// Is overwriteable
		return b, nil
	} else if kindContains(mergeables, aKind) && kindContains(mergeables, bKind) {
		// Is mergeable
		if aKind != bKind {
			return nil, errors.New("For merge, aKind and bKind must be same")
		}

		if aKind == reflect.Array {
			return arrayMerge(a, b)
		} else if aKind == reflect.Slice {
			return sliceMerge(a, b)
		} else if aKind == reflect.Struct {
			return structMerge(a, b)
		} else if aKind == reflect.Map {
			return mapMerge(a, b)
		}
	}

	return nil, fmt.Errorf("Merge of (%v) and (%v) not supported", aKind, bKind)
}

// Field types are taken from `a`
// No type assertion is made
func structMerge(a, b interface{}) (interface{}, error) {
	aV := reflect.ValueOf(a)
	bV := reflect.ValueOf(b)

	var commonFields []reflect.StructField

	// Add fields of A to commonFields
	for i := 0; i < aV.NumField(); i++ {
		x := aV.Type().Field(i)
		if !structFieldContains(commonFields, x.Name) {
			commonFields = append(commonFields, x)
		}
	}

	// Add fields of B to commonFields
	for i := 0; i < bV.NumField(); i++ {
		x := bV.Type().Field(i)
		if !structFieldContains(commonFields, x.Name) {
			commonFields = append(commonFields, x)
		}
	}

	// Construct output struct
	resType := reflect.StructOf(commonFields)
	resValue := reflect.New(resType).Elem()

	for i := range commonFields {
		fieldName := commonFields[i].Name
		fieldA := aV.FieldByName(fieldName)
		fieldB := bV.FieldByName(fieldName)

		if fieldA.IsValid() && fieldB.IsValid() {
			// Field exists in A and B, do a merge
			m, err := Merge(fieldA.Interface(), fieldB.Interface())
			if err != nil {
				return nil, err
			}
			resValue.FieldByName(fieldName).Set(reflect.ValueOf(m))
		} else if fieldA.IsValid() {
			// Field exists in A
			resValue.FieldByName(fieldName).Set(fieldA)
		} else if fieldB.IsValid() {
			// Field exists in A
			resValue.FieldByName(fieldName).Set(fieldB)
		} else {
			return nil, fmt.Errorf("both fields are invalid: (%v), (%v)", fieldA, fieldB)
		}
	}

	return resValue.Interface(), nil
}

// arrayMerge returns effectively a slice
// Result: a || b
// Type assertion only on Elem type
func arrayMerge(a, b interface{}) (interface{}, error) {
	aV := reflect.ValueOf(a)
	bV := reflect.ValueOf(b)

	if aV.Type().Elem() != bV.Type().Elem() {
		return nil, errors.New("Not the same Elem type")
	}

	resType := reflect.ArrayOf(aV.Len()+bV.Len(), aV.Type().Elem())
	resValue := reflect.New(resType).Elem()

	for i := 0; i < aV.Len(); i++ {
		resValue.Index(i).Set(aV.Index(i))
	}
	for i := 0; i < bV.Len(); i++ {
		resValue.Index(aV.Len() + i).Set(bV.Index(i))
	}

	return resValue.Interface(), nil
}

// Result: a || b
// Type assertion on indexing and value type
func mapMerge(a, b interface{}) (interface{}, error) {
	aV := reflect.ValueOf(a)
	bV := reflect.ValueOf(b)
	if aV.Type().Key() != bV.Type().Key() {
		return nil, errors.New("Key type is different")
	}
	if aV.Type().Elem() != bV.Type().Elem() {
		return nil, errors.New("Value type is different")
	}

	res := reflect.MakeMap(aV.Type())

	// Copy everything from a
	for _, k := range aV.MapKeys() {
		res.SetMapIndex(k, aV.MapIndex(k))
	}

	for _, k := range bV.MapKeys() {
		if res.MapIndex(k).Kind() == reflect.Invalid {
			// If key was not in `a` already, add it from `b`
			res.SetMapIndex(k, bV.MapIndex(k))
		} else {
			// Need to merge both values
			resVal := res.MapIndex(k).Interface()
			bVal := bV.MapIndex(k).Interface()
			m, err := Merge(resVal, bVal)
			if err != nil {
				return nil, err
			}
			res.SetMapIndex(k, reflect.ValueOf(m))
		}
	}

	return res.Interface(), nil
}

// Type assertion only on Elem type
func sliceMerge(a, b interface{}) (interface{}, error) {
	aV := reflect.ValueOf(a)
	bV := reflect.ValueOf(b)

	if aV.Type().Elem() != bV.Type().Elem() {
		return nil, errors.New("a or b has not the same type")
	}
	return reflect.AppendSlice(aV, bV).Interface(), nil
}

var overwriteables = []reflect.Kind{
	reflect.Bool,
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Uintptr,
	reflect.Float32,
	reflect.Float64,
	reflect.Complex64,
	reflect.Complex128,
	reflect.String,
}

var mergeables = []reflect.Kind{
	reflect.Array,
	reflect.Map,
	reflect.Slice,
	reflect.Struct,
}

func structFieldContains(s []reflect.StructField, n string) bool {
	for _, z := range s {
		if z.Name == n {
			return true
		}
	}
	return false
}

func kindContains(s []reflect.Kind, k reflect.Kind) bool {
	for _, z := range s {
		if z == k {
			return true
		}
	}
	return false
}
