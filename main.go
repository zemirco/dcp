package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"syscall"
)

type ethernetII struct {
	Destination [6]byte
	Source      [6]byte
	EtherType   uint16
}

type telegram struct {
	FrameID       frameID
	ServiceID     serviceID
	ServiceType   serviceType
	XID           uint32
	ResponseDelay uint16
}

type option byte

const (
	optionIP               option = 1
	optionDeviceProperties option = 2
	optionDeviceInitiative option = 6
	optionAll              option = 255
)

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

type frameID uint16

const (
	frameIDIdentifyRequest  frameID = 0xfefe
	frameIDIdentifyResponse frameID = 0xfeff
)

type serviceID byte

const (
	serviceIDGet      serviceID = 1
	serviceIDSet      serviceID = 2
	serviceIDIdentify serviceID = 5
)

type serviceType byte

const (
	serviceTypeRequest  serviceType = 0
	serviceTypeResponse serviceType = 1
)

type block struct {
	option    option
	suboption uint8
}

// host order (usually little endian) -> network order (big endian)
func htons(n int) int {
	return int(int16(byte(n))<<8 | int16(byte(n>>8)))
}

func main() {

	ifname := "enxa44cc8e54721"

	frame := make([]byte, 30)

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
	copy(frame[0:], buf.Bytes())

	t := telegram{
		FrameID:       frameIDIdentifyRequest,
		ServiceID:     serviceIDIdentify,
		ServiceType:   serviceTypeRequest,
		XID:           rand.Uint32(),
		ResponseDelay: 255,
	}

	buf.Reset()
	binary.Write(&buf, binary.BigEndian, t)
	copy(frame[14:], buf.Bytes())

	// dcp data length

	b := &block{
		option:    optionAll,
		suboption: 255,
	}

	buf.Reset()
	binary.Write(&buf, binary.BigEndian, b)

	// +2 because DCPBlockLength
	binary.BigEndian.PutUint16(frame[24:26], uint16(len(buf.Bytes()))+2)

	copy(frame[26:28], buf.Bytes())

	binary.BigEndian.PutUint16(frame[28:30], 0)

	log.Printf("% x\n", frame)

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, htons(0x8892))

	if err != nil {
		panic(err)
	}

	defer syscall.Close(fd)

	addr := syscall.SockaddrLinklayer{
		Ifindex: interf.Index,
	}

	if err := syscall.Sendto(fd, frame, 0, &addr); err != nil {
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
	option := option(b[0])
	fmt.Println("option", option)

	suboption := b[1]
	fmt.Println("suboption", suboption)

	length := binary.BigEndian.Uint16(b[2:4])
	fmt.Println("length", length)

	info := binary.BigEndian.Uint16(b[4:6])
	fmt.Println("info", info)

	switch {

	// device properties && name of station
	case option == optionDeviceProperties && suboption == 2:

		// info length is 2
		name := string(b[6 : 6+length-2])
		fmt.Println("name", name)

	// ip && ip parameter
	case option == optionIP && suboption == 2:

		ip := net.IP(b[6:10])
		fmt.Println("ip", ip.String())

		subnetmask := net.IP(b[10:14])
		fmt.Println("subnetmask", subnetmask.String())

		gateway := net.IP(b[14:18])
		fmt.Println("gateway", gateway.String())
	}

	return 1 + 1 + 2 + int(length)
}
