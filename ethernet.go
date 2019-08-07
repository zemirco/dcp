package dcp

import (
	"encoding/binary"
	"net"

	"github.com/zemirco/dcp/block"
)

// EthernetII header.
type EthernetII struct {
	Destination net.HardwareAddr
	Source      net.HardwareAddr
	EtherType   uint16
}

// NewEthernetII returns pointer to ethernet II struct.
func NewEthernetII(dst, src net.HardwareAddr) *EthernetII {
	return &EthernetII{
		Destination: dst,
		Source:      src,
		EtherType:   0x8892,
	}
}

var _ block.Block = &EthernetII{}

// MarshalBinary converts struct into byte slice.
func (e *EthernetII) MarshalBinary() ([]byte, error) {
	b := make([]byte, 14)

	copy(b[0:6], e.Destination)
	copy(b[6:12], e.Source)
	binary.BigEndian.PutUint16(b[12:14], e.EtherType)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a EthernetII.
func (e *EthernetII) UnmarshalBinary(b []byte) error {

	e.Destination = b[0:6]
	e.Source = b[6:12]
	e.EtherType = binary.BigEndian.Uint16(b[12:14])

	return nil
}

// Len returns length.
func (e *EthernetII) Len() int {
	return 14
}
