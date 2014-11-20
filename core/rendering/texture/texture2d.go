package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Texture2D struct {
	texture
	Image Image
}

func NewTexture2D(dimensions math.Vector2i, numMipLevels uint32) *Texture2D {
	t := new(Texture2D)

	t.Image.resize(dimensions, numMipLevels)

	return t
}

func NewTexture2DFromDescription(description *description) *Texture2D {
	return NewTexture2D(description.Dimensions.Vector2i(), description.NumMipLevels)
}

func (t *Texture2D) AllocateMipLevels(numMipLevels uint32) {
	t.Image.allocateMipLevels(numMipLevels)
}

func (t *Texture2D) description() description {
	d := description{}

	d.TextureType = Texture_2D

	dimensions := t.Image.Buffers[0].dimensions
	d.Dimensions = math.MakeVector3i(dimensions.X, dimensions.Y, 1)
	d.NumMipLevels = t.Image.NumMipLevels()

	return d
}