package compare

import (
	"image"
	"log"
	"strings"
	"strconv"
	// "github.com/vitali-fedulov/images4"
)


func FetchandCompareImages(images []string, determinedBroken []string, compareImage *image.RGBA) ([]string, []string) {
	foundBroken := []string{}
	brokenURLs := []string{}
	for i, srcImageURL := range images {
		imageURL := srcImageURL
		if !strings.Contains(imageURL, "https://") {
			log.Println("Image does not appear to have a valid URL ", imageURL)
			brokenMessage := "Image at sequence " + strconv.Itoa(i) + " is broken: " + srcImageURL
			foundBroken = append(foundBroken, brokenMessage)
			brokenURLs = append(brokenURLs, srcImageURL)
			continue
		}

		alreadyChecked := Contains(determinedBroken, srcImageURL)
		if alreadyChecked {
			log.Println("Already found to be broken for URL", srcImageURL)
			brokenMessage := "Image at sequence " + strconv.Itoa(i) + " is broken: " + srcImageURL
			foundBroken = append(foundBroken, brokenMessage)
			brokenURLs = append(brokenURLs, srcImageURL)
		} else {
			// Normalize the URL
			// We have to adjust the URL so that the requests are all similar and from a single server
			// dotcom: https://www.staples-3p.com/s7/is/image/Staples/43590_sc7
			// staplesadvantage: https://images.staplesadvantage.com/is/image/Staples/43590_sc7
			if strings.Contains(imageURL, "images.staplesadvantage.com") {
				log.Println("images.staplesadvantage.com URL found. URL will be converted to staples-3p")
			}
			imageURL = strings.ReplaceAll(imageURL, "images.staplesadvantage.com", "www.staples-3p.com/s7")

			// Run fetch and then diff
			fetchedImage, requestErrorMessage := GetImageFromURL(imageURL)
			if fetchedImage != nil {
				diff := DiffImages(fetchedImage, compareImage)
				log.Println("Is fetched image different from the provided? ", diff)
				if !diff {
					brokenMessage := "Image at sequence " + strconv.Itoa(i) + " is broken: " + srcImageURL
					foundBroken = append(foundBroken, brokenMessage)
					brokenURLs = append(brokenURLs, srcImageURL)
				}
			} else {
				// A note: for fetchedimage to be nil, request error message must be populated based off returns
				statusCodeMessage := srcImageURL + requestErrorMessage
				foundBroken = append(foundBroken, statusCodeMessage)
				brokenURLs = append(brokenURLs, srcImageURL)
			}
		}
	}

	return foundBroken, brokenURLs
}

// External Source
//func DiffImagesExternal(fetched, provided image.Image) bool {
//	fetchedIcon := images4.Icon(fetched)
///	providedIcon := images4.Icon(provided)
//	return images4.Similar(fetchedIcon, providedIcon)
//}


// My attempt
func DiffImages(fetched, provided *image.RGBA) bool {
	if fetched.Bounds() != provided.Bounds() {
		log.Printf("Bounds don't match between fetched %v and %v\n", fetched.Bounds(), provided.Bounds())
		return true
	}

	// Adapted from https://stackoverflow.com/questions/32680834/how-to-compare-images-with-go

	imageDiff := int64(0)

	for i:=0; i< len(provided.Pix); i++ {
		imageDiff += int64(sqDiffUInt8(provided.Pix[i], fetched.Pix[i]))
	}

	log.Printf("Image diff value is: %v\n", imageDiff)

	// Due to encoding and HTTP transfer there could be differences in the same 
	// image, so below is an arbitrary number that could throw some false positives, 
	// but in theory should still catch everything 
	// 
	// Note: number should probably never go below 50k
	if imageDiff > 100000 {
		return true
	}

	return false
}

func sqDiffUInt8(x, y uint8) uint64 {   
    d := uint64(x) - uint64(y)
    return d * d
}

func Contains[T comparable](s []T, e T) bool {
    for _, v := range s {
        if v == e {
            return true
        }
    }
    return false
}