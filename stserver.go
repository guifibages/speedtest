package main

import (
	"flag"
	"github.com/guifibages/speedtest/speedtest"
)

var dataPath string

func main() {
	const (
		defaultPath = "/guifibages/var/speedtest"
		usage       = "ruta dels arxius a servir per speedtest"
	)
	//	speedtest.Server()
	flag.StringVar(&dataPath, "data_path", defaultPath, usage)
	flag.StringVar(&dataPath, "d", defaultPath, usage+" (shorthand)")
	flag.Parse()

	speedtest.HTTPServer(dataPath)
}
