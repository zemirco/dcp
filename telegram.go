package dcp

import (
	"encoding/binary"
	"fmt"

	"github.com/zemirco/dcp/block"
	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/suboption"
)

// Telegram is a single telegram.
type Telegram struct {
	FrameID       FrameID
	ServiceID     ServiceID
	ServiceType   ServiceType
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
	ControlResponse      *block.ControlResponse
}

var _ block.Block = &Telegram{}

// UnmarshalBinary unmarshals a byte slice into a EthernetII.
func (t *Telegram) UnmarshalBinary(b []byte) error {
	i := 0

	t.FrameID = FrameID(binary.BigEndian.Uint16(b[i : i+2]))
	i += 2

	t.ServiceID = ServiceID(b[i])
	i++

	t.ServiceType = ServiceType(b[i])
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
		if blockLength%2 != 0 {
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
	if t.IPParameter != nil {
		size += t.IPParameter.Len()
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
		i += t.All.Len()
	}
	if t.IPParameter != nil {
		ipBytes, err := t.IPParameter.MarshalBinary()
		if err != nil {
			return b, err
		}
		copy(b[i:], ipBytes)
		i += t.IPParameter.Len()
	}

	return b, nil
}

// Len returns length.
func (t *Telegram) Len() int {
	length := 12
	if t.IPParameter != nil {
		length += t.IPParameter.Len()
	}
	if t.All != nil {
		length += t.All.Len()
	}
	return length
}

func (t *Telegram) decodeBlock(b []byte) int {
	opt := option.Option(b[0])
	fmt.Println("option", opt)

	subopt := suboption.Suboption(b[1])
	fmt.Println("suboption", subopt)

	length := binary.BigEndian.Uint16(b[2:4])
	fmt.Println("length", length)

	hasInfo := t.ServiceID == Identify && t.ServiceType == Response
	// hasQualifier := false

	switch {

	case opt == option.Properties && subopt == suboption.NameOfStation:

		t.NameOfStation = block.NewNameOfStation(hasInfo)
		if err := t.NameOfStation.UnmarshalBinary(b); err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", t.NameOfStation)
		fmt.Println(t.NameOfStation.NameOfStation)

	case opt == option.IP && subopt == suboption.IPParameter:

		t.IPParameter = &block.IPParameter{}
		t.IPParameter.HasInfo = hasInfo
		if err := t.IPParameter.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", t.IPParameter)
		fmt.Println(t.IPParameter.IPAddress, t.IPParameter.Subnetmask, t.IPParameter.StandardGateway)

	case opt == option.Properties && subopt == suboption.DeviceInstance:

		t.DeviceInstance = block.NewDeviceInstance(hasInfo)
		if err := t.DeviceInstance.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", t.DeviceInstance)
		fmt.Println(t.DeviceInstance.DeviceInstanceHigh, t.DeviceInstance.DeviceInstanceLow)

	case opt == option.Properties && subopt == suboption.ManufacturerSpecific:

		t.ManufacturerSpecific = block.NewManufacturerSpecific(hasInfo)
		if err := t.ManufacturerSpecific.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", t.ManufacturerSpecific)
		fmt.Println(t.ManufacturerSpecific.DeviceVendorValue)

	case opt == option.Initiative && subopt == suboption.DeviceInitiative:

		t.DeviceInitiative = block.NewDeviceInitiative(hasInfo)
		if err := t.DeviceInitiative.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", t.DeviceInitiative)
		fmt.Println(t.DeviceInitiative.Value)

	case opt == option.Control && subopt == suboption.Response:

		t.ControlResponse = block.NewControlResponse(hasInfo)
		if err := t.ControlResponse.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", t.ControlResponse)
		fmt.Println(t.ControlResponse.Error)
	}

	return 1 + 1 + 2 + int(length)
}
