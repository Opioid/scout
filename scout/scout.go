package main

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	"fmt"
	"os"
	"image"
	"image/color"
	"image/png"
)

func main() {

	a := math.Vector3{1.0, 2.0,  3.0}
	b := math.Vector3{4.0, 4.0, -8.0}

	c := a.Add(b)

	s := bounding.Sphere{c, 2.0}

	fmt.Printf("Sphere %v\n", s)

	fmt.Printf("The result is %v, that's a vector\n", c)

	fmt.Printf("a.Dot(b) == %f\n", a.Dot(b))

	fmt.Printf("c.Length == %f\n", c.SquaredLength())

	fo, err := os.Create("output.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

//	out := bufio.NewWriter(fo)

//	defer out.Flush()

//	out.WriteString("Umpa lumpa")

	width := 320
	height := 240

	image := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8(255.0 * (float32(y) / float32(height)))
			g := uint8(255.0 * (float32(x) / float32(width)))

			image.Set(x, y, color.RGBA{r, g, 127, 255})
		}
	}

	png.Encode(fo, image)
}