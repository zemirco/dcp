package block

import (
	"encoding/binary"

	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/suboption"
)

// Header is block header.
type header struct {
	Option    option.Option
	Suboption suboption.Suboption
	Length    uint16

	HasInfo      bool
	Info         uint16
	HasQualifier bool
	Qualifier    uint16
}

// MarshalBinary converts struct into byte slice.
func (h *header) marshalBinary() ([]byte, error) {
	length := 4
	if h.HasInfo {
		length += 2
	}
	if h.HasQualifier {
		length += 2
	}

	b := make([]byte, length)

	offset := 0

	b[offset] = uint8(h.Option)
	offset++

	b[offset] = uint8(h.Suboption)
	offset++

	binary.BigEndian.PutUint16(b[offset:offset+2], h.Length)
	offset += 2

	if h.HasInfo {
		binary.BigEndian.PutUint16(b[offset:offset+2], h.Info)
		offset += 2
	}

	if h.HasQualifier {
		binary.BigEndian.PutUint16(b[offset:offset+2], h.Qualifier)
		offset += 2
	}

	return b, nil
}

// UnmarshalBinary turns bytes into struct.
func (h *header) unmarshalBinary(b []byte) error {
	offset := 0

	h.Option = option.Option(b[offset])
	offset++

	h.Suboption = suboption.Suboption(b[offset])
	offset++

	h.Length = binary.BigEndian.Uint16(b[offset : offset+2])
	offset += 2

	if h.HasInfo {
		h.Info = binary.BigEndian.Uint16(b[offset : offset+2])
		offset += 2
	}

	if h.HasQualifier {
		h.Qualifier = binary.BigEndian.Uint16(b[offset : offset+2])
		offset += 2
	}

	return nil
}

// Len returns block header length.
func (h *header) len() int {
	length := 4
	if h.HasInfo {
		length += 2
	}
	if h.HasQualifier {
		length += 2
	}
	return length
}
