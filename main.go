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

const (
	FileName  = "text.txt"
	PxWidth   = 1920
	PxHeight  = 1080
	Framerate = "30"
)

var numImage = 0

var x = 0
var y = 0

func main() {
	err := os.RemoveAll("./images/")
	if err != nil {
		panic(err)
	}
	err = os.Mkdir("./images/", os.ModePerm)
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile(FileName)
	if err != nil {
		panic(err)
	}

	imgFile, img := newFrame()

	for i := 0; i < len(file); i++ {
		for j := 0; j < 8; j++ {
			zeroOrOne := file[i] >> (7 - j) & 1
			PxColor := color.RGBA{R: 0, G: 0, B: 0}

			if zeroOrOne == 0 {
				PxColor = color.RGBA{R: 255, G: 255, B: 255}
			}

			img.Set(x, y, PxColor)

			if len(file) == i+1 && j == 7 {
				println("last")
				img.Set(x+1, y, color.RGBA{R: 255, G: 0, B: 0})
			}

			if x == PxWidth {
				if y == PxHeight {
					err = bmp.Encode(imgFile, img)
					if err != nil {
						panic(err)
					}

					imgFile, img = newFrame()
					x = 0
					y = 0
				} else {
					x = 0
					y++
				}
			} else {
				x++
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

func newFrame() (imgFile *os.File, img *image.Gray) {
	numImage++

	upLeft := image.Point{}
	lowRight := image.Point{X: PxWidth, Y: PxHeight}

	img = image.NewGray(image.Rectangle{Min: upLeft, Max: lowRight})
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
