package surrounding

import (
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
	gomath "math"
	"os"
	"image/png"
	_ "fmt"
)

func BakeSphereMap(surrounding Surrounding, buffer *texture.Buffer) {
	ray := math.OptimizedRay{}
	ray.MaxT = 1000.0

	dimensions := buffer.Dimensions()

	sx := 1.0 / float32(dimensions.X) * gomath.Pi * 2.0
	sy := 1.0 / float32(dimensions.Y) * gomath.Pi

	for y := 0; y < dimensions.Y; y++ {
		ay := (float32(y) + 0.5) * sy

		vy := math.Cos(ay)
		say := -math.Sin(ay)

		for x := 0; x < dimensions.X; x++ {
			ax := (float32(x) + 0.5) * sx

			vx := say * math.Sin(ax)
			vz := say * math.Cos(ax)

			v := math.MakeVector3(vx, vy, vz)

			ray.SetDirection(v)

			c := surrounding.Sample(&ray)

			buffer.Set(x, y, math.MakeVector4(c.X, c.Y, c.Z, 1.0))
		}
	}

	image := buffer.RGBA()

	fo, err := os.Create("sphere_map.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)
}

func integrateHemisphereSphereMap(surrounding Surrounding, buffer *texture.Buffer) {
	ray := math.OptimizedRay{}
	ray.MaxT = 1000.0

	numSamples := uint32(256)
	numSamplesReciprocal := 1.0 / float32(numSamples)

	basis := math.Matrix3x3{}

	integrateHemisphere := func (d math.Vector3) math.Vector3 {
		basis.SetBasis(d)

		result := math.MakeVector3(0.0, 0.0, 0.0)

		for i := uint32(0); i < numSamples; i++ {
			sample := math.Hammersley(i, numSamples)

			s := math.HemisphereSample_cos(sample.X, sample.Y)

			v := basis.TransformVector3(s)
			ray.SetDirection(v)
			
			c := surrounding.Sample(&ray)

			result.AddAssign(c.Scale(numSamplesReciprocal))
		}

		return result
	}

	dimensions := buffer.Dimensions()

	sx := 1.0 / float32(dimensions.X) * gomath.Pi * 2.0
	sy := 1.0 / float32(dimensions.Y) * gomath.Pi

	for y := 0; y < dimensions.Y; y++ {
		ay := (float32(y) + 0.5) * sy

		vy := math.Cos(ay)
		say := -math.Sin(ay)

		for x := 0; x < dimensions.X; x++ {
			ax := (float32(x) + 0.5) * sx

			vx := say * math.Sin(ax)
			vz := say * math.Cos(ax)

			v := math.MakeVector3(vx, vy, vz)

			c := integrateHemisphere(v)

			buffer.Set(x, y, math.MakeVector4(c.X, c.Y, c.Z, 1.0))
		}
	}

	image := buffer.RGBA()

	fo, err := os.Create("sphere_map.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)
}