package block

// DeviceInstance is a device instance block.
type DeviceInstance struct {
	header
	DeviceInstanceHigh uint8
	DeviceInstanceLow  uint8
}

var _ Block = &DeviceInstance{}

// NewDeviceInstance returns a new block.
func NewDeviceInstance(hasInfo bool) *DeviceInstance {
	return &DeviceInstance{
		header: header{
			HasInfo: hasInfo,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (d *DeviceInstance) UnmarshalBinary(b []byte) error {
	if err := d.header.unmarshalBinary(b); err != nil {
		return err
	}

	i := d.header.len()
	d.DeviceInstanceHigh = b[i]
	i++
	d.DeviceInstanceLow = b[i]

	return nil
}

// MarshalBinary converts struct into byte slice.
func (d *DeviceInstance) MarshalBinary() ([]byte, error) {
	b := make([]byte, d.Len())

	bh, err := d.header.marshalBinary()
	if err != nil {
		return b, err
	}
	offset := 0

	copy(b[offset:], bh)
	offset += d.header.len()

	b[offset] = d.DeviceInstanceHigh
	offset++

	b[offset] = d.DeviceInstanceLow

	return b, nil
}

// Len returns length for name of station block.
func (d *DeviceInstance) Len() int {
	return d.header.len() + 1 + 1
}
