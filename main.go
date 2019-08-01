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
	"github.com/zemirco/dcp/option"
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

type myBlock struct {
	option    option.Option
	suboption uint8
}

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

	b := &myBlock{
		option:    option.All,
		suboption: 255,
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

		length := binary.BigEndian.Uint16(buffer[24:26])
		fmt.Println("length", length)

		// fmt.Println(int(dcpDataLength))

		blockLength := decodeBlock(buffer[26:])

		if int(length)-blockLength > 0 {
			blockLength = decodeBlock(buffer[26+blockLength:])

		}

	}

}

func decodeBlock(b []byte) int {
	optSubopt := block.OptionSuboption(binary.BigEndian.Uint16(b[0:2]))
	fmt.Println("option suboption", optSubopt)

	length := binary.BigEndian.Uint16(b[2:4])
	fmt.Println("length", length)

	info := binary.BigEndian.Uint16(b[4:6])
	fmt.Println("info", info)

	switch optSubopt {

	case block.DevicePropertiesNameOfStation:

		// info length is 2
		name := string(b[6 : 6+length-2])
		fmt.Println("name", name)

	case block.IPIPParameter:

		ip := net.IP(b[6:10])
		fmt.Println("ip", ip.String())

		subnetmask := net.IP(b[10:14])
		fmt.Println("subnetmask", subnetmask.String())

		gateway := net.IP(b[14:18])
		fmt.Println("gateway", gateway.String())
	}

	return 1 + 1 + 2 + int(length)
}
