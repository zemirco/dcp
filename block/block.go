package block

import (
	"encoding/binary"
	"net"
)

// Device describes a real device.
type Device struct {
	HardwareAddr net.HardwareAddr
	IPParameter
	NameOfStation
	DeviceID
	DeviceInstance
	ManufacturerSpecific
	DeviceInitiative
}

// OptionSuboption is two bytes long. First byte is option and second byte is suboption.
// Combined they are easier to use since suboption reuses values depending on the option.
// suboption == 1 is mac address for option ip
// suboption == 1 is device vendor for option device properties.
type OptionSuboption uint16

// Option and Suboption combined.
const (
	// option IP
	IPMACAddress  OptionSuboption = 0x0101
	IPIPParameter OptionSuboption = 0x0102
	IPFullIPSuite OptionSuboption = 0x0103

	// option DeviceProperties
	DevicePropertiesManufacturerSpecific OptionSuboption = 0x0201
	DevicePropertiesNameOfStation        OptionSuboption = 0x0202
	DevicePropertiesDeviceID             OptionSuboption = 0x0203
	DevicePropertiesDeviceRole           OptionSuboption = 0x0204
	DevicePropertiesDeviceOptions        OptionSuboption = 0x0205
	DevicePropertiesAliasName            OptionSuboption = 0x0206
	DevicePropertiesDeviceInstance       OptionSuboption = 0x0207
	DevicePropertiesOEMDeviceID          OptionSuboption = 0x0208

	// option DHCP
	DHCPHostName                  OptionSuboption = 0x030C
	DHCPVendorSpecificInformation OptionSuboption = 0x032B
	DHCPServerIdentifier          OptionSuboption = 0x0336
	DHCPParameterRequestList      OptionSuboption = 0x0337
	DHCPClassIdentifier           OptionSuboption = 0x033C
	DHCPDHCPClientIdentifier      OptionSuboption = 0x033D
	DHCPFullyQualifiedDomainName  OptionSuboption = 0x0351
	DHCPUUIDClientIdentifier      OptionSuboption = 0x0361
	DHCPDHCP                      OptionSuboption = 0x03FF

	// option Control
	ControlStart          OptionSuboption = 0x0501
	ControlStop           OptionSuboption = 0x0502
	ControlSignal         OptionSuboption = 0x0503
	ControlResponse       OptionSuboption = 0x0504
	ControlFactoryReset   OptionSuboption = 0x0505
	ControlResetToFactory OptionSuboption = 0x0506

	// option DeviceInitiative
	DeviceInitiativeDeviceInitiative OptionSuboption = 0x0601

	// option All
	AllSelectorAllSelector OptionSuboption = 0xFFFF
)

// Something is common interface.
type Something interface {
	Unmarshal(b []byte) error
	Len() int
}

// Header is block header.
type Header struct {
	Option    uint8
	Suboption uint8
	Length    uint16
	Info      uint16
}

// Unmarshal turns bytes into struct.
func (h *Header) Unmarshal(b []byte) error {
	h.Option = b[0]
	h.Suboption = b[1]
	h.Length = binary.BigEndian.Uint16(b[2:4])
	h.Info = binary.BigEndian.Uint16(b[4:6])
	return nil
}

// Len returns block header length.
func (h *Header) Len() int {
	return 6
}

// IPParameter is an ip parameter block.
type IPParameter struct {
	Header
	IPAddress       net.IP
	Subnetmask      net.IP
	StandardGateway net.IP
}

// Unmarshal turns bytes into struct.
func (i *IPParameter) Unmarshal(b []byte) error {
	if err := i.Header.Unmarshal(b); err != nil {
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

// Unmarshal turns bytes into struct.
func (n *NameOfStation) Unmarshal(b []byte) error {
	if err := n.Header.Unmarshal(b); err != nil {
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

// Unmarshal turns bytes into struct.
func (d *DeviceID) Unmarshal(b []byte) error {
	if err := d.Header.Unmarshal(b); err != nil {
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

// Unmarshal turns bytes into struct.
func (d *DeviceInstance) Unmarshal(b []byte) error {
	if err := d.Header.Unmarshal(b); err != nil {
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

// Unmarshal turns bytes into struct.
func (m *ManufacturerSpecific) Unmarshal(b []byte) error {
	if err := m.Header.Unmarshal(b); err != nil {
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

// Unmarshal turns bytes into struct.
func (d *DeviceInitiative) Unmarshal(b []byte) error {
	if err := d.Header.Unmarshal(b); err != nil {
		return err
	}

	d.Value = binary.BigEndian.Uint16(b[6:8])

	return nil
}
