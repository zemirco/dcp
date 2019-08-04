package option

// Option is a single byte
type Option uint8

// All options
const (
	IP               Option = 0x01
	DeviceProperties Option = 0x02
	DHCP             Option = 0x03
	Control          Option = 0x05
	DeviceInitiative Option = 0x06
	All              Option = 0xFF
)
