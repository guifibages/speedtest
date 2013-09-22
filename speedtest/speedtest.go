package speedtest

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var sizes = map[string]int{
		"random350x350.jpg":   245388,
		"random500x500.jpg":   505544,
		"random750x750.jpg":   1118012,
		"random1000x1000.jpg": 1986284,
		"random1500x1500.jpg": 4468241,
		"random2000x2000.jpg": 7907740,
		"random2500x2500.jpg": 12407926,
		"random3000x3000.jpg": 17816816,
		"random3500x3500.jpg": 24262167,
		"24262167":            24262167,
		"random4000x4000.jpg": 31625365,
	}
	file := path.Base(r.URL.Path)
	if file == "upload.php" {
		fmt.Println("Upload")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading body", err)
		}
		fmt.Println("Upload:", string(body))
		return
	}
	size := sizes[string(file)]
	fmt.Printf("Serving %s (%d)B\n", file, size)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(size))
	randomSrc := randomDataMaker{rand.NewSource(1028890720402726901)}
	buf := make([]byte, 1024)
	for i := 0; i < size/1024; i++ {
		io.ReadAtLeast(&randomSrc, buf, 1024)
		w.Write(buf)
		if i%(1024) == 0 {
			fmt.Printf("Served %dKB\n", i)
		}
	}
}

func HTTPServer(dataPath string) {
	http.ListenAndServe(":59330", http.FileServer(http.Dir(dataPath)))
}

func Server() {
	port := ":59330"
	log.Printf("Listening on %s", port)
	http.HandleFunc("/speedtest/", handler)
	http.Handle("/", http.FileServer(http.Dir(".")))

	http.ListenAndServe(port, nil)
}

// http://stackoverflow.com/a/12810288
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
