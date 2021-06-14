package serial

import "io"

// Conversion represents a specific encoding (e.g. BSON, JSON, YAML).
// A conversion object must be
type Conversion interface {
	// TypeName gets the type name for the object to be unmarshaled.
	// The reader must be io.ReadSeeker to enable
	// the stream to be reset after acquiring the type name if necessary.
	// Implementations of TypeName must reset the stream before successful return.
	TypeName(reader io.ReadSeeker) (string, error)

	Decode(item interface{}, reader io.Reader) error
	Encode(item interface{}, writer io.Writer) error
}
