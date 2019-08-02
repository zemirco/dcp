package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"syscall"

	"github.com/zemirco/dcp/block"
	"github.com/zemirco/dcp/frame"
	"github.com/zemirco/dcp/service"
)

type ethernetII struct {
	Destination [6]byte
	Source      [6]byte
	EtherType   uint16
}

type telegram struct {
	FrameID       frame.ID
	ServiceID     service.ID
	ServiceType   service.Type
	XID           uint32
	ResponseDelay uint16
}

var destination = [6]byte{
	0x01, 0x0e, 0xcf, 0x00, 0x00, 0x00,
}

// var request = []byte{
// 	0x01, 0x0e, 0xcf, 0x00, 0x00, 0x00, 0xa4, 0x4c,
// 	0xc8, 0xe6, 0xd7, 0x89, 0x88, 0x92, 0xfe, 0xfe,
// 	0x05, 0x00, 0x00, 0x00, 0x07, 0xe8, 0x00, 0xff,
// 	0x00, 0x04, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
// 	0x00, 0x00, 0x00, 0x00,
// }

const etherType uint16 = 0x8892

// host order (usually little endian) -> network order (big endian)
func htons(n int) int {
	return int(int16(byte(n))<<8 | int16(byte(n>>8)))
}

func main() {

	ifname := "enxa44cc8e54721"

	f := make([]byte, 30)

	interf, err := net.InterfaceByName(ifname)
	if err != nil {
		panic(err)
	}

	var source [6]byte
	copy(source[:], interf.HardwareAddr)

	e := ethernetII{
		Destination: destination,
		Source:      source,
		EtherType:   etherType,
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, e)
	copy(f[0:], buf.Bytes())

	t := telegram{
		FrameID:       frame.IdentifyRequest,
		ServiceID:     service.Identify,
		ServiceType:   service.Request,
		XID:           rand.Uint32(),
		ResponseDelay: 255,
	}

	buf.Reset()
	binary.Write(&buf, binary.BigEndian, t)
	copy(f[14:], buf.Bytes())

	// dcp data length

	b := &block.Header{
		Option:    255,
		Suboption: 255,
	}

	buf.Reset()
	binary.Write(&buf, binary.BigEndian, b)

	// +2 because DCPBlockLength
	binary.BigEndian.PutUint16(f[24:26], uint16(len(buf.Bytes()))+2)

	copy(f[26:28], buf.Bytes())

	binary.BigEndian.PutUint16(f[28:30], 0)

	log.Printf("% x\n", f)

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, htons(0x8892))

	if err != nil {
		panic(err)
	}

	defer syscall.Close(fd)

	addr := syscall.SockaddrLinklayer{
		Ifindex: interf.Index,
	}

	if err := syscall.Sendto(fd, f, 0, &addr); err != nil {
		panic(err)
	}

	// start reading incoming data
	for {
		buffer := make([]byte, 256)

		n, from, err := syscall.Recvfrom(fd, buffer, 0)
		if err != nil {
			panic(err)
		}

		fmt.Println(n)
		fmt.Println(from)

		// fmt.Printf("% x\n", buffer[:n])

		e := ethernetII{}
		binary.Read(bytes.NewBuffer(buffer[0:]), binary.BigEndian, &e)
		fmt.Printf("%#v\n", e)

		frameID := buffer[14:16]
		fmt.Printf("frame id % x\n", frameID)

		serviceID := buffer[16:17]
		fmt.Printf("service id % x\n", serviceID)

		serviceType := buffer[17:18]
		fmt.Printf("service type % x\n", serviceType)

		xid := buffer[18:22]
		fmt.Printf("xid % x\n", xid)

		reserved := buffer[22:24]
		fmt.Printf("reserved % x\n", reserved)

		// dcpDataLength := buffer[24:26]
		// fmt.Printf("dcp data length % x\n", dcpDataLength)

		length := int(binary.BigEndian.Uint16(buffer[24:26]))
		fmt.Println("length", length)

		offset := 0

		for length > 0 {
			blockLength := decodeBlock(buffer[26+offset:])

			// add padding for odd length block
			if blockLength%2 == 1 {
				blockLength++
			}

			length -= blockLength
			offset += blockLength
		}

	}

}

func decodeBlock(b []byte) int {
	optSubopt := block.OptionSuboption(binary.BigEndian.Uint16(b[0:2]))
	fmt.Println("option suboption", optSubopt)

	length := binary.BigEndian.Uint16(b[2:4])
	fmt.Println("length", length)

	// switch b := event.(type) {
	// case *block.NameOfStation:
	// 	fmt.Println(b.NameOfStation)
	// case *block.IPParameter:
	// 	fmt.Println(b.IPAddress, b.Subnetmask, b.StandardGateway)
	// }

	switch optSubopt {

	case block.DevicePropertiesNameOfStation:

		var bnos block.NameOfStation
		if err := bnos.Unmarshal(b); err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", bnos)
		fmt.Println(bnos.NameOfStation)

	case block.IPIPParameter:

		var bip block.IPParameter
		if err := bip.Unmarshal(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bip)
		fmt.Println(bip.IPAddress, bip.Subnetmask, bip.StandardGateway)

	case block.DevicePropertiesDeviceInstance:

		var bdi block.DeviceInstance
		if err := bdi.Unmarshal(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bdi)
		fmt.Println(bdi.DeviceInstanceHigh, bdi.DeviceInstanceLow)

	case block.DevicePropertiesManufacturerSpecific:

		var bms block.ManufacturerSpecific
		if err := bms.Unmarshal(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bms)
		fmt.Println(bms.DeviceVendorValue)

	case block.DeviceInitiativeDeviceInitiative:

		var bdi block.DeviceInitiative
		if err := bdi.Unmarshal(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bdi)
		fmt.Println(bdi.Value)
	}

	return 1 + 1 + 2 + int(length)
}
