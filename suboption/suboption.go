package suboption

// Suboption is a single byte
type Suboption uint8

// All suboptions
const (
	MACAddress  Suboption = 0x01
	IPParameter Suboption = 0x02
	FullIPSuite Suboption = 0x03

	ManufacturerSpecific Suboption = 0x01
	NameOfStation        Suboption = 0x02
	DeviceID             Suboption = 0x03
	DeviceRole           Suboption = 0x04
	DeviceOptions        Suboption = 0x05
	AliasName            Suboption = 0x06
	DeviceInstance       Suboption = 0x07
	OEMDeviceID          Suboption = 0x08

	HostName                  Suboption = 0x0C
	VendorSpecificInformation Suboption = 0x2B
	ServerIdentifier          Suboption = 0x36
	ParameterRequestList      Suboption = 0x37
	ClassIdentifier           Suboption = 0x3C
	DHCPClientIdentifier      Suboption = 0x3D
	FullyQualifiedDomainName  Suboption = 0x51
	UUIDClientIdentifier      Suboption = 0x61
	DHCP                      Suboption = 0xFF

	Start          Suboption = 0x01
	Stop           Suboption = 0x02
	Signal         Suboption = 0x03
	Response       Suboption = 0x04
	FactoryReset   Suboption = 0x05
	ResetToFactory Suboption = 0x06

	DeviceInitiative Suboption = 0x01

	AllSelector Suboption = 0xFF
)
