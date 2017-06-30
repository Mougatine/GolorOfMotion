package main

import (
	"flag"

	"image/color"
	"image/jpeg"

	"image"

	"os"

	"github.com/lazywei/go-opencv/opencv"
)

type imageAverage struct {
	index   int
	average color.RGBA
}

func main() {
	flag.Parse()
	src := flag.Arg(0)

	srcCap := opencv.NewFileCapture(src)
	defer srcCap.Release()

	ch := make(chan imageAverage)
	i := 0
	for img := srcCap.QueryFrame(); img != nil; img = srcCap.QueryFrame() {
		go processImageAverage(img, ch, i)
		i++
	}

	imgs := make([]color.RGBA, i)
	for j := 0; j < i; j++ {
		p := <-ch
		imgs[p.index] = p.average
	}

	img := buildImage(imgs)

	toimg, _ := os.Create("out.jpg")
	jpeg.Encode(toimg, img, nil)
}

func buildImage(averages []color.RGBA) *image.RGBA {
	wa := 5
	waa := wa * len(averages)
	h := 100

	rect := image.Rectangle{image.Point{0, 0}, image.Point{waa, h}}
	img := image.NewRGBA(rect)

	for slice := 0; slice < len(averages); slice++ {
		pos := slice * wa
		for x := pos; x < pos+wa; x++ {
			for y := 0; y < h; y++ {
				img.Set(x, y, averages[slice])
			}
		}
	}

	return img
}

func processImageAverage(src *opencv.IplImage, ch chan imageAverage, index int) {
	r, g, b := 0., 0., 0.

	for x := 0; x < src.Width(); x++ {
		for y := 0; y < src.Height(); y++ {
			color := src.Get2D(x, y).Val()
			r += color[0]
			g += color[1]
			b += color[2]
		}
	}

	total := float64(src.Width() * src.Height())
	r, g, b = r/total, g/total, b/total
	r1, g1, b1, a1 := uint8(r), uint8(g), uint8(b), uint8(255)

	ch <- imageAverage{index, color.RGBA{r1, g1, b1, a1}}
}
