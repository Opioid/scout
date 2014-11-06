package texture

import (
	"github.com/Opioid/scout/base/math"
	gomath "math"
)

type image struct {
	buffers []Buffer
}

func (i *image) resize(dimensions math.Vector2i, mipLevels int) {
	if mipLevels <= 0 {
		mipLevels = countMipLevels(dimensions)
	} else {
		mipLevels = math.Mini(countMipLevels(dimensions), mipLevels)
	}

	i.buffers = make([]Buffer, mipLevels)

	for l := 0; l < mipLevels; l++ {
		i.buffers[l].Resize(dimensions)

		dimensions.X = math.Maxi(dimensions.X / 2, 1)
		dimensions.Y = math.Maxi(dimensions.Y / 2, 1)
	}
}

func (i *image) mipLevels() int {
	return len(i.buffers)
}

func (i *image) allocateMipLevels(mipLevels int) {
	buffers := make([]Buffer, mipLevels)

	copy(buffers, i.buffers)

	previousMipLevels := i.mipLevels()

	dimensions := i.buffers[previousMipLevels - 1].dimensions

	for l := previousMipLevels; l < mipLevels; l++ {
		buffers[l].Resize(dimensions)

		dimensions.X = math.Maxi(dimensions.X / 2, 1)
		dimensions.Y = math.Maxi(dimensions.Y / 2, 1)
	}

	i.buffers = buffers
}

func countMipLevels(dimensions math.Vector2i) int {
	m := math.Maxi(dimensions.X, dimensions.Y)

	return 1 + int(gomath.Log2(float64(m)))
}