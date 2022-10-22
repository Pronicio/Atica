package main

import (
	"fmt"
	"golang.org/x/image/bmp"
	"os"
	"os/exec"
)

func decode() {
	err := os.RemoveAll("./out/imgs/")
	err = os.Mkdir("./out/imgs/", os.ModePerm)

	cmd := exec.Command("ffmpeg", "-i", "output.mp4", "./imgs/img-%d.bmp")
	cmd.Dir = "out/"
	err = cmd.Run()

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

		for i := 1; i < PxWidth*PxHeight; i++ {
			co := rgbaToDecimal(img.At(x, y).RGBA())

			if co == 0 {
				continue
			}

			if co == 255 {
				break
			}

			fmt.Print(co, " ")
			binary = append(binary, co)

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
	}

	writeInFile(binary)
}

func rgbaToDecimal(r uint32, g uint32, b uint32, a uint32) uint8 {
	return uint8(r / 257)
}

func writeInFile(bin []byte) {
	err := os.WriteFile("./out/result.txt", bin, 0644)
	if err != nil {
		panic(err)
	}
}
