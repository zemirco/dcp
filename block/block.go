package block

import (
	"encoding"
	"encoding/binary"
	"net"

	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/suboption"
)

// // Device describes a real device.
// type Device struct {
// 	HardwareAddr net.HardwareAddr
// 	IPParameter
// 	NameOfStation
// 	DeviceID
// 	DeviceInstance
// 	ManufacturerSpecific
// 	DeviceInitiative
// }

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
	Info      uint16
}

// MarshalBinary converts struct into byte slice.
func (h *Header) MarshalBinary() ([]byte, error) {
	length := 4
	if h.Length != 0 {
		length += 2
	}

	b := make([]byte, length)
	b[0] = uint8(h.Option)
	b[1] = uint8(h.Suboption)
	binary.BigEndian.PutUint16(b[2:4], h.Length)

	if h.Length != 0 {
		binary.BigEndian.PutUint16(b[4:6], h.Info)
	}

	return b, nil
}

// UnmarshalBinary turns bytes into struct.
func (h *Header) UnmarshalBinary(b []byte) error {
	h.Option = option.Option(b[0])
	h.Suboption = suboption.Suboption(b[1])
	h.Length = binary.BigEndian.Uint16(b[2:4])

	if h.Length != 0 {
		h.Info = binary.BigEndian.Uint16(b[4:6])
	}

	return nil
}

// Len returns block header length.
func (h *Header) Len() int {
	length := 4
	if h.Length != 0 {
		length += 2
	}
	return length
}

// IPParameter is an ip parameter block.
type IPParameter struct {
	Header
	IPAddress       net.IP
	Subnetmask      net.IP
	StandardGateway net.IP
}

// UnmarshalBinary turns bytes into struct.
func (i *IPParameter) UnmarshalBinary(b []byte) error {
	if err := i.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	i.IPAddress = net.IP(b[6:10])
	i.Subnetmask = net.IP(b[10:14])
	i.StandardGateway = net.IP(b[14:18])

	return nil
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

// UnmarshalBinary turns bytes into struct.
func (n *NameOfStation) UnmarshalBinary(b []byte) error {
	if err := n.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	n.NameOfStation = string(b[6 : 6+n.Header.Length-2])

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

	d.VendorID = binary.BigEndian.Uint16(b[6:8])
	d.DeviceID = binary.BigEndian.Uint16(b[8:10])

	return nil
}

// DeviceInstance is a device instance block.
type DeviceInstance struct {
	Header
	DeviceInstanceHigh uint8
	DeviceInstanceLow  uint8
}

// UnmarshalBinary turns bytes into struct.
func (d *DeviceInstance) UnmarshalBinary(b []byte) error {
	if err := d.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	d.DeviceInstanceHigh = b[6]
	d.DeviceInstanceLow = b[7]

	return nil
}

// ManufacturerSpecific is a manufacturer specific block.
type ManufacturerSpecific struct {
	Header
	DeviceVendorValue string
}

// UnmarshalBinary turns bytes into struct.
func (m *ManufacturerSpecific) UnmarshalBinary(b []byte) error {
	if err := m.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	m.DeviceVendorValue = string(b[6 : 6+m.Header.Length-2])

	return nil
}

// DeviceInitiative is a manufacturer specific block.
type DeviceInitiative struct {
	Header
	Value uint16
}

// UnmarshalBinary turns bytes into struct.
func (d *DeviceInitiative) UnmarshalBinary(b []byte) error {
	if err := d.Header.UnmarshalBinary(b); err != nil {
		return err
	}

	d.Value = binary.BigEndian.Uint16(b[6:8])

	return nil
}
