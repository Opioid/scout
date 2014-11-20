package texture

import (
	"github.com/Opioid/scout/base/math"
)

const (
	Texture_1D   = iota
	Texture_2D   = iota
	Texture_3D   = iota
	Texture_cube = iota
)

type description struct {
	textureType uint32
	format uint32

	dimensions math.Vector3i

	numLayers uint32
	numSamples uint32

	shaderResource uint8
}

type texture struct {
	
}