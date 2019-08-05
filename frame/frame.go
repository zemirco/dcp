package frame

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"

	"github.com/zemirco/dcp/block"
	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/service"
	"github.com/zemirco/dcp/suboption"
)

// ID is two bytes.
type ID uint16

// Known frame ids.
const (
	IdentifyRequest  ID = 0xfefe
	IdentifyResponse ID = 0xfeff
)

// EthernetII header.
type EthernetII struct {
	Destination net.HardwareAddr
	Source      net.HardwareAddr
	EtherType   uint16
}

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

// Telegram is a single telegram.
type Telegram struct {
	FrameID       ID
	ServiceID     service.ID
	ServiceType   service.Type
	XID           uint32
	ResponseDelay uint16
	DCPDataLength uint16
}

// UnmarshalBinary unmarshals a byte slice into a EthernetII.
func (t *Telegram) UnmarshalBinary(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.BigEndian, t)
}

// MarshalBinary converts struct into byte slice.
func (t *Telegram) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, t)
	return b.Bytes(), nil
}

// Len returns length.
func (t *Telegram) Len() int {
	return 12
}

// Frame is a single frame.
type Frame struct {
	EthernetII
	Telegram
	Block block.Block
}

// NewIdentifyRequest returns an identify request.
func NewIdentifyRequest(source net.HardwareAddr) *Frame {

	b := &block.Header{
		Option:    option.All,
		Suboption: suboption.All,
	}

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
		},
		Block: b,
	}
}

// MarshalBinary converts struct into byte slice.
func (f *Frame) MarshalBinary() ([]byte, error) {
	size := f.EthernetII.Len() + f.Telegram.Len() + f.Block.Len()
	b := make([]byte, size)

	ethernetIIBytes, err := f.EthernetII.MarshalBinary()
	if err != nil {
		return b, err
	}
	copy(b, ethernetIIBytes)

	telegramBytes, err := f.Telegram.MarshalBinary()
	if err != nil {
		return b, err
	}
	copy(b[f.EthernetII.Len():], telegramBytes)

	blockBytes, err := f.Block.MarshalBinary()
	if err != nil {
		return b, err
	}
	copy(b[f.EthernetII.Len()+f.Telegram.Len():], blockBytes)

	return b, nil
}

// UnmarshalBinary unmarshals a byte slice into a EthernetII.
func (f *Frame) UnmarshalBinary(b []byte) error {

	if err := f.EthernetII.UnmarshalBinary(b); err != nil {
		return err
	}

	if err := f.Telegram.UnmarshalBinary(b[f.EthernetII.Len():]); err != nil {
		return err
	}

	// todo: continue here

	// e.Destination = b[0:6]
	// e.Source = b[6:12]
	// e.EtherType = binary.BigEndian.Uint16(b[12:14])

	return nil
}
