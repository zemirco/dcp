package option

// // Option is a single byte.
// type Option byte

// // Known options.
// const (
// 	IP               Option = 1
// 	DeviceProperties Option = 2
// 	DeviceInitiative Option = 6
// 	All              Option = 255
// )

// Option is a single byte
type Option uint8

// Suboption is a single byte
type Suboption uint8

// All options
const (
	IP               Option = 1
	DeviceProperties Option = 2
)

// All suboptions
const (
	MACAddress  Suboption = 1
	IPParameter Suboption = 2
	FullIPSuite Suboption = 3

	ManufacturerSpecific Suboption = 1
	NameOfStation        Suboption = 2
)
