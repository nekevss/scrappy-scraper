package core

import (
	"image_scrapper/scrapper"
	"image_scrapper/compare"
	"image"
	"log"
)


func ScrapeAndCompare(items [][]string, compareImage *image.RGBA) [][]string {
	// iterate through items
	outputRecords := [][]string {
		{"SKU", "Broken 1", "Broken 2", "Broken 3", "Broken 4", "Broken 5", "Broken 6", "Broken 7", "Broken 8", "Broken 9"},
	}
	
	determinedBroken := []string{}

	for _, row := range items {
		log.Println("Beginning process for item:", row[0])
		images, successful:= scrapper.GetImageURLS(row[0])

		if !successful {
			log.Println("Was not successful retrieving product from site")
			message:= "Check manually. There was a server error or the item was not actually live."
			newRow:=[]string{row[0]}
			newRow = append(newRow, message)
			outputRecords = append(outputRecords, newRow)
			continue;
		}

		log.Printf("Scrapped Images: %+q\n", images)
		potentialRow, brokenURLs := compare.FetchandCompareImages(images, determinedBroken, compareImage)
		
		for _, broken := range brokenURLs {
			if broken != "" {
				determinedBroken = append(determinedBroken, broken)
			}
		}

		if len(potentialRow) != 0 {
			log.Printf("Row provided to add: %+q\n", potentialRow)
			newRow := []string{row[0]}
			newRow = append(newRow, potentialRow...)
			outputRecords = append(outputRecords, newRow)
		}
	}
	
	return outputRecords
}