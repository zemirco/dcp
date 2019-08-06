package block

// ManufacturerSpecific is a manufacturer specific block.
type ManufacturerSpecific struct {
	header
	DeviceVendorValue string
}

var _ Block = &ManufacturerSpecific{}

// NewManufacturerSpecific returns a new block.
func NewManufacturerSpecific(hasInfo bool) *ManufacturerSpecific {
	return &ManufacturerSpecific{
		header: header{
			HasInfo: hasInfo,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (m *ManufacturerSpecific) UnmarshalBinary(b []byte) error {
	if err := m.header.unmarshalBinary(b); err != nil {
		return err
	}

	i := m.header.len()
	m.DeviceVendorValue = string(b[i : i+int(m.header.Length)-2])

	return nil
}

// MarshalBinary converts struct into byte slice.
func (m *ManufacturerSpecific) MarshalBinary() ([]byte, error) {
	b := make([]byte, m.Len())

	bh, err := m.header.marshalBinary()
	if err != nil {
		return b, err
	}
	offset := 0

	copy(b[offset:], bh)
	offset += m.header.len()

	copy(b[offset:], m.DeviceVendorValue)

	return b, nil
}

// Len returns length for name of station block.
func (m *ManufacturerSpecific) Len() int {
	return m.header.len() + len(m.DeviceVendorValue)
}