package main

import (
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/draw"
	"os"
	"os/exec"
	"strconv"
)

func encode() {
	err := os.RemoveAll("./images/")
	err = os.Mkdir("./images/", os.ModePerm)

	file, err := os.ReadFile(FileName)
	if err != nil {
		panic(err)
	}

	imgFile, img := newFrame()

	sizeBoard := 15
	sizeBlock := PxWidth / sizeBoard

	for i := 0; i < len(file); i += 3 {
		var code1 byte = 0
		var code2 byte = 0
		var code3 byte = 0

		if !(i >= len(file)) {
			code1 = file[i]
		}

		if !(i+1 >= len(file)) {
			code2 = file[i+1]
		}

		if !(i+2 >= len(file)) {
			code3 = file[i+2]
		}

		Color := color.RGBA{R: code1, G: code2, B: code3, A: 255}

		draw.Draw(img, image.Rect(x, y, x+sizeBlock, y+sizeBlock),
			&image.Uniform{Color}, image.ZP, draw.Src)

		x += sizeBlock

		if i >= len(file) || i+1 >= len(file) || i+2 >= len(file) || i+3 >= len(file) {
			draw.Draw(img, image.Rect(x, y, x+sizeBlock, y+sizeBlock),
				&image.Uniform{color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.ZP, draw.Src)
		}

		if x == PxWidth {
			if y == PxHeight {
				x = 0
				y = 0

				err = bmp.Encode(imgFile, img)
				if err != nil {
					panic(err)
				}
				imgFile, img = newFrame()
			} else {
				x = 0
				y += sizeBlock
			}
		}
	}

	err = bmp.Encode(imgFile, img)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()

	toVideo()
}

func newFrame() (imgFile *os.File, img *image.RGBA) {
	numImage++

	upLeft := image.Point{}
	lowRight := image.Point{X: PxWidth, Y: PxHeight}

	img = image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	filename := "./images/img-" + strconv.Itoa(numImage) + ".bmp"

	imgFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	draw.Draw(img, img.Bounds(), img, image.Point{}, draw.Over)
	return
}

func toVideo() {
	cmd := exec.Command("ffmpeg", "-framerate", Framerate, "-i", "img-%d.bmp", "-crf", "0", "output.mp4")
	cmd.Dir = "images/"
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
