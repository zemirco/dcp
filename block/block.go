package block

import (
	"encoding"
)

// Block interface.
type Block interface {
	Len() int
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}
