package texture

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	gomath "math"
	_ "fmt"
)

type SamplerSpherical_linear struct {
}

func NewSamplerSpherical_linear() *SamplerSpherical_linear {
	sampler := new(SamplerSpherical_linear)
	
	return sampler
}

func (sampler *SamplerSpherical_linear) Sample(texture *Texture2D, xyz math.Vector3) math.Vector4 {
	// atan2(0, 0) is undefined, this is an easy way to fix this for a notorious problems with xyz == [0, 1, 0] in scenes with a ground plane
	if xyz.X == 0.0 {
		xyz.X = 2.0 * gomath.Pi 
	}

	uv := math.MakeVector2((math32.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math32.Acos(xyz.Y) / gomath.Pi)

	return sampler.sampleLevel(texture, uv, 0)
}

func (sampler *SamplerSpherical_linear) Sample3(texture *Texture2D, xyz math.Vector3) math.Vector3 {
	// atan2(0, 0) is undefined, this is an easy way to fix this for a notorious problems with xyz == [0, 1, 0] in scenes with a ground plane
	if xyz.X == 0.0 {
		xyz.X = 2.0 * gomath.Pi 
	}

	uv := math.MakeVector2((math32.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math32.Acos(xyz.Y) / gomath.Pi)

	return sampler.sampleLevel3(texture, uv, 0)
}

func (sampler *SamplerSpherical_linear) SampleLod(texture *Texture2D, xyz math.Vector3, mipLevel float32) math.Vector4 {
	// atan2(0, 0) is undefined, this is an easy way to fix this for a notorious problems with xyz == [0, 1, 0] in scenes with a ground plane
	if xyz.X == 0.0 {
		xyz.X = 2.0 * gomath.Pi 
	}

	uv := math.MakeVector2((math32.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math32.Acos(xyz.Y) / gomath.Pi)

	l0 := math32.Floor(mipLevel)

	l0i := uint32(l0)
	l1i := math.Minui(l0i + 1, texture.MaxMipLevel)

	c0 := sampler.sampleLevel(texture, uv, l0i)
	c1 := sampler.sampleLevel(texture, uv, l1i)

	return c0.Lerp(c1, mipLevel - l0)
}

func (sampler *SamplerSpherical_linear) SampleLod3(texture *Texture2D, xyz math.Vector3, mipLevel float32) math.Vector3 {
	// atan2(0, 0) is undefined, this is an easy way to fix this for a notorious problems with xyz == [0, 1, 0] in scenes with a ground plane
	if xyz.X == 0.0 {
		xyz.X = 2.0 * gomath.Pi 
	}

	uv := math.MakeVector2((math32.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math32.Acos(xyz.Y) / gomath.Pi)

	l0 := math32.Floor(mipLevel)

	l0i := uint32(l0)
	l1i := math.Minui(l0i + 1, texture.MaxMipLevel)

	c0 := sampler.sampleLevel3(texture, uv, l0i)
	c1 := sampler.sampleLevel3(texture, uv, l1i)

	return c0.Lerp(c1, mipLevel - l0)
}

func (sampler *SamplerSpherical_linear) sampleLevel(texture *Texture2D, uv math.Vector2, mipLevel uint32) math.Vector4 {
	b := texture.Image.Buffers[mipLevel]
	dimensions := b.Dimensions()

	u := uv.X * float32(dimensions.X) - 0.5
	v := uv.Y * float32(dimensions.Y) - 0.5

	fu := math32.Floor(u)
	fv := math32.Floor(v)

	x := int32(fu)
	y := int32(fv)

	if x < 0 {
		x = dimensions.X - 1
	}

	x1 := x + 1

	if x1 >= dimensions.X {
		x1 = 0
	}

	y1 := math.Mini(y + 1, dimensions.Y - 1)

	y = math.Maxi(y, 0)	

	c00 := b.At(x,  y)	
	c01 := b.At(x,  y1)	
	c10 := b.At(x1, y)	
	c11 := b.At(x1, y1)

	s := u - fu
	t := v - fv

	return bilinear(c00, c01, c10, c11, s, t)
}

func (sampler *SamplerSpherical_linear) sampleLevel3(texture *Texture2D, uv math.Vector2, mipLevel uint32) math.Vector3 {
	b := texture.Image.Buffers[mipLevel]
	dimensions := b.Dimensions()

	u := uv.X * float32(dimensions.X) - 0.5
	v := uv.Y * float32(dimensions.Y) - 0.5

	fu := math32.Floor(u)
	fv := math32.Floor(v)

	x := int32(fu)
	y := int32(fv)

	if x < 0 {
		x = dimensions.X - 1
	}

	x1 := x + 1

	if x1 >= dimensions.X {
		x1 = 0
	}

	y1 := math.Mini(y + 1, dimensions.Y - 1)

	y = math.Maxi(y, 0)	

	c00 := b.At3(x,  y)	
	c01 := b.At3(x,  y1)	
	c10 := b.At3(x1, y)	
	c11 := b.At3(x1, y1)

	s := u - fu
	t := v - fv

	return bilinear3(c00, c01, c10, c11, s, t)
}