package service

// ID is a single byte.
type ID byte

// Known ids.
const (
	Get      ID = 1
	Set      ID = 4
	Identify ID = 5
)

// Type is a single byte.
type Type byte

// Known types.
const (
	Request  Type = 0
	Response Type = 1
)
