package scrapper

import (
	"log"
	//"io"

	"github.com/PuerkitoBio/goquery"
)

func GetImageURLS(sku string) ([]string, bool) {
	log.Printf("Requesting \"/product_%v\" from sites", sku)
	res, successful := fetchBody(sku)

	if !successful {
		return nil, false
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var images []string

	// Find the review items
	doc.Find("#image_gallery_container").Each(func(i int, gallery *goquery.Selection) {
		// For each item found, get the title
		log.Println("Image Gallery Found for SKU")
		gallery.Find("img").Each(func(i int, img *goquery.Selection) {
			//parentEl := img.Parent()
			imgSource, _ := img.Attr("alt")
			images = append(images, imgSource)
		})
	})

	return images, true
}