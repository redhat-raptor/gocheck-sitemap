package urlstatus

import (
	"sync"
)

var mu sync.Mutex
var urlStatuses = map[string]int{}

func SetUrlHTTPStatus(url string, httpCode int) {
	mu.Lock()
	urlStatuses[url] = httpCode
	mu.Unlock()
}

func GetUrlStatuses() map[string]int {
	mu.Lock()
	defer mu.Unlock()
	return urlStatuses
}