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
	TextureType uint32
	Format uint32

	Dimensions math.Vector3i

	NumLayers    uint32
	NumMipLevels uint32
	NumSamples   uint32

	ShaderResource uint8
}

type texture struct {
	
}