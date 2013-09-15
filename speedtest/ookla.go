package speedtest

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	ClientConfigURL = "http://www.speedtest.net/speedtest-config.php"
)

//	PrefServer      = "http://speedtest.bcn.adamo.es/speedtest"

type OoklaClient struct {
	XMLName xml.Name          `xml:"settings"`
	Client  OoklaClientConfig `xml:"client"`
	Server  string
}

type OoklaClientConfig struct {
	IP  string `xml:"ip,attr"`
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
}

func (c *OoklaClient) GetConfig() error {
	res, err := http.Get(ClientConfigURL)
	if err != nil {
		return err
	}
	configxml, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}

	err = xml.Unmarshal(configxml, &c)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	return nil
}

type DownloadResult struct {
	Speed   float64
	Seconds float64
	Size    float64
	File    string
	Latency float64
}

func (c *OoklaClient) Download() []DownloadResult {
	sizes := [10]int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
	result := make([]DownloadResult, 10)
	for i, v := range sizes {
		url := fmt.Sprintf("http://%s/speedtest/random%dx%d.jpg", c.Server, v, v)
		start := time.Now()
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		size, err := strconv.ParseFloat(res.Header["Content-Length"][0], 64)
		if err != nil {
			log.Fatal(err)
		}
		latency := time.Since(start).Seconds() * 1000
		fmt.Printf("Downloading %.2fMB file (Latency: %.2fms) from %s\n", size/1024/1024, latency, c.Server)
		downTimer := make(chan int)
		go func(res *http.Response) {
			res.Write(ioutil.Discard)
			downTimer <- 1
		}(res)
		select {
		case _ = <-downTimer:
			lapse := time.Since(start).Seconds()
			speed := size * 8 / lapse
			result[i] = DownloadResult{Speed: speed, Seconds: lapse, Size: size, File: url, Latency: latency}
			//fmt.Printf("\tURL:%s (%s) %dbytes in %fseconds (%fbps)\n", url, res.Status, size, lapse, speed)
		case <-time.After(20 * time.Second):
			fmt.Printf("Timed out on %.2fMB file\n", size/1024/1024)
			return result // Timed out
		}
		//			fmt.Printf("\tURL:%s (%s) %s\n\t%v\n", url, res.Status, size, res.Header)
	}
	return result
}

func (c *OoklaClient) TestServer(server string) {
	c.Server = server
	res := c.Download()
	var j int
	var totaltime, totalbytes, maxspeed, minspeed, maxlat, minlat, totallat float64
	for i, v := range res {
		j = i
		if v.File == "" {
			break
		}
		if minspeed == 0.0 || v.Speed < minspeed {
			minspeed = v.Speed
		}
		if minlat == 0.0 || v.Latency < minlat {
			minlat = v.Latency
		}
		if v.Speed > maxspeed {
			maxspeed = v.Speed
		}
		if v.Latency > maxlat {
			maxlat = v.Latency
		}

		totallat += v.Latency
		totalbytes += v.Size
		totaltime += v.Seconds
	}
	j += 1
	medlat := totallat / float64(j)
	medspeed := (totalbytes * 8) / totaltime

	fmt.Printf("Summary: %.2fMB transferred in %.2f seconds from %s\n", totalbytes/1024/1024, totaltime, c.Server)
	fmt.Printf("  Latency: min: %.2fms med: %.2fms max: %.2fms\n", minlat, medlat, maxlat)
	fmt.Printf("  Speed: min: %.2fmbps med: %.2fmbps max: %.2fmbps\n", minspeed/1000000, medspeed/1000000, maxspeed/1000000)
}
