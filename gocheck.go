package main

import (
	"net/http"
	"log"
	"os"
	"fmt"
	"time"
	"encoding/json"

	"./urlstatus"
	"./sitemap"
)

const CheckInterval = 2

var httpAddr = fmt.Sprintf(":%s", getEnv("PORT", "3000"))


func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func checkSitemap(urls sitemap.URLs) {
	for _, anUrl := range urls.Locs {
		time.Sleep(CheckInterval * time.Second)

		statusCode, err := getHTTPStatus(anUrl)
		if err != nil {
			continue
		}

		urlstatus.SetUrlHTTPStatus(anUrl, statusCode)
		fmt.Println(anUrl, statusCode)
	}
}

func getHTTPStatus(anUrl string) (int, error) {
	resp, err := http.Get(anUrl)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func serveHTTPStatuses(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(urlstatus.GetUrlStatuses())
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

	sitemap.InitSitemap(siteMapUrls)

	//Listen to port to serve url statuses
	http.HandleFunc("/", serveHTTPStatuses)
	log.Println("Starting server ", httpAddr)

	go func() {
		for {
			checkSitemap(sitemap.GetSitemapURLs())
		}
	}()

	log.Println(http.ListenAndServe(httpAddr, nil))
}
