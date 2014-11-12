package surrounding

import (
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	gomath "math"
	"os"
	"image/png"
	"runtime"
	"sync"	
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

			c, _ := surrounding.Sample(&ray)

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

func calculateSphereMapSolidAngleWeights(buffer *texture.Buffer) {
	dimensions := buffer.Dimensions()

//	sx := 1.0 / float32(dimensions.X) * gomath.Pi * 2.0
	sy := 1.0 / float32(dimensions.Y) * gomath.Pi

	for y := 0; y < dimensions.Y; y++ {
		ay := (float32(y) + 0.5) * sy
	//	ay1 := (float32(y) + 1.0) * sy

	//	vy := math.Cos(ay)
		say := /*float32(1.0)//*/math.Sin(ay)

		for x := 0; x < dimensions.X; x++ {
		//	ax0 := (float32(x) + 0.5) * sx
		//	ax1 := (float32(x) + 1.0) * sx

		//	angles := math.MakeVector2(ax0, ay0)
		//	fmt.Println(math.MakeVector2(ax0, ay0))
		//	fmt.Println(math.MakeVector2(ax1, ay1))

		//	vx := say * math.Sin(ax)
		//	vz := say * math.Cos(ax)

		//	sa := (math.Sin(ay1) - math.Sin(ay0)) * (ax1 - ax0)

		//	fmt.Println(ay0 - ay1)

		//	sa := float32(0.5)

			buffer.SetChannel(x, y, 3, say)
		}

		//fmt.Println(math.Sin(ay))
	}

	image := buffer.RGBA()

	fo, err := os.Create("solid_angles.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)
}

func integrateHemisphereSphereMap(surrounding Surrounding, numSamples uint32, buffer *texture.Buffer) {
	numTaks := runtime.GOMAXPROCS(0)

	dimensions := buffer.Dimensions()

	a := dimensions.Y / numTaks

	start := math.MakeVector2i(0, 0)
	end   := math.MakeVector2i(dimensions.X, a)

	wg := sync.WaitGroup{}

	for i := 0; i < numTaks; i++ {
		wg.Add(1)

		go func (s, e math.Vector2i) {
			integrateHemisphereSphereMapTask(surrounding, numSamples, s, e, buffer)
			wg.Done()
		}(start, end)

		start.Y += a

		if i == numTaks - 2 {
			end.Y = dimensions.Y
		} else {
			end.Y += a
		}
	}

	wg.Wait()

	image := buffer.RGBA()

	fo, err := os.Create("sphere_map.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)
}

func integrateHemisphereSphereMapTask(surrounding Surrounding, numSamples uint32, start, end math.Vector2i, buffer *texture.Buffer) {
	rng := random.Generator{}
	rng.Seed(uint32(start.X) + 0, uint32(start.Y) + 1, uint32(start.X) + 2, uint32(start.Y) + 3)	

	ray := math.OptimizedRay{}
	ray.MaxT = 1000.0

//	numSamplesReciprocal := 1.0 / float32(numSamples)

	basis := math.Matrix3x3{}

	integrateHemisphere := func (d math.Vector3) math.Vector3 {
		basis.SetBasis(d)

		result := math.MakeVector3(0.0, 0.0, 0.0)

		weightSum := float32(0.0)

		rn := rng.RandomUint32()

		for i := uint32(0); i < numSamples; i++ {
			sample := math.ScrambledHammersley(i, numSamples, rn)

			s := math.HemisphereSample_cos(sample.X, sample.Y)

			v := basis.TransformVector3(s)
			ray.SetDirection(v)
			
			c, w := surrounding.Sample(&ray)

			weightSum += w

			result.AddAssign(c.Scale(w))
		//	result.AddAssign(c.Scale(numSamplesReciprocal))
		}

		result.ScaleAssign(1.0 / weightSum)

		return result
	}

	dimensions := buffer.Dimensions()

	sx := 1.0 / float32(dimensions.X) * gomath.Pi * 2.0
	sy := 1.0 / float32(dimensions.Y) * gomath.Pi

	for y := start.Y; y < end.Y; y++ {
		ay := (float32(y) + 0.5) * sy

		vy := math.Cos(ay)
		say := -math.Sin(ay)

		for x := start.X; x < end.X; x++ {
			ax := (float32(x) + 0.5) * sx

			vx := say * math.Sin(ax)
			vz := say * math.Cos(ax)

			v := math.MakeVector3(vx, vy, vz)

			c := integrateHemisphere(v)

			buffer.Set(x, y, math.MakeVector4(c.X, c.Y, c.Z, 1.0))
		}
	}
}