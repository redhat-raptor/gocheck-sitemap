package main

import (
	"net/http"
	"log"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/xml"
	"time"
	"encoding/json"
	"sync"
	"sync/atomic"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

const CheckInterval = 10

type URLs struct {
	Locs    []string    `xml:"url>loc"`
}

var URLStatuses = map[string]int{}

var httpAddr = fmt.Sprintf(":%s", getEnv("PORT", "3000"))


func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getSitemap(sitemapUrl string) []byte {
	resp, err := httpClient.Get(sitemapUrl)
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

func checkSitemap(c chan string) {
	for anUrl := range c {
		statusCode, err := getHTTPStatus(anUrl)
		if err != nil {
			continue
		}

		mu.Lock()
		URLStatuses[anUrl] = statusCode
		mu.Unlock()
		fmt.Println(anUrl, statusCode)
	}

	URLStatuses[anUrl] = statusCode
	fmt.Println(anUrl, statusCode)
}

func getHTTPStatus(anUrl string) (int, error) {
	resp, err := http.Get(anUrl)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func serveHTTPStatuses(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(URLStatuses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	siteMapUrls := os.Getenv("SITEMAP")

	if siteMapUrls == "" {
		log.Fatal("No sitemap url passed!")
		os.Exit(1)
	}

	sitemap := getSitemap(siteMapUrls)

	var urls URLs
	xml.Unmarshal(sitemap, &urls)

	//Listen to port to serve url statuses
	http.HandleFunc("/", serveHTTPStatuses)
	log.Println("Starting server ", httpAddr)

	for {
		left := urls.Locs
		for {
			var wg sync.WaitGroup
			wg.Add(20)
		
			for i := 0; i < 20; i++ {
				if i >= len(left) {
					break
				}
				go func(i int) {
					defer wg.Done()
					checkSitemap(left[i])
				}(i)
			}
			
			wg.Wait()
			time.Sleep(CheckInterval * time.Second)
			
			if len(left) < 20 {
				left = urls.Locs
			} else {
				left = left[20:]
			}
		}
	}

	for i := 0; i < 10; i++ {
		go checkSitemap(c)
	}

	log.Println(http.ListenAndServe(httpAddr, nil))
}
