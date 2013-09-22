package main

import (
	//	"flag"
	"flag"
	"fmt"
	"github.com/guifibages/speedtest/speedtest"
	"log"
)

/*
* Perpetrated in 2013 by http://ignacio.torresmasdeu.name/
*
* Clever ideas from https://github.com/sivel/speedtest-cli/
* All the bugs are mine.
 */

/* <client ip="83.50.183.140" lat="41.5637" lon="2.0167"
* isp="Telefonica de Espana" isprating="2.5" rating="0"
* ispdlavg="7018" ispulavg="894" loggedin="0" />

 */

func main() {

	var listen bool
	flag.BoolVar(&listen, "listen", false, "act as a server")
	flag.BoolVar(&listen, "l", false, "act as a server (shorthand)")

	c := new(speedtest.OoklaClient)
	flag.StringVar(&c.Server, "server", "speedtest.guifibages.net", "server to test against")
	flag.StringVar(&c.Server, "s", "localhost:12345", "server to test against (shorthand)")
	flag.IntVar(&c.Timeout, "timeout", 30, "timeout in seconds")
	flag.IntVar(&c.Timeout, "t", 30, "timeout in seconds (shorthand)")
	flag.Parse()
	if listen {
		speedtest.Server()
	} else {
		err := c.GetConfig()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Speedtest against %s with a %d seconds timeout\n", c.Server, c.Timeout)
		//fmt.Printf("IP: %s\nLon:%s\nLat:%s\n", c.Client.IP, c.Client.Lon, c.Client.Lat)
		c.TestServer()
	}
}
