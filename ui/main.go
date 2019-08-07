package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/zemirco/dcp"
)

// host order (usually little endian) -> network order (big endian)
func htons(n int) int {
	return int(int16(byte(n))<<8 | int16(byte(n>>8)))
}

var db = make(map[string]dcp.Frame)

var (
	t *template.Template
)

func init() {
	t = template.Must(template.ParseFiles("src/index.html"))
}

func main() {

	r := mux.NewRouter()

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	r.HandleFunc("/api/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(db)
	})

	r.Methods(http.MethodGet).Path("/api/{mac}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		mac := vars["mac"]
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(db[mac])
	})

	r.Methods(http.MethodPost).Path("/api/{mac}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var f dcp.Frame
		err := json.NewDecoder(r.Body).Decode(&f)
		if err != nil {
			panic(err)
		}
		spew.Dump(f)
	})

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, nil)
	})

	// start server
	go func() {
		fmt.Println("server running at http://localhost:8085. ctrl+c to stop it.")
		log.Fatal(http.ListenAndServe(":8085", r))
	}()

	ifname := "enxa44cc8e54721"

	interf, err := net.InterfaceByName(ifname)
	if err != nil {
		panic(err)
	}

	f := dcp.NewIdentifyRequest(interf.HardwareAddr)
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

		f := dcp.Frame{}
		if err := f.UnmarshalBinary(buffer); err != nil {
			panic(err)
		}

		spew.Dump(f)

		db[f.Source.String()] = f

	}

}