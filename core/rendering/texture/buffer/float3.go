package buffer

import (
	"github.com/Opioid/scout/base/math"
)

type float3 struct {
	buffer
	data []math.Vector3
}

func newFloat3(dimensions math.Vector2i) *float3 {
	b := float3{}
	b.Resize(dimensions)
	return &b
}

func (b *float3) Type() uint32 {
	return Float3
}

func (b *float3) Resize(dimensions math.Vector2i) {
	b.dimensions = dimensions
	b.data = make([]math.Vector3, dimensions.X * dimensions.Y)
}

func (b *float3) At(x, y int32) math.Vector4 {
	return math.MakeVector4FromVector3(b.data[b.dimensions.X * y + x], 1.0)
}

func (b *float3) Set(x, y int32, color math.Vector4) {
	v := &b.data[b.dimensions.X * y + x]

	v.X = color.X
	v.Y = color.Y
	v.Z = color.Z
}

func (b *float3) SetRgb(x, y int32, color math.Vector3) {
	b.data[b.dimensions.X * y + x] = color
}

func (b *float3) SetChannel(x, y, c int32, value float32) {
	b.data[b.dimensions.X * y + x].Set(c, value)
}