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
type Header struct {
	Option    option.Option
	Suboption suboption.Suboption
	Length    uint16

	HasInfo      bool
	Info         uint16
	HasQualifier bool
	Qualifier    uint16
}

// MarshalBinary converts struct into byte slice.
func (h *Header) MarshalBinary() ([]byte, error) {
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
func (h *Header) UnmarshalBinary(b []byte) error {
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
func (h *Header) Len() int {
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
	Header
}

// UnmarshalBinary turns bytes into struct.
func (a *All) UnmarshalBinary(b []byte) error {
	return a.Header.UnmarshalBinary(b)
}

// MarshalBinary converts struct into byte slice.
func (a *All) MarshalBinary() ([]byte, error) {
	return a.Header.MarshalBinary()
}

// Len returns length for ip parameter block.
func (a *All) Len() int {
	return a.Header.Len()
}

// IPParameter is an ip parameter block.
type IPParameter struct {
	Header
	IPAddress       net.IP
	Subnetmask      net.IP
	StandardGateway net.IP
}

// NewIPParameter returns a new block.
func NewIPParameter(hasInfo, hasQualifier bool) *IPParameter {
	return &IPParameter{
		Header: Header{
			Option:       option.IP,
			Suboption:    suboption.IPParameter,
			HasInfo:      hasInfo,
			HasQualifier: hasQualifier,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (i *IPParameter) UnmarshalBinary(b []byte) error {
	if err := i.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	o := i.Header.Len()
	i.IPAddress = net.IP(b[o : o+4])
	o += 4
	i.Subnetmask = net.IP(b[o : o+4])
	o += 4
	i.StandardGateway = net.IP(b[o : o+4])

	return nil
}

// MarshalBinary converts struct into byte slice.
func (i *IPParameter) MarshalBinary() ([]byte, error) {
	size := i.Header.Len() + i.Len()

	b := make([]byte, size)

	bh, err := i.Header.MarshalBinary()
	if err != nil {
		return b, err
	}
	offset := 0

	copy(b[offset:], bh)
	offset += i.Header.Len()

	copy(b[offset:offset+4], i.IPAddress)
	offset += 4

	copy(b[offset:offset+4], i.Subnetmask)
	offset += 4

	copy(b[offset:offset+4], i.StandardGateway)

	return b, nil
}

// Len returns length for ip parameter block.
func (i *IPParameter) Len() int {
	return i.Header.Len() + 4 + 4 + 4
}

// NameOfStation is a name of station block.
type NameOfStation struct {
	Header
	NameOfStation string
}

// NewNameOfStation returns a new block.
func NewNameOfStation(hasInfo bool) *NameOfStation {
	return &NameOfStation{
		Header: Header{
			HasInfo: hasInfo,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (n *NameOfStation) UnmarshalBinary(b []byte) error {
	if err := n.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	i := n.Header.Len()
	n.NameOfStation = string(b[i : int(i)+int(n.Header.Length)-2])

	return nil
}

// Len returns length for name of station block.
func (n *NameOfStation) Len() int {
	return n.Header.Len() + len(n.NameOfStation)
}

// DeviceID is a device id block.
type DeviceID struct {
	Header
	VendorID uint16
	DeviceID uint16
}

// UnmarshalBinary turns bytes into struct.
func (d *DeviceID) UnmarshalBinary(b []byte) error {
	if err := d.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	i := d.Header.Len()
	d.VendorID = binary.BigEndian.Uint16(b[i : i+2])
	i += 2
	d.DeviceID = binary.BigEndian.Uint16(b[i : i+2])

	return nil
}

// DeviceInstance is a device instance block.
type DeviceInstance struct {
	Header
	DeviceInstanceHigh uint8
	DeviceInstanceLow  uint8
}

// NewDeviceInstance returns a new block.
func NewDeviceInstance(hasInfo bool) *DeviceInstance {
	return &DeviceInstance{
		Header: Header{
			HasInfo: hasInfo,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (d *DeviceInstance) UnmarshalBinary(b []byte) error {
	if err := d.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	i := d.Header.Len()
	d.DeviceInstanceHigh = b[i]
	i++
	d.DeviceInstanceLow = b[i]

	return nil
}

// ManufacturerSpecific is a manufacturer specific block.
type ManufacturerSpecific struct {
	Header
	DeviceVendorValue string
}

// NewManufacturerSpecific returns a new block.
func NewManufacturerSpecific(hasInfo bool) *ManufacturerSpecific {
	return &ManufacturerSpecific{
		Header: Header{
			HasInfo: hasInfo,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (m *ManufacturerSpecific) UnmarshalBinary(b []byte) error {
	if err := m.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	i := m.Header.Len()
	m.DeviceVendorValue = string(b[i : i+int(m.Header.Length)-2])

	return nil
}

// DeviceInitiative is a manufacturer specific block.
type DeviceInitiative struct {
	Header
	Value uint16
}

// NewDeviceInitiative returns a new block.
func NewDeviceInitiative(hasInfo bool) *DeviceInitiative {
	return &DeviceInitiative{
		Header: Header{
			HasInfo: hasInfo,
		},
	}
}

// UnmarshalBinary turns bytes into struct.
func (d *DeviceInitiative) UnmarshalBinary(b []byte) error {
	if err := d.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	i := d.Header.Len()
	d.Value = binary.BigEndian.Uint16(b[i : i+2])

	return nil
}
