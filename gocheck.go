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
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

const CheckInterval = 2

type URLs struct {
	Locs    []string    `xml:"url>loc"`
}

type HTTPStatus struct {
	Status int
	Error error
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

func checkSitemap(done chan bool, urls URLs) {
	httpStatusChannel := make(chan HTTPStatus)

	for _, anUrl := range urls.Locs {
		time.Sleep(CheckInterval * time.Second)

		getHTTPStatus(httpStatusChannel, anUrl)
		httpStatus := <- httpStatusChannel

		if httpStatus.Error != nil {
			log.Println("Error checking http status", httpStatus.Error)
			continue
		}

		URLStatuses[anUrl] = httpStatus.Status
		fmt.Println(anUrl, httpStatus.Status)
	}

	done <- true
}

func getHTTPStatus(ch chan HTTPStatus, anUrl string) {
	resp, err := http.Get(anUrl)

	ch <- HTTPStatus{
		Error: err,
		Status: resp.StatusCode,
	}
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

func keepChecking(urls URLs) {
	done := make(chan bool)

	for {
		checkSitemap(done, urls)
		completed := <- done

		if completed == true {
			os.Exit(0)
		}
	}
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

	go keepChecking(urls)

	log.Println(http.ListenAndServe(httpAddr, nil))
}