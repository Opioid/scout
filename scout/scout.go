package main

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	"github.com/Opioid/scout/core/rendering"
	"fmt"
	"os"
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

	dimensions := math.Vector2{320, 240}
	buffer := rendering.NewPixelBuffer(dimensions)

	buffer.Set(0, 0, math.Vector4{0.5, 0.3, 0.2, 1.0})

	fo, err := os.Create("output.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

//	out := bufio.NewWriter(fo)

//	defer out.Flush()

//	out.WriteString("Umpa lumpa")


//	image := image.NewRGBA(image.Rect(0, 0, dimensions.X, dimensions.Y))

	for y := 0; y < dimensions.Y; y++ {
		for x := 0; x < dimensions.X; x++ {
			r := float32(y) / float32(dimensions.Y)
			g := float32(x) / float32(dimensions.X)

			buffer.Set(x, y, math.Vector4{r, g, 0.5, 1.0})
		}
	}
	
	image := buffer.RGBA()

	png.Encode(fo, image)
}