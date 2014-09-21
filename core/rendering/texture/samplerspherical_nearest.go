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
	uv := math.MakeVector2((math.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, 1.0 - (xyz.Y + 1.0) * 0.5)
	x := int(uv.X * float32(s.texture.images[0].dimensions.X - 1))
	y := int(uv.Y * float32(s.texture.images[0].dimensions.Y - 1))
	return s.texture.images[0].at(x, y)
}