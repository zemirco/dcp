package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/zemirco/dcp"
)

// host order (usually little endian) -> network order (big endian)
func htons(n int) int {
	return int(int16(byte(n))<<8 | int16(byte(n>>8)))
}

var (
	t    *template.Template
	db   = make(map[string]dcp.Frame)
	last time.Time
)

func init() {
	t = template.Must(template.ParseFiles("src/index.html"))
}

func main() {

	r := mux.NewRouter()

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	r.HandleFunc("/api/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(db); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/api/last", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(last); err != nil {
			panic(err)
		}
	})

	r.Methods(http.MethodGet).Path("/api/{mac}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		mac := vars["mac"]
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(db[mac]); err != nil {
			panic(err)
		}
	})

	r.Methods(http.MethodPost).Path("/api/{mac}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var f dcp.Frame
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
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

	request := dcp.NewIdentifyRequest(interf.HardwareAddr)
	b, err := request.MarshalBinary()
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

	last = time.Now()

	// start reading incoming data
	for {
		buffer := make([]byte, 256)

		n, _, err := syscall.Recvfrom(fd, buffer, 0)
		if err != nil {
			panic(err)
		}

		fmt.Println(n)

		response := dcp.Frame{}
		if err := response.UnmarshalBinary(buffer); err != nil {
			panic(err)
		}

		spew.Dump(response)

		// save device to db in case we have an answer to our identify request
		if request.XID == response.XID {
			db[response.Source.String()] = response
		}

	}

}
