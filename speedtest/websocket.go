package speedtest

import (
	"code.google.com/p/go.net/websocket"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func Generate() []byte {
	b := make([]byte, 1024)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		fmt.Println("error:", err)
	}
	return b
}

// Echo the data received on the WebSocket.
func SpeedTest(ws *websocket.Conn) {
	var t Test
	var err error
	var r Result
	b := Generate()
	fmt.Println("Server: Receiving JSON")
	err = websocket.JSON.Receive(ws, &t)
	if err != nil {
		fmt.Println("JSON Error:", t)
		panic("Receiving JSON: " + err.Error())

	}
	fmt.Println("Server: Received", t)
	start := time.Now()

	fmt.Printf("Sending: %.0fKB.\n", t.Size)
	for i := 0.0; i < t.Size; i++ {
		websocket.Message.Send(ws, b)
	}
	fmt.Printf("Sent: %.0fKB in %.3f seconds.\n", t.Size, time.Since(start).Seconds())
	err = websocket.JSON.Receive(ws, &r)
	if err != nil {
		panic("Receiving Result: " + err.Error())

	}
	fmt.Printf("Server Result: Size: %.0fKB transfered in %.3f seconds. Latency: %fms\n", r.Size, r.Seconds, r.Latency)

}

func SpeedTestServer() {
	//http.Handle("/control", websocket.Handler(Control))
	http.Handle("/send", websocket.Handler(SpeedTest))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func ClientTest(address string) {
	fmt.Printf("Starting client test\n")

	origin := "http://" + address + "/"
	url := "ws://" + address + "/send"
	start := time.Now()
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	latency := time.Since(start).Seconds() * 1000
	t := Test{Size: 100, Direction: "both"}
	fmt.Println("Sending JSON")
	websocket.JSON.Send(ws, &t)
	var data []byte
	fmt.Println("Receiving Message")
	for i := 0.0; i < t.Size; i++ {
		websocket.Message.Receive(ws, &data)
	}
	duration := time.Since(start).Seconds()
	fmt.Printf("Received: %dMB in %.3f seconds.\n", t.Size, duration)
	result := Result{Size: t.Size, Seconds: duration, Latency: latency, Sender: address}
	websocket.JSON.Send(ws, result)
	fmt.Printf("Client Result: Size: %dMB Latency: %.3fms Time: %.3fs\n", result.Size, result.Latency, result.Seconds)
}
