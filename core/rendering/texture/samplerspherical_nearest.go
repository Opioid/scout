package texture

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	gomath "math"
)

type SamplerSpherical_nearest struct {
}

func NewSamplerSpherical_nearest() *SamplerSpherical_nearest {
	s := new(SamplerSpherical_nearest)
	return s
}

func (sampler *SamplerSpherical_nearest) Sample(texture *Texture2D, xyz math.Vector3) math.Vector4 {
	// atan2(0, 0) is undefined, this is an easy way to fix this for a notorious problems with xyz == [0, 1, 0] in scenes with a ground plane
	if xyz.X == 0.0 {
		xyz.X = 0.000001//2.0 * gomath.Pi
	}

	uv := math.MakeVector2((math32.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math32.Acos(xyz.Y) / gomath.Pi)

	d := texture.Image.Buffers[0].Dimensions()

	x := math.Mini(int32(uv.X * float32(d.X)), d.X - 1)
	y := math.Mini(int32(uv.Y * float32(d.Y)), d.Y - 1)

	return texture.Image.Buffers[0].At(x, y)
}

func (sampler *SamplerSpherical_nearest) Sample3(texture *Texture2D, xyz math.Vector3) math.Vector3 {
	// atan2(0, 0) is undefined, this is an easy way to fix this for a notorious problems with xyz == [0, 1, 0] in scenes with a ground plane
	if xyz.X == 0.0 {
		xyz.X = 0.000001//2.0 * gomath.Pi
	}

	uv := math.MakeVector2((math32.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1.0) * 0.5, math32.Acos(xyz.Y) / gomath.Pi)

	d := texture.Image.Buffers[0].Dimensions()

	x := math.Mini(int32(uv.X * float32(d.X)), d.X - 1)
	y := math.Mini(int32(uv.Y * float32(d.Y)), d.Y - 1)

	return texture.Image.Buffers[0].At3(x, y)
}

func (sampler *SamplerSpherical_nearest) SampleLod(texture *Texture2D, xyz math.Vector3, mipLevel float32) math.Vector4 {
	// atan2(0, 0) is undefined, this is an easy way to fix this for a notorious problems with xyz == [0, 1, 0] in scenes with a ground plane
	if xyz.X == 0.0 {
		xyz.X = 0.000001//2.0 * gomath.Pi
	}

	uv := math.MakeVector2((math32.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1) * 0.5, math32.Acos(xyz.Y) / gomath.Pi)

	b := texture.Image.Buffers[int(mipLevel)]
	dimensions := b.Dimensions()

	x := int32(uv.X * float32(dimensions.X - 1) + 0.5)
	y := int32(uv.Y * float32(dimensions.Y - 1) + 0.5)

	return b.At(x, y)
}

func (sampler *SamplerSpherical_nearest) SampleLod3(texture *Texture2D, xyz math.Vector3, mipLevel float32) math.Vector3 {
	// atan2(0, 0) is undefined, this is an easy way to fix this for a notorious problems with xyz == [0, 1, 0] in scenes with a ground plane
	if xyz.X == 0.0 {
		xyz.X = 0.000001//2.0 * gomath.Pi
	}

	uv := math.MakeVector2((math32.Atan2(xyz.X, xyz.Z) / gomath.Pi + 1) * 0.5, math32.Acos(xyz.Y) / gomath.Pi)

	b := texture.Image.Buffers[int(mipLevel)]
	dimensions := b.Dimensions()

	x := int32(uv.X * float32(dimensions.X - 1) + 0.5)
	y := int32(uv.Y * float32(dimensions.Y - 1) + 0.5)

	return b.At3(x, y)
}