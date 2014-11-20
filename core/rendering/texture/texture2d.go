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