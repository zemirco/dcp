package block

import "encoding/binary"

// DeviceInitiative is a manufacturer specific block.
type DeviceInitiative struct {
	header
	Value uint16
}

var _ Block = &DeviceInitiative{}

// NewDeviceInitiative returns a new block.
func NewDeviceInitiative(hasInfo bool) *DeviceInitiative {
	return &DeviceInitiative{
		header: header{
			HasInfo: hasInfo,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (d *DeviceInitiative) UnmarshalBinary(b []byte) error {
	if err := d.header.unmarshalBinary(b); err != nil {
		return err
	}

	i := d.header.len()
	d.Value = binary.BigEndian.Uint16(b[i : i+2])

	return nil
}

// MarshalBinary converts struct into byte slice.
func (d *DeviceInitiative) MarshalBinary() ([]byte, error) {
	b := make([]byte, d.Len())

	bh, err := d.header.marshalBinary()
	if err != nil {
		return b, err
	}
	offset := 0

	copy(b[offset:], bh)
	offset += d.header.len()

	binary.BigEndian.PutUint16(b[offset:offset+2], d.Value)

	return b, nil
}

// Len returns length for name of station block.
func (d *DeviceInitiative) Len() int {
	return d.header.len() + 2
}
