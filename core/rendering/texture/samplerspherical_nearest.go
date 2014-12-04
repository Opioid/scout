package texture

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	gomath "math"
)

type SamplerSpherical_nearest struct {
	texture *Texture2D
}

func NewSamplerSpherical_nearest(t *Texture2D) *SamplerSpherical_nearest {
	s := new(SamplerSpherical_nearest)
	s.SetTexture(t)
	return s
}

func (sampler *SamplerSpherical_nearest) Texture() *Texture2D {
	return sampler.texture
}

func (sampler *SamplerSpherical_nearest) SetTexture(t *Texture2D) {
	sampler.texture = t
}

func (sampler *SamplerSpherical_nearest) Sample(xyz math.Vector3) math.Vector4 {
	uv := math.MakeVector2((math32.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1) * 0.5, math32.Acos(xyz.Y) / gomath.Pi)

	d := sampler.texture.Image.Buffers[0].dimensions

	x := math.Mini(int32(uv.X * float32(d.X)), d.X - 1)
	y := math.Mini(int32(uv.Y * float32(d.Y)), d.Y - 1)

	return sampler.texture.Image.Buffers[0].At(x, y)
}

func (sampler *SamplerSpherical_nearest) SampleLod(xyz math.Vector3, mipLevel float32) math.Vector4 {
	uv := math.MakeVector2((math32.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1) * 0.5, math32.Acos(xyz.Y) / gomath.Pi)

	b := &sampler.texture.Image.Buffers[int(mipLevel)]

	x := int32(uv.X * float32(b.dimensions.X - 1) + 0.5)
	y := int32(uv.Y * float32(b.dimensions.Y - 1) + 0.5)

	return b.At(x, y)
}