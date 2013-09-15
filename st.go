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
	var server_address string
	usage := "server to test"
	flag.StringVar(&server_address, "server", "speedtest.guifibages.net", usage)
	flag.StringVar(&server_address, "s", "speedtest.guifibages.net", usage+" (shorthand)")
	flag.Parse()

	c := new(speedtest.OoklaClient)
	err := c.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Speedtest against", server_address)
	fmt.Printf("IP: %s\nLon:%s\nLat:%s\n", c.Client.IP, c.Client.Lon, c.Client.Lat)
	c.TestServer(server_address)
}
