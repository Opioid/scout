package texture

import (
	"github.com/Opioid/scout/base/math"
	gomath "math"
)

type Image struct {
	Buffers []Buffer
}

func (i *Image) resize(dimensions math.Vector2i, mipLevels int32) {
	if mipLevels <= 0 {
		mipLevels = countMipLevels(dimensions)
	} else {
		mipLevels = math.Mini(countMipLevels(dimensions), mipLevels)
	}

	i.Buffers = make([]Buffer, mipLevels)

	for l := int32(0); l < mipLevels; l++ {
		i.Buffers[l].Resize(dimensions)

		dimensions.X = math.Maxi(dimensions.X / 2, 1)
		dimensions.Y = math.Maxi(dimensions.Y / 2, 1)
	}
}

func (i *Image) MipLevels() int32 {
	return int32(len(i.Buffers))
}

func (i *Image) allocateMipLevels(mipLevels int32) {
	buffers := make([]Buffer, mipLevels)

	copy(buffers, i.Buffers)

	previousMipLevels := i.MipLevels()

	dimensions := i.Buffers[previousMipLevels - 1].dimensions

	for l := previousMipLevels; l < mipLevels; l++ {
		dimensions.X = math.Maxi(dimensions.X / 2, 1)
		dimensions.Y = math.Maxi(dimensions.Y / 2, 1)

		buffers[l].Resize(dimensions)
	}

	i.Buffers = buffers
}

func countMipLevels(dimensions math.Vector2i) int32 {
	m := math.Maxi(dimensions.X, dimensions.Y)

	return 1 + int32(gomath.Log2(float64(m)))
}