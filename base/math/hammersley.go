package math

func radicalInverseVdC(bits uint32) float32 {
	bits = (bits << 16) | (bits >> 16)
	bits = ((bits & 0x55555555) << 1) | ((bits & 0xAAAAAAAA) >> 1)
	bits = ((bits & 0x33333333) << 2) | ((bits & 0xCCCCCCCC) >> 2)
	bits = ((bits & 0x0F0F0F0F) << 4) | ((bits & 0xF0F0F0F0) >> 4)
	bits = ((bits & 0x00FF00FF) << 8) | ((bits & 0xFF00FF00) >> 8)
	return float32(bits) * 2.3283064365386963e-10; // / 0x100000000
}

func Hammersley(i, numSamples uint32) Vector2 {
	return MakeVector2(float32(i) / float32(numSamples), radicalInverseVdC(i))
}