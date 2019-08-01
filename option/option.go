package option

// Option is a single byte.
type Option byte

// Known options.
const (
	IP               Option = 1
	DeviceProperties Option = 2
	DeviceInitiative Option = 6
	All              Option = 255
)
