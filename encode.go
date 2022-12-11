package main

import (
	"github.com/sergeymakinen/go-bmp"
	"image"
	"image/color"
	"image/draw"
	"os"
	"os/exec"
	"strconv"
)

var (
	Black = color.Gray16{Y: 0}
	White = color.Gray16{Y: 0xffff}
)

func encode() {
	err := os.RemoveAll("./images/")
	err = os.Mkdir("./images/", os.ModePerm)

	file, err := os.ReadFile(FileName)
	if err != nil {
		panic(err)
	}

	imgFile, img := newFrame()

	sizeBoard := 160
	sizeBlock := PxWidth / sizeBoard

	for i := 0; i < len(file); i++ {
		bin := decimalToBinary(file[i])

		for _, letter := range bin {
			binNumber := string(letter)

			if binNumber == "1" {
				draw.Draw(img, image.Rect(x, y, x+sizeBlock, y+sizeBlock),
					&image.Uniform{White}, image.ZP, draw.Src)
			} else {
				draw.Draw(img, image.Rect(x, y, x+sizeBlock, y+sizeBlock),
					&image.Uniform{Black}, image.ZP, draw.Src)
			}

			x += sizeBlock

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

		// Last loop :
		if i == (len(file) - 1) {
			for j := 0; j < 8; j++ {
				draw.Draw(img, image.Rect(x, y, x+sizeBlock, y+sizeBlock),
					&image.Uniform{White}, image.ZP, draw.Src)

				x += sizeBlock

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
		}
	}

	err = bmp.Encode(imgFile, img)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()

	toVideo()
}

func newFrame() (imgFile *os.File, img *image.Paletted) {
	numImage++

	palette := color.Palette([]color.Color{
		color.Gray16{Y: 0},
		color.Gray16{Y: 0xffff},
	})

	upLeft := image.Point{}
	lowRight := image.Point{X: PxWidth, Y: PxHeight}

	img = image.NewPaletted(image.Rectangle{Min: upLeft, Max: lowRight}, palette)
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

func decimalToBinary(decimal byte) string {
	return strconv.FormatInt(int64(decimal), 2)
}
