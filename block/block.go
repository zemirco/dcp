package block

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
	DevicePropertiesDeviceVendor   OptionSuboption = 0x0201
	DevicePropertiesNameOfStation  OptionSuboption = 0x0202
	DevicePropertiesDeviceID       OptionSuboption = 0x0203
	DevicePropertiesDeviceRole     OptionSuboption = 0x0204
	DevicePropertiesDeviceOptions  OptionSuboption = 0x0205
	DevicePropertiesAliasName      OptionSuboption = 0x0206
	DevicePropertiesDeviceInstance OptionSuboption = 0x0207
	DevicePropertiesOEMDeviceID    OptionSuboption = 0x0208

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
