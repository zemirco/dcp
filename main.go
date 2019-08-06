package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/rakyll/statik/fs"
	"github.com/zemirco/dcp/frame"
)

const etherType uint16 = 0x8892

// host order (usually little endian) -> network order (big endian)
func htons(n int) int {
	return int(int16(byte(n))<<8 | int16(byte(n>>8)))
}

var db = make(map[string]frame.Frame)

func main() {

	mode := flag.String("mode", "development", "switch between local file system and embedded one.")
	flag.Parse()

	// add file server
	if *mode == "production" {
		statikFS, err := fs.New()
		if err != nil {
			panic(err)
		}
		http.Handle("/", http.FileServer(statikFS))
	} else {
		fmt.Println("serving files from local file system")
		http.Handle("/", http.FileServer(http.Dir("ui/public")))
	}

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(db)
	})

	// start server
	go func() {
		fmt.Println("server running at http://localhost:8085. ctrl+c to stop it.")
		log.Fatal(http.ListenAndServe(":8085", nil))
	}()

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

	// // request block
	// rb := block.NewIPParameterQualifier()
	// rb.IPAddress = []byte{0xac, 0x13, 0x68, 0x03}
	// rb.Subnetmask = []byte{0xff, 0xff, 0x00, 0x00}
	// rb.StandardGateway = []byte{0x00, 0x00, 0x00, 0x00}

	// destination := []byte{0x00, 0x09, 0xe5, 0x00, 0x9a, 0x20}

	// req := frame.NewSetIPParameterRequest(destination, interf.HardwareAddr, rb)
	// b, err := req.MarshalBinary()
	// if err != nil {
	// 	panic(err)
	// }

	// spew.Dump(req)
	// spew.Dump(b)

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

		db[f.Source.String()] = f

	}

}
