// Secret is a package for types that make it harder to accidentally expose
// secret variables when passing them around.
package secret

type String string

const redacted = "REDACTED"

// String implements fmt.Stringer and redacts the sensitive value.
func (s String) String() string {
	return redacted
}

// GoString implements fmt.GoStringer and redacts the sensitive value.
func (s String) GoString() string {
	return redacted
}

// Value returns the sensitive value as a string.
func (s String) Value() string {
	return string(s)
}

func (s String) MarshalJSON() ([]byte, error) {
	return []byte(`"` + redacted + `"`), nil
}
