package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"syscall"

	"github.com/zemirco/dcp/block"
	"github.com/zemirco/dcp/frame"
	"github.com/zemirco/dcp/option"
	"github.com/zemirco/dcp/suboption"
)

var destination = []byte{
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

	interf, err := net.InterfaceByName(ifname)
	if err != nil {
		panic(err)
	}

	f := frame.NewIdentifyRequest(interf.HardwareAddr)
	b, err := f.MarshalBinary()
	if err != nil {
		panic(err)
	}

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, htons(0x8892))

	if err != nil {
		panic(err)
	}

	defer syscall.Close(fd)

	addr := syscall.SockaddrLinklayer{
		Ifindex: interf.Index,
	}

	if err := syscall.Sendto(fd, b, 0, &addr); err != nil {
		panic(err)
	}

	// start reading incoming data
	for {
		buffer := make([]byte, 256)

		// var device block.Device

		n, from, err := syscall.Recvfrom(fd, buffer, 0)
		if err != nil {
			panic(err)
		}

		fmt.Println(n)
		fmt.Println(from)

		// fmt.Printf("% x\n", buffer[:n])

		e := frame.EthernetII{}
		if err := e.UnmarshalBinary(buffer); err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", e)

		// device.HardwareAddr = e.Source

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

		// fmt.Printf("%#v\n", device)

	}

}

func decodeBlock(b []byte) int {
	opt := option.Option(b[0])
	fmt.Println("option", opt)

	subopt := suboption.Suboption(b[1])
	fmt.Println("suboption", subopt)

	length := binary.BigEndian.Uint16(b[2:4])
	fmt.Println("length", length)

	switch {

	case opt == option.DeviceProperties && subopt == suboption.NameOfStation:

		var bnos block.NameOfStation
		if err := bnos.UnmarshalBinary(b); err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", bnos)
		fmt.Println(bnos.NameOfStation)

	case opt == option.IP && subopt == suboption.IPParameter:

		var bip block.IPParameter
		if err := bip.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bip)
		fmt.Println(bip.IPAddress, bip.Subnetmask, bip.StandardGateway)

	case opt == option.DeviceProperties && subopt == suboption.DeviceInstance:

		var bdi block.DeviceInstance
		if err := bdi.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bdi)
		fmt.Println(bdi.DeviceInstanceHigh, bdi.DeviceInstanceLow)

	case opt == option.DeviceProperties && subopt == suboption.ManufacturerSpecific:

		var bms block.ManufacturerSpecific
		if err := bms.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bms)
		fmt.Println(bms.DeviceVendorValue)

	case opt == option.DeviceInitiative && subopt == suboption.DeviceInitiative:

		var bdi block.DeviceInitiative
		if err := bdi.UnmarshalBinary(b); err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", bdi)
		fmt.Println(bdi.Value)
	}

	return 1 + 1 + 2 + int(length)
}
