package dcp

// ServiceID is a single byte.
type ServiceID byte

// Known ids.
const (
	Get      ServiceID = 1
	Set      ServiceID = 4
	Identify ServiceID = 5
)

// ServiceType is a single byte.
type ServiceType byte

// Known types.
const (
	Request  ServiceType = 0
	Response ServiceType = 1
)
