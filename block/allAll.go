package block

import (
	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/suboption"
)

// All is an all block.
type All struct {
	header
}

var _ Block = &All{}

// NewAll returns new all block.
func NewAll() *All {
	return &All{
		header: header{
			Option:    option.All,
			Suboption: suboption.All,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (a *All) UnmarshalBinary(b []byte) error {
	return a.header.unmarshalBinary(b)
}

// MarshalBinary converts struct into byte slice.
func (a *All) MarshalBinary() ([]byte, error) {
	return a.header.marshalBinary()
}

// Len returns length for ip parameter block.
func (a *All) Len() int {
	return a.header.len()
}
