package block

import (
	"encoding"
	"encoding/binary"
	"net"

	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/suboption"
)

// Block interface.
type Block interface {
	Len() int
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

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

// IPParameter is an ip parameter block.
type IPParameter struct {
	header
	IPAddress       net.IP
	Subnetmask      net.IP
	StandardGateway net.IP
}

var _ Block = &IPParameter{}

// NewIPParameter returns a new block.
func NewIPParameter(hasInfo, hasQualifier bool) *IPParameter {
	return &IPParameter{
		header: header{
			Option:       option.IP,
			Suboption:    suboption.IPParameter,
			HasInfo:      hasInfo,
			HasQualifier: hasQualifier,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (i *IPParameter) UnmarshalBinary(b []byte) error {
	if err := i.header.unmarshalBinary(b); err != nil {
		return err
	}

	o := i.header.len()
	i.IPAddress = net.IP(b[o : o+4])
	o += 4
	i.Subnetmask = net.IP(b[o : o+4])
	o += 4
	i.StandardGateway = net.IP(b[o : o+4])

	return nil
}

// MarshalBinary converts struct into byte slice.
func (i *IPParameter) MarshalBinary() ([]byte, error) {
	b := make([]byte, i.Len())

	bh, err := i.header.marshalBinary()
	if err != nil {
		return b, err
	}
	offset := 0

	copy(b[offset:], bh)
	offset += i.header.len()

	copy(b[offset:offset+4], i.IPAddress)
	offset += 4

	copy(b[offset:offset+4], i.Subnetmask)
	offset += 4

	copy(b[offset:offset+4], i.StandardGateway)

	return b, nil
}

// Len returns length for ip parameter block.
func (i *IPParameter) Len() int {
	return i.header.len() + 4 + 4 + 4
}

// NameOfStation is a name of station block.
type NameOfStation struct {
	header
	NameOfStation string
}

var _ Block = &NameOfStation{}

// NewNameOfStation returns a new block.
func NewNameOfStation(hasInfo bool) *NameOfStation {
	return &NameOfStation{
		header: header{
			HasInfo: hasInfo,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (n *NameOfStation) UnmarshalBinary(b []byte) error {
	if err := n.header.unmarshalBinary(b); err != nil {
		return err
	}

	i := n.header.len()
	n.NameOfStation = string(b[i : int(i)+int(n.header.Length)-2])

	return nil
}

// MarshalBinary converts struct into byte slice.
func (n *NameOfStation) MarshalBinary() ([]byte, error) {
	b := make([]byte, n.Len())

	bh, err := n.header.marshalBinary()
	if err != nil {
		return b, err
	}
	offset := 0

	copy(b[offset:], bh)
	offset += n.header.len()

	copy(b[offset:], n.NameOfStation)

	return b, nil
}

// Len returns length for name of station block.
func (n *NameOfStation) Len() int {
	return n.header.len() + len(n.NameOfStation)
}

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
