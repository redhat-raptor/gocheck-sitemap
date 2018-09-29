package sitemap

import (
	"log"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/xml"
	"sync"
)

var mu sync.Mutex

type URLs struct {
	Locs    []string    `xml:"url>loc"`
}

var urls URLs

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func InitSitemap(sitemapUrl string) {
	resp, err := httpClient.Get(sitemapUrl)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	mu.Lock()
	defer mu.Unlock()
	xml.Unmarshal(bodyBytes, &urls)
}

func GetSitemapURLs() URLs {
	mu.Lock()
	defer mu.Unlock()
	return urls
}