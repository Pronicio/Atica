package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/draw"
	"math"
	"os"
	"os/exec"
	"strconv"
)

func decode() {
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

	var binary []byte

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

		println(len(bin))

		for i := 0; i < len(bin); i += 8 {
			end := false
			if i+8 >= len(bin) {
				i--
				end = true
			}

			var oct string
			for j := 0; j < 8; j++ {
				oct = oct + strconv.Itoa(bin[i+j])
			}

			octInt, _ := strconv.Atoi(oct)
			if octInt == 11111111 {
				continue
			}

			decimal := convertBinaryToDecimal(octInt)
			binary = append(binary, uint8(decimal))

			if end {
				println(i, i+8)
				break
			}
		}
	}

	writeInFile(binary)
}

func convertBinaryToDecimal(number int) int {
	decimal := 0
	counter := 0.0
	remainder := 0

	for number != 0 {
		remainder = number % 10
		decimal += remainder * int(math.Pow(2.0, counter))
		number = number / 10
		counter++
	}
	return decimal
}

func writeInFile(bin []byte) {
	err := os.WriteFile("./out/result.txt", bin, 0644)
	if err != nil {
		panic(err)
	}
}
