package block

import (
	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/suboption"
)

// ControlResponse is a control response block.
type ControlResponse struct {
	header
	Response  option.Option
	Suboption suboption.Suboption
	Error     uint8
}

var _ Block = &ControlResponse{}

// NewControlResponse returns a new block.
func NewControlResponse(hasInfo bool) *ControlResponse {
	return &ControlResponse{
		header: header{
			HasInfo: hasInfo,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (c *ControlResponse) UnmarshalBinary(b []byte) error {
	if err := c.header.unmarshalBinary(b); err != nil {
		return err
	}

	offset := c.header.len()

	c.Response = option.Option(b[offset])
	offset++

	c.Suboption = suboption.Suboption(b[offset])
	offset++

	c.Error = b[offset]

	return nil
}

// MarshalBinary converts struct into byte slice.
func (c *ControlResponse) MarshalBinary() ([]byte, error) {
	b := make([]byte, c.Len())

	bh, err := c.header.marshalBinary()
	if err != nil {
		return b, err
	}
	offset := 0

	copy(b[offset:], bh)
	offset += c.header.len()

	b[offset] = byte(c.Response)
	offset++

	b[offset] = byte(c.Suboption)
	offset++

	b[offset] = c.Error

	return b, nil
}

// Len returns length for name of station block.
func (c *ControlResponse) Len() int {
	return c.header.len() + 1 + 1 + 1
}
