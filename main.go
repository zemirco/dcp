package main

import (
	"fmt"
	"net"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/zemirco/dcp/frame"
)

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

		n, _, err := syscall.Recvfrom(fd, buffer, 0)
		if err != nil {
			panic(err)
		}

		fmt.Println(n)

		f := frame.Frame{}
		if err := f.UnmarshalBinary(buffer); err != nil {
			panic(err)
		}

		spew.Dump(f)

	}

}
