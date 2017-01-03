package merger

const (
	errDiffKindText          = "merge of (%s) and (%s) is invalid"
	errMergeUnsupportedText  = "merge of (%s) and (%s) is unsupported"
	errInvalidFieldsText     = "field with name (%s) is invalid in both structs"
	errDiffArrayTypesText    = "different types (%s) and (%s) for array"
	errDiffSliceTypesText    = "different types (%s) and (%s) for slice"
	errDiffMapKeyTypesText   = "different types (%s) and (%s) for map key"
	errDiffMapValueTypesText = "different types (%s) and (%s) for map value"
)

// MergeError types with value 2ⁿ
const (
	ErrDiffKind          = 1 << iota
	ErrMergeUnsupported  = 1 << iota
	ErrInvalidFields     = 1 << iota
	ErrDiffArrayTypes    = 1 << iota
	ErrDiffSliceTypes    = 1 << iota
	ErrDiffMapKeyTypes   = 1 << iota
	ErrDiffMapValueTypes = 1 << iota
)

// TODO: Make a way to determine where exactly it has failed. Perhaps together with tracer

// MergeError represents an error which has accured while merging
type MergeError struct {
	errString string
	errType   int
}

func (e *MergeError) Error() string {
	return e.errString
}

// Type returns the MergeError type which is one or multiple of Err constants by 2ⁿ
func (e *MergeError) Type() int {
	return e.errType
}
