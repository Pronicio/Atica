package main

import (
	"crypto/rand"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/draw"
	"math/big"
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

	sizeBoard := 30
	sizeBlock := PxWidth / sizeBoard

	for i := 0; i < len(file); i++ {
		code := file[i]
		//print(code, " ")

		nBig, err := rand.Int(rand.Reader, big.NewInt(255))
		random1 := uint8(nBig.Int64())

		nBig, err = rand.Int(rand.Reader, big.NewInt(255))
		random2 := uint8(nBig.Int64())

		if err != nil {
			panic(err)
		}

		Color := color.RGBA{R: code, G: random1, B: random2, A: 255}

		draw.Draw(img, image.Rect(x, y, x+sizeBlock, y+sizeBlock),
			&image.Uniform{Color}, image.ZP, draw.Src)

		x += sizeBlock

		if len(file) == i+1 {
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
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	println(out)
}
