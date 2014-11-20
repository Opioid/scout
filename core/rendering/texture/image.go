package texture

import (
	"github.com/Opioid/scout/base/math"
	gomath "math"
)

type Image struct {
	Buffers []Buffer
}

func (i *Image) resize(dimensions math.Vector2i, numMipLevels uint32) {
	if numMipLevels <= 0 {
		numMipLevels = countMipLevels(dimensions)
	} else {
		numMipLevels = math.Minui(countMipLevels(dimensions), numMipLevels)
	}

	i.Buffers = make([]Buffer, numMipLevels)

	for l := uint32(0); l < numMipLevels; l++ {
		i.Buffers[l].Resize(dimensions)

		dimensions.X = math.Maxi(dimensions.X / 2, 1)
		dimensions.Y = math.Maxi(dimensions.Y / 2, 1)
	}
}

func (i *Image) NumMipLevels() uint32 {
	return uint32(len(i.Buffers))
}

func (i *Image) allocateMipLevels(numMipLevels uint32) {
	buffers := make([]Buffer, numMipLevels)

	copy(buffers, i.Buffers)

	previousMipLevels := i.NumMipLevels()

	dimensions := i.Buffers[previousMipLevels - 1].dimensions

	for l := previousMipLevels; l < numMipLevels; l++ {
		dimensions.X = math.Maxi(dimensions.X / 2, 1)
		dimensions.Y = math.Maxi(dimensions.Y / 2, 1)

		buffers[l].Resize(dimensions)
	}

	i.Buffers = buffers
}

func countMipLevels(dimensions math.Vector2i) uint32 {
	m := math.Maxi(dimensions.X, dimensions.Y)

	return 1 + uint32(gomath.Log2(float64(m)))
}