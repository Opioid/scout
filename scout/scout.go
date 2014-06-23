package main

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	"fmt"
)

func main() {

	a := math.Vector3{1.0, 2.0,  3.0}
	b := math.Vector3{4.0, 4.0, -8.0}

	c := math.V3Add(a, b)

	s := bounding.Sphere{c, 2.0}

	fmt.Printf("Sphere %v\n", s)

	fmt.Printf("The result is %v, that's a vector\n", c)

	fmt.Printf("dot(a, b) == %f\n", math.V3Dot(a, b))
}