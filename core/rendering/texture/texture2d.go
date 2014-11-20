package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Texture2D struct {
	texture
	Image Image
}

func NewTexture2D(dimensions math.Vector2i, mipLevels int32) *Texture2D {
	t := new(Texture2D)

	t.Image.resize(dimensions, mipLevels)

	return t
}

func (t *Texture2D) AllocateMipLevels(mipLevels int32) {
	t.Image.allocateMipLevels(mipLevels)
}

func (t *Texture2D) description() description {
	d := description{}

	d.textureType = Texture_2D

	dimensions := t.Image.Buffers[0].dimensions
	d.dimensions = math.MakeVector3i(dimensions.X, dimensions.Y, 1)

	return d
}