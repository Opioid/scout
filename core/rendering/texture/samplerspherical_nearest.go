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

func (s *SamplerSpherical_nearest) Sample(xyz math.Vector3) math.Vector4 {
	uv := math.MakeVector2((math.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math.Acos(xyz.Y) / gomath.Pi)

	d := s.texture.Image.Buffers[0].dimensions

	x := math.Mini(int(uv.X * float32(d.X)), d.X - 1)
	y := math.Mini(int(uv.Y * float32(d.Y)), d.Y - 1)

	return s.texture.Image.Buffers[0].At(x, y)
}

func (s *SamplerSpherical_nearest) SampleLod(xyz math.Vector3, mipLevel int) math.Vector4 {
	uv := math.MakeVector2((math.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math.Acos(xyz.Y) / gomath.Pi)

	b := &s.texture.Image.Buffers[mipLevel]

	x := int(uv.X * float32(b.dimensions.X - 1) + 0.5)
	y := int(uv.Y * float32(b.dimensions.Y - 1) + 0.5)

	return b.At(x, y)
}