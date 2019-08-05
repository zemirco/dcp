package frame

import (
	"encoding/binary"
	"fmt"
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
	GetSet           ID = 0xfefd
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

	// blocks
	All                  *block.All
	NameOfStation        *block.NameOfStation
	IPParameter          *block.IPParameter
	DeviceInstance       *block.DeviceInstance
	ManufacturerSpecific *block.ManufacturerSpecific
	DeviceInitiative     *block.DeviceInitiative
}

// UnmarshalBinary unmarshals a byte slice into a EthernetII.
func (t *Telegram) UnmarshalBinary(b []byte) error {
	i := 0

	t.FrameID = ID(binary.BigEndian.Uint16(b[i : i+2]))
	i += 2

	t.ServiceID = service.ID(b[i])
	i++

	t.ServiceType = service.Type(b[i])
	i++

	t.XID = binary.BigEndian.Uint32(b[i : i+4])
	i += 4

	t.ResponseDelay = binary.BigEndian.Uint16(b[i : i+2])
	i += 2

	t.DCPDataLength = binary.BigEndian.Uint16(b[i : i+2])
	i += 2

	length := int(t.DCPDataLength)
	offset := 0

	// fmt.Println("####", length)

	for length > 0 {
		blockLength := t.decodeBlock(b[i+offset:])

		// add padding for odd length block
		if blockLength%2 == 1 {
			blockLength++
		}

		length -= blockLength
		offset += blockLength
	}

	return nil
}

// MarshalBinary converts struct into byte slice.
func (t *Telegram) MarshalBinary() ([]byte, error) {
	size := 12
	if t.All != nil {
		size += t.All.Len()
	}

	b := make([]byte, size)
	i := 0

	binary.BigEndian.PutUint16(b[i:i+2], uint16(t.FrameID))
	i += 2

	b[i] = byte(t.ServiceID)
	i++

	b[i] = byte(t.ServiceType)
	i++

	binary.BigEndian.PutUint32(b[i:i+4], t.XID)
	i += 4

	binary.BigEndian.PutUint16(b[i:i+2], t.ResponseDelay)
	i += 2

	binary.BigEndian.PutUint16(b[i:i+2], t.DCPDataLength)
	i += 2

	if t.All != nil {
		allBytes, err := t.All.MarshalBinary()
		if err != nil {
			return b, err
		}
		copy(b[i:], allBytes)
	}

	return b, nil
}

// Len returns length.
func (t *Telegram) Len() int {
	return 12
}

func (t *Telegram) decodeBlock(b []byte) int {
	opt := option.Option(b[0])
	fmt.Println("option", opt)

	subopt := suboption.Suboption(b[1])
	fmt.Println("suboption", subopt)

	length := binary.BigEndian.Uint16(b[2:4])
	fmt.Println("length", length)

	switch {

	case opt == option.Properties && subopt == suboption.NameOfStation:

		var bnos block.NameOfStation
		if err := bnos.UnmarshalBinary(b); err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", bnos)
		fmt.Println(bnos.NameOfStation)

		t.NameOfStation = &bnos

	case opt == option.IP && subopt == suboption.IPParameter:

		var bip block.IPParameter
		if err := bip.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bip)
		fmt.Println(bip.IPAddress, bip.Subnetmask, bip.StandardGateway)

		t.IPParameter = &bip

	case opt == option.Properties && subopt == suboption.DeviceInstance:

		var bdi block.DeviceInstance
		if err := bdi.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bdi)
		fmt.Println(bdi.DeviceInstanceHigh, bdi.DeviceInstanceLow)

		t.DeviceInstance = &bdi

	case opt == option.Properties && subopt == suboption.ManufacturerSpecific:

		var bms block.ManufacturerSpecific
		if err := bms.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bms)
		fmt.Println(bms.DeviceVendorValue)

		t.ManufacturerSpecific = &bms

	case opt == option.Initiative && subopt == suboption.DeviceInitiative:

		var bdi block.DeviceInitiative
		if err := bdi.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bdi)
		fmt.Println(bdi.Value)

		t.DeviceInitiative = &bdi
	}

	return 1 + 1 + 2 + int(length)
}

// Frame is a single frame.
type Frame struct {
	EthernetII
	Telegram
}

// NewIdentifyRequest returns an identify request.
func NewIdentifyRequest(source net.HardwareAddr) *Frame {

	b := &block.All{
		Header: block.Header{
			Option:    option.All,
			Suboption: suboption.All,
		},
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
			All:           b,
		},
	}
}

// // NewSetRequest returns a set request.
// func NewSetRequest(source net.HardwareAddr, b block.Block) *Frame {
// 	return &Frame{
// 		EthernetII: EthernetII{
// 			Destination: []byte{0x01, 0x0e, 0xcf, 0x00, 0x00, 0x00},
// 			Source:      source,
// 			EtherType:   0x8892,
// 		},
// 		Telegram: Telegram{
// 			FrameID:       IdentifyRequest,
// 			ServiceID:     service.Identify,
// 			ServiceType:   service.Request,
// 			XID:           rand.Uint32(),
// 			ResponseDelay: 255,
// 			DCPDataLength: uint16(b.Len()),
// 		},
// 		All: b,
// 	}
// }

// MarshalBinary converts struct into byte slice.
func (f *Frame) MarshalBinary() ([]byte, error) {
	size := f.EthernetII.Len() + f.Telegram.Len()
	if f.All != nil {
		size += f.All.Len()
	}

	b := make([]byte, size)
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

	if err := f.Telegram.UnmarshalBinary(b[f.EthernetII.Len():]); err != nil {
		return err
	}

	return nil
}
