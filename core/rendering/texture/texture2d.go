package texture

import (
	"github.com/Opioid/scout/core/rendering/texture/buffer"
	"github.com/Opioid/scout/base/math"
)

type Texture2D struct {
	texture
	Image Image
	MaxMipLevel uint32
}

func NewTexture2D(t uint32, dimensions math.Vector2i, numMipLevels uint32) *Texture2D {
	tex := Texture2D{}
	tex.Image = MakeImage(t, dimensions, numMipLevels)
	return &tex
}

func NewTexture2DFromDescription(description *description) *Texture2D {
	return NewTexture2D(buffer.Float4, description.Dimensions.Vector2i(), description.NumMipLevels)
}

func (t *Texture2D) AllocateMipLevels(numMipLevels uint32) {
	t.Image.allocateMipLevels(numMipLevels)
	t.updateMaxMipLevel()
}

func (t *Texture2D) AllocateMipLevelsDownTo(dimensions math.Vector2i) {
	t.Image.allocateMipLevelsDownTo(dimensions)
	t.updateMaxMipLevel()
}

func (t *Texture2D) updateMaxMipLevel() {
	t.MaxMipLevel = t.Image.NumMipLevels() - 1
}

func (t *Texture2D) description() description {
	d := description{}

	d.TextureType = Texture_2D

	dimensions := t.Image.Buffers[0].Dimensions()
	d.Dimensions = math.MakeVector3i(dimensions.X, dimensions.Y, 1)
	d.NumMipLevels = t.Image.NumMipLevels()

	return d
}