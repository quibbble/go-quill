package errors

// common errors
var (
	ErrIndexOutOfBounds    = Errorf("index out of bounds")
	ErrInvalidSliceLength  = Errorf("slice has invalid length")
	ErrInterfaceConversion = Errorf("failed to convert interface")
	ErrMissingMapKey       = Errorf("key not found in map")
)