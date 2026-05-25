package ocr

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func preprocessImageFile(sourcePath string, threshold uint8) (string, error) {
	source, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer source.Close()

	img, _, err := image.Decode(source)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	output := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			if gray.Y >= threshold {
				output.Set(x, y, color.Gray{Y: 255})
			} else {
				output.Set(x, y, color.Gray{Y: 0})
			}
		}
	}

	tempFile, err := os.CreateTemp("", "ocr_preprocessed_*.png")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	if err := png.Encode(tempFile, output); err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}
