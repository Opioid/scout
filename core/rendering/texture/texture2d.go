package texture

import (
	"github.com/Opioid/scout/base/math"
)

type Texture2D struct {
	Image Image
}

func NewTexture2D(dimensions math.Vector2i, mipLevels int) *Texture2D {
	t := new(Texture2D)

	t.Image.resize(dimensions, mipLevels)

	return t
}

func (t *Texture2D) AllocateMipLevels(mipLevels int) {
	t.Image.allocateMipLevels(mipLevels)
}