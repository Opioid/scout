package texture

import (
	"github.com/Opioid/scout/base/math"
	gomath "math"
	_ "fmt"
)

type SamplerSpherical_linear struct {
	texture *Texture2D
	maxMipLevel int32
}

func NewSamplerSpherical_linear(t *Texture2D) *SamplerSpherical_linear {
	sampler := new(SamplerSpherical_linear)
	sampler.texture = t
	sampler.maxMipLevel = sampler.texture.Image.MipLevels() - 1
	return sampler
}

func (sampler *SamplerSpherical_linear) Texture() *Texture2D {
	return sampler.texture
}

func (sampler *SamplerSpherical_linear) Sample(xyz math.Vector3) math.Vector4 {
	uv := math.MakeVector2((math.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math.Acos(xyz.Y) / gomath.Pi)

	return sampler.sampleLevel(uv, 0)
}

func (sampler *SamplerSpherical_linear) SampleLod(xyz math.Vector3, mipLevel float32) math.Vector4 {
	uv := math.MakeVector2((math.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math.Acos(xyz.Y) / gomath.Pi)

	l0 := math.Floor(mipLevel)

	l0i := int32(l0)
	l1i := math.Maxi(int32(l0) + 1, sampler.maxMipLevel)

	c0 := sampler.sampleLevel(uv, l0i)
	c1 := sampler.sampleLevel(uv, l1i)

//	fmt.Printf("%v %v\n", c0, c1)

	return c0.Lerp(c1, mipLevel - l0)
}

func (sampler *SamplerSpherical_linear) sampleLevel(uv math.Vector2, mipLevel int32) math.Vector4 {
	b := &sampler.texture.Image.Buffers[mipLevel]

	u := uv.X * float32(b.dimensions.X) - 0.5
	v := uv.Y * float32(b.dimensions.Y) - 0.5

	fu := math.Floor(u)
	fv := math.Floor(v)

	x := int32(fu)
	y := int32(fv)

	x1 := x + 1

	if x1 >= b.dimensions.X {
		x1 = 0
	}

	y1 := math.Mini(y + 1, b.dimensions.Y - 1)

	if x < 0 {
		x = b.dimensions.X - 1
	}

	y = math.Maxi(y, 0)	

	c00 := b.At(x,  y)	
	c01 := b.At(x,  y1)	
	c10 := b.At(x1, y)	
	c11 := b.At(x1, y1)

	s := u - fu
	t := v - fv

	return bilinear(c00, c01, c10, c11, s, t)
}