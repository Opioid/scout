package tonemapping

import (
	"github.com/Opioid/scout/base/math"
)

type filmic struct {
	linearWhite math.Vector3
}

func NewFilmic(linearWhite math.Vector3) *filmic {
	f := new(filmic)
	f.linearWhite = linearWhite
	return f
}

func (f *filmic) Tonemap(color math.Vector3) math.Vector3 {
	numerator   := tonemapFunction(color) 
	denominator := tonemapFunction(f.linearWhite)

	return numerator.DivV(denominator)
}

// Function used by the Uncharte2 tone mapping curve
func tonemapFunction(color math.Vector3) math.Vector3 {
	/*
	float A = 0.15f;
	float B = 0.50f;
	float C = 0.10f;
	float D = 0.20f;
	float E = 0.02f;
	float F = 0.30f;
	*/

	A := float32(0.22)
	B := float32(0.30)
	C := float32(0.10)
	D := float32(0.20)
	E := float32(0.01)
	F := float32(0.30)

	A_color := color.Scale(A)

	return ((color.Mul(A_color.AddS(C * B)).AddS(D * E)).DivV(color.Mul(A_color.AddS(B)).AddS(D * F))).SubS(E / F)


// return ((color * (A * color + C * B) + D * E) / (color * (A * color + B) + D * F)) - E / F;
}