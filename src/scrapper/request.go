package scrapper

import (
	"log"
	"net/http"
)

// Want to note that this is the dumbest approach...but 
// this thing just needs to work
func fetchBody(sku string) (*http.Response, bool) {
	// Check dotcom
	composedUrl := "https://www.staples.com/product_" + sku
	comRes, err := http.Get(composedUrl)
	if err != nil {
		log.Fatal(err)
	}

	if comRes.StatusCode == 200 {
		return comRes, true
	}

	// Check SA
	composedUrl = "https://www.staplesadvantage.com/product_" + sku

	saRes, err := http.Get(composedUrl)
	if err != nil {
		log.Fatal(err)
	}

	if saRes.StatusCode == 200 {
		return saRes, true
	}

	return saRes, false
} 