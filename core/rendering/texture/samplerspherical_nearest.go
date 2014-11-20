package texture

import (
	"github.com/Opioid/scout/base/math"
	gomath "math"
)

type SamplerSpherical_nearest struct {
	texture *Texture2D
}

func NewSamplerSpherical_nearest(t *Texture2D) *SamplerSpherical_nearest {
	s := new(SamplerSpherical_nearest)
	s.texture = t
	return s
}

func (sampler *SamplerSpherical_nearest) Texture() *Texture2D {
	return sampler.texture
}

func (sampler *SamplerSpherical_nearest) Sample(xyz math.Vector3) math.Vector4 {
	uv := math.MakeVector2((math.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math.Acos(xyz.Y) / gomath.Pi)

	d := sampler.texture.Image.Buffers[0].dimensions

	x := math.Mini(int32(uv.X * float32(d.X)), d.X - 1)
	y := math.Mini(int32(uv.Y * float32(d.Y)), d.Y - 1)

	return sampler.texture.Image.Buffers[0].At(x, y)
}

func (sampler *SamplerSpherical_nearest) SampleLod(xyz math.Vector3, mipLevel float32) math.Vector4 {
	uv := math.MakeVector2((math.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math.Acos(xyz.Y) / gomath.Pi)

	b := &sampler.texture.Image.Buffers[int(mipLevel)]

	x := int32(uv.X * float32(b.dimensions.X - 1) + 0.5)
	y := int32(uv.Y * float32(b.dimensions.Y - 1) + 0.5)

	return b.At(x, y)
}