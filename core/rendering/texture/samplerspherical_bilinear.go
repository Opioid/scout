package texture

import (
	"github.com/Opioid/scout/base/math"
	gomath "math"
)

type SamplerSpherical_bilinear struct {
	texture *Texture2D
}

func NewSamplerSpherical_bilinear(t *Texture2D) *SamplerSpherical_bilinear {
	s := new(SamplerSpherical_bilinear)
	s.texture = t
	return s
}

func (sampler *SamplerSpherical_bilinear) Sample(xyz math.Vector3) math.Vector4 {
	uv := math.MakeVector2((math.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, 1.0 - (xyz.Y + 1.0) * 0.5)

	d := sampler.texture.Image.Buffers[0].dimensions

	u := uv.X * float32(d.X) - 0.5
	v := uv.Y * float32(d.Y) - 0.5

	fu := math.Floor(u)
	fv := math.Floor(v)

	x := int(fu)
	y := int(fv)

	x1 := x + 1

	if x1 >= d.X {
		x1 = 0
	}

	y1 := math.Mini(y + 1, d.Y - 1)

	if x < 0 {
		x = d.X - 1
	}

	y = math.Maxi(y, 0)	

	c00 := sampler.texture.Image.Buffers[0].At(x,  y)	
	c01 := sampler.texture.Image.Buffers[0].At(x,  y1)	
	c10 := sampler.texture.Image.Buffers[0].At(x1, y)	
	c11 := sampler.texture.Image.Buffers[0].At(x1, y1)

	s := u - fu
	t := v - fv

	return bilinear(c00, c01, c10, c11, s, t)
	
/*	uv := math.MakeVector2((math.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, 1.0 - (xyz.Y + 1.0) * 0.5)

	d := sampler.texture.Image.Buffers[0].dimensions

	x := math.Mini(int(uv.X * float32(d.X)), d.X - 1)
	y := math.Mini(int(uv.Y * float32(d.Y)), d.Y - 1)

	return sampler.texture.Image.Buffers[0].At(x, y)
	*/
}

func (s *SamplerSpherical_bilinear) SampleLod(xyz math.Vector3, mipLevel int) math.Vector4 {
	uv := math.MakeVector2((math.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, 1.0 - (xyz.Y + 1.0) * 0.5)

	b := &s.texture.Image.Buffers[mipLevel]

	x := int(uv.X * float32(b.dimensions.X - 1))
	y := int(uv.Y * float32(b.dimensions.Y - 1))

	return b.At(x, y)
}