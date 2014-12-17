package texture

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/core/rendering/texture/buffer"
	gomath "math"
)

type Image struct {
	Buffers []buffer.Buffer
}

func MakeImage(t uint32, dimensions math.Vector2i, numMipLevels uint32) Image {
	i := Image{}

	if numMipLevels <= 0 {
		numMipLevels = countMipLevels(dimensions)
	} else {
		numMipLevels = math.Minui(countMipLevels(dimensions), numMipLevels)
	}

	i.Buffers = make([]buffer.Buffer, numMipLevels)

	for l := uint32(0); l < numMipLevels; l++ {
		i.Buffers[l] = buffer.New(t, dimensions)

		dimensions.X = math.Maxi(dimensions.X / 2, 1)
		dimensions.Y = math.Maxi(dimensions.Y / 2, 1)
	}

	return i
}

func (i *Image) NumMipLevels() uint32 {
	return uint32(len(i.Buffers))
}

func (i *Image) allocateMipLevels(numMipLevels uint32) {
	buffers := make([]buffer.Buffer, numMipLevels)

	copy(buffers, i.Buffers)

	previousMipLevels := i.NumMipLevels()

	dimensions := i.Buffers[previousMipLevels - 1].Dimensions()

	for l := previousMipLevels; l < numMipLevels; l++ {
		dimensions.X = math.Maxi(dimensions.X / 2, 1)
		dimensions.Y = math.Maxi(dimensions.Y / 2, 1)

		i.Buffers[l] = buffer.New(buffer.Float4, dimensions)
	}

	i.Buffers = buffers
}

func (i *Image) allocateMipLevelsDownTo(bottom math.Vector2i) {
	numMipLevels := countMipLevelsDownTo(i.Buffers[0].Dimensions(), bottom)
	
	i.allocateMipLevels(numMipLevels)
}

func countMipLevels(dimensions math.Vector2i) uint32 {
	m := math.Maxi(dimensions.X, dimensions.Y)

	return 1 + uint32(gomath.Log2(float64(m)))
}

func countMipLevelsDownTo(top, bottom math.Vector2i) uint32 {
	numMipLevels := uint32(0)

	for {
		if top.X < bottom.X || top.Y < bottom.Y {
			break
		}

		top.X = math.Maxi(top.X / 2, 1)
		top.Y = math.Maxi(top.Y / 2, 1)

		numMipLevels++
	}

	return numMipLevels
}