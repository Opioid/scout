package texture

import (
	_ "github.com/Opioid/scout/base/math"
)

type Texture2D struct {
	images []image
}

func NewTexture2D(numImages int) *Texture2D {
	t := new(Texture2D)
	t.images = make([]image, numImages)
	return t
}