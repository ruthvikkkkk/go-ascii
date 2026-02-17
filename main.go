package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/nfnt/resize"
)

func main() {

	outputPtr := flag.String("o", "", "Output File Name")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: go run main.go <image_file> [-o output_file]")
		os.Exit(1)
	}

	filename := args[0]

	asciiChars := ".:-=+*#%@"
	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	newWidth := uint(40)

	ratio := float64(img.Bounds().Dy()) / float64(img.Bounds().Dx())
	newHeight := uint(float64(newWidth) * ratio * 0.5)

	resizedImage := resize.Resize(newWidth, newHeight, img, resize.Lanczos2)
	bounds := resizedImage.Bounds()

	w, h := bounds.Max.X, bounds.Max.Y

	var writer io.Writer = os.Stdout

	if *outputPtr != "" {
		f, err := os.Create(*outputPtr)
		if err != nil {
			fmt.Println("Error creating output file: ", err)
			return
		}
		defer f.Close()
		writer = f
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := resizedImage.At(x, y)

			gray := color.GrayModel.Convert(c).(color.Gray)

			i := int(gray.Y) * (len(asciiChars) - 1) / 255

			fmt.Fprintf(writer, "%c", asciiChars[i])
		}
		fmt.Fprintln(writer)
	}

}
