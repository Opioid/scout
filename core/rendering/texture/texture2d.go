package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Texture2D struct {
	image image
}

func NewTexture2D(dimensions math.Vector2i, mipLevels int) *Texture2D {
	t := new(Texture2D)

	t.image.resize(dimensions, mipLevels)

	return t
}

func (t *Texture2D) AllocateMipLevels(mipLevels int) {
	t.image.allocateMipLevels(mipLevels)
}