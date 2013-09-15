package speedtest

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func HTTPServer(dataPath string) {
	port := 59330
	log.Printf("Listening on %d", port)
	http.ListenAndServe(":59330", http.FileServer(http.Dir(dataPath)))
}

func Server() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func EchoFunc(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		println("Error reading:", err.Error())
		return
	}

	println("In EchoFunc")
	//send reply
	megas, err := strconv.Atoi(string(buf[0]))
	if err != nil {
		println("Error reading:", err.Error())
		return
	}
	kilobytes := 1024 * megas

	for i := 0; i < kilobytes; i++ {
		randomSrc := randomDataMaker{rand.NewSource(1028890720402726901)}
		io.ReadAtLeast(&randomSrc, buf, 1024)
		_, err = conn.Write(buf)
		if err != nil {
			println("Error send reply:", err.Error())
		}
	}
	conn.Close()
}

type randomDataMaker struct {
	src rand.Source
}

func (r *randomDataMaker) Read(p []byte) (n int, err error) {
	todo := len(p)
	offset := 0
	for {
		val := int64(r.src.Int63())
		for i := 0; i < 8; i++ {
			p[offset] = byte(val)
			todo--
			if todo == 0 {
				return len(p), nil
			}
			offset++
			val >>= 8
		}
	}

	panic("unreachable")
}
