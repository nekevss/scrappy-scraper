package compare

import (
	"net/http"
	"log"
	"image"
	"image/jpeg"
	"image/png"
	//"golang.org/x/image/webp"
	"image/draw"
)

// Add an auto fail to this function for unsupported types
func GetImageFromURL(imageURL string) (*image.RGBA, string) {
	res, err := http.Get(imageURL)
	if err != nil {
		log.Print("Image URL Request failed\n")
		return nil, " did not fetch properly from CDN."
	}

	if res.StatusCode != 200 {
		log.Printf("%s returned status code: %d", imageURL, res.StatusCode)
		log.Print("Cannot compare ", imageURL)
		return nil, " did not return a 200 status code."
	}

	var fetchedImage image.Image

	switch res.Header.Get("Content-Type") {
	case "image/jpeg":
		log.Println("Decoding jpeg image.")
		fetchedImage, err = jpeg.Decode(res.Body)
		if err != nil {
			log.Println("Error decoding jpeg mime type:", err)
			log.Println("Auto failing image")
			return nil, " jpg did not decode successfully."
		}
	case "image/png":
		log.Println("Decoding png image.")
		fetchedImage, err = png.Decode(res.Body)
		if err != nil {
			log.Println("Error decode png mime type:", err)
			log.Println("Auto failing image")
			return nil, " png did not decode successfully."
		}
	default:
		log.Println("URL Response Mime-type is not supported")
		log.Println("Failing image for manual check")
		return nil, " unsupported image type"
	}

	// Convert the image.Image into a image.RGBA
	b := fetchedImage.Bounds()
	im := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(im, im.Bounds(), fetchedImage, b.Min, draw.Src)

	// Below is a horrid idea, but nil in general is a bad idea, so too bad haters
	return im, "stringed-nil"
}

