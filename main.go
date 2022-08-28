package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"os/exec"
	"strconv"
)

const (
	FileName  = "text.txt"
	PxWidth   = 1080
	PxHeight  = 720
	Framerate = "10"
)

var numImage = 0

var x = 0
var y = 0

func main() {
	file, err := os.ReadFile(FileName)
	if err != nil {
		fmt.Print(err)
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

			if x == PxWidth {
				if y == PxHeight {
					opt := jpeg.Options{
						Quality: 100,
					}

					err = jpeg.Encode(imgFile, img, &opt)
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

	opt := jpeg.Options{
		Quality: 100,
	}

	err = jpeg.Encode(imgFile, img, &opt)
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
	filename := "./images/img-" + strconv.Itoa(numImage) + ".jpg"

	imgFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	draw.Draw(img, img.Bounds(), img, image.Point{}, draw.Over)
	return
}

func toVideo() {
	cmd := exec.Command("ffmpeg", "-framerate", Framerate, "-i", "img-%d.jpg", "output.mp4")
	cmd.Dir = "images/"
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	println(out)
}
