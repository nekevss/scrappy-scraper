package cmd

import (
	"fmt"
	"log"
	"os"
	"io"
	"time"
	"encoding/csv"
	"strings"
	"image"
	"image/jpeg"
	"image/draw"
	"image_scrapper/core"
	"path/filepath"
	"github.com/spf13/cobra"
)

var inputPath string
var compareImagePath string
var outputPath string

var rootCmd = &cobra.Command{
	Use: "image_scrapper",
	Short: "Utility tool for checking broken images",
	Run: func(cmd *cobra.Command, args[]string) {
		// Start time for the operation is taken here
		start := time.Now()

		fileName := filepath.Base(inputPath)
		fileExtension := filepath.Ext(inputPath)
		fileName = strings.ReplaceAll(fileName, fileExtension, "")


		logFile, err := os.OpenFile(fileName + "_run_info.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666) // Setting the bits to read/write for everything
		if err != nil {
			log.Fatal("Logging the failure to initiate log file...sort of ironic XD")
		}

		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)

		// Load core image
		compareRaw, err := os.Open(compareImagePath)
		if err != nil {
			log.Print("Was unable to load the image for comparing\n")
			log.Fatal(err)
		}
		defer compareRaw.Close()

		jpegImage, err := jpeg.Decode(compareRaw)
		if err != nil {
			log.Print("Error decoding provided jpg")
			log.Fatal(err)
		}

		bound := jpegImage.Bounds()
		compareImage := image.NewRGBA(image.Rect(0, 0, bound.Dx(), bound.Dy()))
		draw.Draw(compareImage, compareImage.Bounds(), jpegImage, bound.Min, draw.Src)

		log.Println("Beginnning work on submitted file", inputPath)
		// Load the inputed CSV
		inputCsv, err := os.Open(inputPath)
		if err != nil {
			log.Print("Unable to read inputted CSV file\n")
			log.Fatal(err)
		}
		defer inputCsv.Close()

		csvReader := csv.NewReader(inputCsv)

		inputItems, err := csvReader.ReadAll()
		if err != nil {
			log.Print("Unable to read CSV file into [][]string\n")
			log.Fatal(err)
		}

		outputFilepath := filepath.Join(outputPath,fileName + "_output.csv")
		outputFile, err := os.Create(outputFilepath)
		if err != nil {
			log.Print("Unable to create output csv file\n")
			log.Fatal(err)
		}
		
		outputRecords := core.ScrapeAndCompare(inputItems, compareImage)

		writer := csv.NewWriter(outputFile)
		err = writer.WriteAll(outputRecords)
		if err != nil {
			log.Print("Unable to write\n")
			log.Fatal(err)
		}

		log.Println("Completed! :)")
		log.Println("Run time was ", time.Since(start))

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Input file with items to check")
	rootCmd.Flags().StringVarP(&compareImagePath, "output", "o", "", "Output Directory")
	rootCmd.Flags().StringVarP(&compareImagePath, "compare-input", "c", "", "Image to compare against (jpg encoded)")
}