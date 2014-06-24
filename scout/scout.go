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

	dimensions := math.Vector2{1280, 720}
	buffer := rendering.NewPixelBuffer(dimensions)

	renderer := rendering.Renderer{}

	renderer.Render(buffer)

	fo, err := os.Create("output.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

//	out := bufio.NewWriter(fo)

//	defer out.Flush()

//	out.WriteString("Umpa lumpa")
	
	image := buffer.RGBA()

	png.Encode(fo, image)
}