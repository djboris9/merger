package merger

const (
	errDiffKindText          = "cannot resolve merge of (%s) and (%s)"
	errMergeUnsupportedText  = "merge of (%s) and (%s) is unsupported"
	errInvalidFieldsText     = "field with name (%s) is invalid in both structs"
	errDiffArrayTypesText    = "different types (%s) and (%s) for array"
	errDiffSliceTypesText    = "different types (%s) and (%s) for slice"
	errDiffMapKeyTypesText   = "different types (%s) and (%s) for map key"
	errDiffMapValueTypesText = "different types (%s) and (%s) for map value"
)

const (
	ErrDiffKind          = 1 << iota
	ErrMergeUnsupported  = 1 << iota
	ErrInvalidFields     = 1 << iota
	ErrDiffArrayTypes    = 1 << iota
	ErrDiffSliceTypes    = 1 << iota
	ErrDiffMapKeyTypes   = 1 << iota
	ErrDiffMapValueTypes = 1 << iota
)

// TODO: Make a way to determine where exactly it has failed
type MergeError struct {
	errString string
	errType   int
}

func (e *MergeError) Error() string {
	return e.errString
}

func (e *MergeError) Type() int {
	return e.errType
}
