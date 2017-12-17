package merger

import (
	"fmt"
	"reflect"
)

// Merge merges a and b together where b has precedence.
// It is not destructive on their parameters, but these must not be modified while
// merge is in progress.
func Merge(a interface{}, b interface{}) (interface{}, error) {
	aKind := reflect.ValueOf(a).Kind()
	bKind := reflect.ValueOf(b).Kind()

	if kindContains(overwriteables, aKind) && kindContains(overwriteables, bKind) {
		// Is overwriteable
		return b, nil
	} else if kindContains(mergeables, aKind) && kindContains(mergeables, bKind) {
		// Is mergeable
		if aKind != bKind {
			return nil, &MergeError{
				errType:   ErrDiffKind,
				errString: fmt.Sprintf(errDiffKindText, aKind, bKind),
			}
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

	return nil, &MergeError{
		errType:   ErrMergeUnsupported,
		errString: fmt.Sprintf(errMergeUnsupportedText, aKind, bKind),
	}
}

// Field order is taken from a, additional fields of b are appended.
// No type assertion is made
// BUG(djboris): When merging structs fields are appended (A.Fields ∥ B.Fields).
// But field types should be taken from b as it has precedence.
func structMerge(a, b interface{}) (interface{}, error) {
	aV := reflect.ValueOf(a)
	bV := reflect.ValueOf(b)

	var commonFields []reflect.StructField

	// Add fields of A to commonFields
	for i := 0; i < aV.NumField(); i++ {
		x := aV.Type().Field(i)
		commonFields = append(commonFields, x)
	}

	// Add fields of B to commonFields
	for i := 0; i < bV.NumField(); i++ {
		x := bV.Type().Field(i)
		if !structFieldContains(commonFields, x.Name) {
			commonFields = append(commonFields, x)
		} else {
			// TODO Overwrite type, according to BUG for this function
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
			// Field exists in B
			resValue.FieldByName(fieldName).Set(fieldB)
		} else {
			return nil, &MergeError{
				errType:   ErrInvalidFields,
				errString: fmt.Sprintf(errInvalidFieldsText, fieldName),
			}
		}
	}

	return resValue.Interface(), nil
}

// arrayMerge returns effectively a slice
// Result: a ∥ b
// Type assertion only on Elem type
func arrayMerge(a, b interface{}) (interface{}, error) {
	aV := reflect.ValueOf(a)
	bV := reflect.ValueOf(b)

	if aV.Type().Elem() != bV.Type().Elem() {
		return nil, &MergeError{
			errType:   ErrDiffArrayTypes,
			errString: fmt.Sprintf(errDiffArrayTypesText, aV.Type().Elem(), bV.Type().Elem()),
		}
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

// Result: a ∪ b
// Type assertion on indexing and value type
func mapMerge(a, b interface{}) (interface{}, error) {
	aV := reflect.ValueOf(a)
	bV := reflect.ValueOf(b)
	if aV.Type().Key() != bV.Type().Key() {
		return nil, &MergeError{
			errType:   ErrDiffMapKeyTypes,
			errString: fmt.Sprintf(errDiffMapKeyTypesText, aV.Type().Key(), bV.Type().Key()),
		}
	}
	if aV.Type().Elem() != bV.Type().Elem() {
		return nil, &MergeError{
			errType:   ErrDiffMapValueTypes,
			errString: fmt.Sprintf(errDiffMapValueTypesText, aV.Type().Elem(), bV.Type().Elem()),
		}
	}

	res := reflect.MakeMap(aV.Type())

	// Copy everything from a
	for _, k := range aV.MapKeys() {
		res.SetMapIndex(k, aV.MapIndex(k))
	}

	for _, k := range bV.MapKeys() {
		if res.MapIndex(k).Kind() == reflect.Invalid {
			// If key was not in a already, add it from b
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
		return nil, &MergeError{
			errType:   ErrDiffSliceTypes,
			errString: fmt.Sprintf(errDiffSliceTypesText, aV.Type().Elem(), bV.Type().Elem()),
		}
	}
	return reflect.AppendSlice(aV, bV).Interface(), nil
}

// Boolean, numeric and string types (overwrite)
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

// Array and slice types (concat) + map and struct types (union)
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
