package ibl

import (
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/material/ggx"
	"github.com/Opioid/scout/core/scene/surrounding"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	gomath "math"
	"runtime"
	"sync"	
	_ "os"
	_ "image/png"	
	_ "strconv"
	_ "fmt"
)

func CalculateSphereMapSolidAngleWeights(buffer *texture.Buffer) {
	dimensions := buffer.Dimensions()

//	sx := 1.0 / float32(dimensions.X) * gomath.Pi * 2.0
	sy := 1.0 / float32(dimensions.Y) * gomath.Pi

	for y := int32(0); y < dimensions.Y; y++ {
		ay := (float32(y) + 0.5) * sy
	//	ay1 := (float32(y) + 1.0) * sy

	//	vy := math.Cos(ay)
		say := /*float32(1.0)//*/math.Sin(ay)

		for x := int32(0); x < dimensions.X; x++ {
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
/*
	image := buffer.RGBA()

	fo, err := os.Create("solid_angles.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)
	*/
}

func IntegrateHemisphereSphereMap(surrounding surrounding.Surrounding, numSamples uint32, buffer *texture.Buffer) {
	dimensions := buffer.Dimensions()

	numTaks := int32(runtime.GOMAXPROCS(0))

	a := dimensions.Y / numTaks

	start := math.MakeVector2i(0, 0)
	end   := math.MakeVector2i(dimensions.X, a)

	wg := sync.WaitGroup{}

	for i := int32(0); i < numTaks; i++ {
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
/*
	image := buffer.RGBA()

	fo, err := os.Create("sphere_map_diffuse.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)
	*/
}

func integrateHemisphereSphereMapTask(surrounding surrounding.Surrounding, numSamples uint32, start, end math.Vector2i, buffer *texture.Buffer) {
	rng := random.Generator{}
	rng.Seed(uint32(start.X) + 0, uint32(start.Y) + 1, uint32(start.X) + 2, uint32(start.Y) + 3)	

	ray := math.OptimizedRay{}
	ray.MaxT = 1000.0

//	numSamplesReciprocal := 1.0 / float32(numSamples)

	basis := math.Matrix3x3{}

	integrateHemisphere := func (n math.Vector3) math.Vector3 {
		basis.SetBasis(n)

		result := math.MakeVector3(0, 0, 0)

		weightSum := float32(0)

		rn := rng.RandomUint32()

		for i := uint32(0); i < numSamples; i++ {
			sample := math.ScrambledHammersley(i, numSamples, rn)

			s := math.HemisphereSample_cos(sample.X, sample.Y)

			v := basis.TransformVector3(s)
			ray.SetDirection(v)
			
			c, w := surrounding.SampleSecondary(&ray)

			weightSum += w

			result.AddAssign(c.Scale(w))
		//	result.AddAssign(c.Scale(numSamplesReciprocal))
		}

		result.ScaleAssign(1.0 / weightSum)

		return result
	}

	dimensions := buffer.Dimensions()

	sx := 1 / float32(dimensions.X) * gomath.Pi * 2
	sy := 1 / float32(dimensions.Y) * gomath.Pi

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

			buffer.Set(x, y, math.MakeVector4(c.X, c.Y, c.Z, 1))
		}
	}
}

func IntegrateConeSphereMap(surrounding surrounding.Surrounding, roughness float32, numSamples uint32, buffer *texture.Buffer) {
	dimensions := buffer.Dimensions()

	numTaks := int32(runtime.GOMAXPROCS(0))

	a := dimensions.Y / numTaks

	start := math.MakeVector2i(0, 0)
	end   := math.MakeVector2i(dimensions.X, a)

	wg := sync.WaitGroup{}

	for i := int32(0); i < numTaks; i++ {
		wg.Add(1)

		go func (s, e math.Vector2i) {
			integrateConeSphereMapTask(surrounding, roughness, numSamples, s, e, buffer)
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
/*
	image := buffer.RGBA()

	fo, err := os.Create("sphere_map_specular_" + strconv.FormatFloat(float64(roughness), 'f', 2, 32) + ".png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)
	*/
}

func integrateConeSphereMapTask(surrounding surrounding.Surrounding, roughness float32, numSamples uint32, start, end math.Vector2i, buffer *texture.Buffer) {
	rng := random.Generator{}
	rng.Seed(uint32(start.X) + 0, uint32(start.Y) + 1, uint32(start.X) + 2, uint32(start.Y) + 3)	

	ray := math.OptimizedRay{}
	ray.MaxT = 1000.0

//	numSamplesReciprocal := 1.0 / float32(numSamples)

	integrateCone := func (n math.Vector3) math.Vector3 {
		v := n

		result := math.MakeVector3(0, 0, 0)

		weightSum := float32(0)

		rn := rng.RandomUint32()

		for i := uint32(0); i < numSamples; i++ {
			xi := math.ScrambledHammersley(i, numSamples, rn)
			h  := ggx.ImportanceSample(xi, roughness, n)

			// normalizing here prevents some NaN where l.Y is beyond either -1 or 1
			l := h.Scale(2 * v.Dot(h)).Sub(v).Normalized()

			n_dot_l := math.Saturate(n.Dot(l))

			if n_dot_l > 0 {
				ray.SetDirection(l)
			
				c, _ := surrounding.SampleSecondary(&ray)

				result.AddAssign(c.Scale(n_dot_l))

				weightSum += n_dot_l
			}

		//	vec3 L = 2 * dot( V, H ) * H - V;
		//	float NoL = saturate( dot( N, L ) );
		//	if (NoL > 0.f)
		//	{
		//		PrefilteredColor += texture(g_cubemap, L).rgb * NoL;
		//		TotalWeight += NoL;
		//	}
		}

		result.ScaleAssign(1 / weightSum)

		return result
	}

	dimensions := buffer.Dimensions()

	sx := 1 / float32(dimensions.X) * gomath.Pi * 2
	sy := 1 / float32(dimensions.Y) * gomath.Pi

	for y := start.Y; y < end.Y; y++ {
		ay := (float32(y) + 0.5) * sy

		vy := math.Cos(ay)
		say := -math.Sin(ay)

		for x := start.X; x < end.X; x++ {
			ax := (float32(x) + 0.5) * sx

			vx := say * math.Sin(ax)
			vz := say * math.Cos(ax)

			v := math.MakeVector3(vx, vy, vz)

			c := integrateCone(v)

			buffer.Set(x, y, math.MakeVector4(c.X, c.Y, c.Z, 1))
		}
	}
}