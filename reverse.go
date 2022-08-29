package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/draw"
	"os"
	"os/exec"
	"path/filepath"
)

func reverse() {
	cmd := exec.Command("ffmpeg", "-i", "output.mp4", "-pix_fmt", "bgr8", "./imgs/img-%d.bmp")
	cmd.Dir = "out/"
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	files, err := os.ReadDir("./out/imgs/")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		data, err := os.Open("./out/imgs/" + file.Name())
		if err != nil {
			panic(err)
		}

		img, err := bmp.Decode(data)
		if err != nil {
			panic(err)
		}

		rect := img.Bounds()
		cimg := image.NewGray(rect)
		draw.Draw(cimg, rect, img, rect.Min, draw.Src)

		var bin []int

		for i := 1; i < PxWidth*PxHeight; i++ {
			co := cimg.At(x, y)
			ct := fmt.Sprint(co.RGBA())

			if ct == "65021 65021 65021 65535" {
				bin = append(bin, 0)
			} else if ct == "0 0 0 65535" {
				bin = append(bin, 1)
			} else if ct == "26985 26985 26985 65535" {
				break
			} else {
				bin = append(bin, 1)
			}

			if x == PxWidth {
				if y == PxHeight {
					x = 0
					y = 0
					println("Finish")
				} else {
					x = 0
					y++
				}
			} else {
				x++
			}
		}

		println(bin)
	}
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		*files = append(*files, path)
		return nil
	}
}
