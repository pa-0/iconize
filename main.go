package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"time"

	"git.shangtai.net/staffan/go-ico"
	"golang.org/x/image/draw"
)

func main() {

	start := time.Now()

	if len(os.Args) < 2 {
		os.Exit(1)
	}
	path := os.Args[1]

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		log.Fatal(err)
	}

	width, height := size(img)
	size := max(width, height)

	alpha := image.NewRGBA(image.Rect(0, 0, size, size))

	rgba := color.RGBA{0, 0, 0, 0}
	draw.Draw(alpha, alpha.Bounds(), &image.Uniform{C: rgba}, image.Point{}, draw.Src)

	draw.Draw(alpha, img.Bounds().Add(
		image.Point{
			X: (size - img.Bounds().Dx()) / 2,
			Y: (size - img.Bounds().Dy()) / 2,
		}), img, image.Point{}, draw.Over)

	basename := name(path)
	icoFile, err := os.Create(basename + ".ico")
	if err != nil {
		log.Fatal(err)
	}
	defer icoFile.Close()

	icon := ico.NewIcon()
	// https://learn.microsoft.com/en-us/windows/win32/uxguide/vis-icons
	sizes := []int{128, 256}
	for _, size := range sizes {
		resizedImg := scale(alpha, size, size)
		icon.AddPng(resizedImg)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Elapsed:", elapsed)

	enc, err := icon.Encode()
	if err != nil {
		log.Fatal(err)
	}

	icoFile.Write(enc)
	os.Exit(0)

}

func max(width, height int) int {
	size := width
	if height > width {
		size = height
	}
	return size
}

func size(img image.Image) (int, int) {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	return width, height
}

func scale(img image.Image, width, height int) image.Image {
	resizedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	return resizedImg
}

func name(filePath string) string {
	fileName := filepath.Base(filePath)
	extension := filepath.Ext(fileName)
	return fileName[0 : len(fileName)-len(extension)]
}