// Package merger can merge different Go types together.
//
// This can be useful for:
//   * Merging configurations from different sources
//   * Deduplicate items before serialization
//   * Setting field preferences
//
// Merge Algorithm
//
// Let A be the first arbitrary value and B be the second arbitrary value with precendence.
//
// boolean, numeric and string types will be overwritten by the argument with precendence (B).
// slice and array types will be concatenated (A ∥ B).
// struct and map types will be merged together giving a union of all fields, where the values of them are merged too (A ∪ B)
//
// Merge Examples
//
// Merging boolean, numeric, string types (overwrite):
//   A := "a"
//   B :=  4
//   V, _ := merger.Merge(A, B)
//   // V: 4
//
// Merging map, struct types (union):
//   A := map[string]int{
//       "x": 1,
//       "y": 2,
//   }
//   B := map[string]int{
//       "a": 1,
//       "b": 2,
//   }
//   V, _ := merger.Merge(A, B)
//   // V: map[string]int{"a": 1, "b": 2, "x": 1, "y": 2}
//
// Merging slice, array types (concat):
//   A := []int{1, 2, 3}
//   B := []int{4, 5, 6}
//   V, _ := merger.Merge(A, B)
//   // V: []int{1, 2, 3, 4, 5, 6}
//
// Merging embedded types:
//   A := map[string]interface{}{
//       "a": []int{1, 2},
//       "b": "aAaAaA",
//       "c": map[int]int{1: 10, 2: 20},
//   }
//   B := map[string]interface{}{
//       "a": []int{3, 4},
//       "b": "bBbBbB",
//       "c": map[int]int{1: 100, 3: 30},
//   }
//   V, _ := merger.Merge(A, B)
//   // V: map[string]interface{}{"a": []int{1,2,3,4}, "b": "bBbBbB", "c": map[int]int{1:100, 2:20, 3:30}}
package merger
