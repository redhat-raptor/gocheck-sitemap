package main

import (
	"net/http"
	"log"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/xml"
	"time"
)

const CheckInterval = 2

type URLs struct {
	Locs    []string    `xml:"url>loc"`
	httpStatus	int8
}

func getSitemap(sitemapUrl string) []byte {
	resp, err := http.Get(sitemapUrl)
	if err != nil {
		log.Fatal("Error: ", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return bodyBytes
}

func checkSiteamp(urls URLs) {
	for _, anUrl := range urls.Locs {
		time.Sleep(CheckInterval * time.Second)
		fmt.Println(anUrl)
	}

	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No sitemap url passed!")
		os.Exit(0)
	}

	siteMapUrls := os.Args[1]
	sitemap := getSitemap(siteMapUrls)

	var urls URLs
	xml.Unmarshal(sitemap, &urls)

	for {
		checkSiteamp(urls)
	}
}