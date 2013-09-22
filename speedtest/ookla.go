package speedtest

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	OoklaClientConfigURL = "http://www.speedtest.net/speedtest-config.php"
)

//	PrefServer      = "http://speedtest.bcn.adamo.es/speedtest"

type OoklaClient struct {
	XMLName xml.Name          `xml:"settings"`
	Client  OoklaClientConfig `xml:"client"`
	Server  string
	Timeout int
}

type OoklaClientConfig struct {
	IP  string `xml:"ip,attr"`
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
}

func (c *OoklaClient) GetConfig() error {
	res, err := http.Get(OoklaClientConfigURL)
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

func (c *OoklaClient) Download() []*Result {
	sizes := [10]int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
	result := make([]*Result, 10)
	for i, v := range sizes {
		file := fmt.Sprintf("random%dx%d.jpg", v, v)
		url := fmt.Sprintf("http://%s/speedtest/%s", c.Server, file)
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
		fmt.Printf("Downloading %.2fB file (%s) (Latency: %.2fms) from %s\n", size, file, latency, c.Server)
		downTimer := make(chan int)
		go func(res *http.Response) {
			res.Write(ioutil.Discard)
			downTimer <- 1
		}(res)
		select {
		case _ = <-downTimer:
			lapse := time.Since(start).Seconds()
			result[i] = NewResult(size, lapse)
			result[i].Latency = latency
			result[i].File = url
			//fmt.Printf("\tURL:%s (%s) %dbytes in %fseconds (%fbps)\n", url, res.Status, size, lapse, speed)
		case <-time.After(time.Duration(c.Timeout) * time.Second):
			fmt.Printf("Timed out on %.2fMB file\n", size/1024/1024)
			return result // Timed out
		}
		//			fmt.Printf("\tURL:%s (%s) %s\n\t%v\n", url, res.Status, size, res.Header)
	}
	return result
}

func (c *OoklaClient) Upload() {

	buf := make([]byte, (1024))
	randomSrc := randomDataMaker{rand.NewSource(1028890720402726901)}
	io.ReadFull(&randomSrc, buf)
	data := string(buf)
	extension := "php"
	uploadurl := fmt.Sprintf("http://%s/speedtest/upload.%s", c.Server, extension)
	fmt.Println("Uploading to", uploadurl, "data", len(data))
	v := url.Values{}
	v.Add("content1", data)
	resp, err := http.PostForm(uploadurl,
		v)
	if err != nil {
		fmt.Println("Error posting", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body", err)
	}
	fmt.Println("Response:", string(body))
}

func (c *OoklaClient) TestServer() {
	c.Upload()
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
