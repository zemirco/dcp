package dcp

import (
	"math/rand"
	"net"

	"github.com/zemirco/dcp/block"
	"github.com/zemirco/dcp/service"
)

// ID is two bytes.
type ID uint16

// Known frame ids.
const (
	IdentifyRequest  ID = 0xfefe
	IdentifyResponse ID = 0xfeff
	GetSet           ID = 0xfefd
)

// Frame is a single frame.
type Frame struct {
	EthernetII
	Telegram
}

var _ block.Block = &Frame{}

// NewIdentifyRequest returns an identify request.
func NewIdentifyRequest(source net.HardwareAddr) *Frame {

	b := block.NewAll()

	return &Frame{
		EthernetII: EthernetII{
			Destination: []byte{0x01, 0x0e, 0xcf, 0x00, 0x00, 0x00},
			Source:      source,
			EtherType:   0x8892,
		},
		Telegram: Telegram{
			FrameID:       IdentifyRequest,
			ServiceID:     service.Identify,
			ServiceType:   service.Request,
			XID:           rand.Uint32(),
			ResponseDelay: 255,
			DCPDataLength: uint16(b.Len()),
			All:           b,
		},
	}
}

// NewSetIPParameterRequest returns a set request.
func NewSetIPParameterRequest(dst, src net.HardwareAddr, b *block.IPParameter) *Frame {
	return &Frame{
		EthernetII: EthernetII{
			Destination: dst,
			Source:      src,
			EtherType:   0x8892,
		},
		Telegram: Telegram{
			FrameID:       GetSet,
			ServiceID:     service.Set,
			ServiceType:   service.Request,
			XID:           rand.Uint32(),
			ResponseDelay: 255,
			DCPDataLength: uint16(b.Len()),
			IPParameter:   b,
		},
	}
}

// MarshalBinary converts struct into byte slice.
func (f *Frame) MarshalBinary() ([]byte, error) {
	b := make([]byte, f.Len())
	i := 0

	ethernetIIBytes, err := f.EthernetII.MarshalBinary()
	if err != nil {
		return b, err
	}
	copy(b, ethernetIIBytes)
	i += f.EthernetII.Len()

	telegramBytes, err := f.Telegram.MarshalBinary()
	if err != nil {
		return b, err
	}
	copy(b[i:], telegramBytes)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a EthernetII.
func (f *Frame) UnmarshalBinary(b []byte) error {

	if err := f.EthernetII.UnmarshalBinary(b); err != nil {
		return err
	}

	return f.Telegram.UnmarshalBinary(b[f.EthernetII.Len():])
}

// Len returns length for name of station block.
func (f *Frame) Len() int {
	return f.EthernetII.Len() + f.Telegram.Len()
}
