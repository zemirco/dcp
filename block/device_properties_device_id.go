package block

import "encoding/binary"

// DeviceID is a device id block.
type DeviceID struct {
	header
	VendorID uint16
	DeviceID uint16
}

var _ Block = &DeviceID{}

// UnmarshalBinary turns bytes into struct.
func (d *DeviceID) UnmarshalBinary(b []byte) error {
	if err := d.header.unmarshalBinary(b); err != nil {
		return err
	}

	i := d.header.len()
	d.VendorID = binary.BigEndian.Uint16(b[i : i+2])
	i += 2
	d.DeviceID = binary.BigEndian.Uint16(b[i : i+2])

	return nil
}

// MarshalBinary converts struct into byte slice.
func (d *DeviceID) MarshalBinary() ([]byte, error) {
	b := make([]byte, d.Len())

	bh, err := d.header.marshalBinary()
	if err != nil {
		return b, err
	}
	offset := 0

	copy(b[offset:], bh)
	offset += d.header.len()

	binary.BigEndian.PutUint16(b[offset:offset+2], d.VendorID)
	offset += 2

	binary.BigEndian.PutUint16(b[offset:offset+2], d.DeviceID)

	return b, nil
}

// Len returns length for name of station block.
func (d *DeviceID) Len() int {
	return d.header.len() + 2 + 2
}
