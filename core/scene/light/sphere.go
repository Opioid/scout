package light

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	gomath "math"
	_ "fmt"
)

type Sphere struct {
	light
}

func NewSphere(shape shape.Shape) *Sphere {
	l := Sphere{}
	l.prop.SetVisible(prop.IsLight, true)
	l.prop.Shape = shape
	return &l
}

const (
	hemispherePdf = 1.0 / (2.0 * gomath.Pi)
	hemisphereArea = 2.0 * gomath.Pi
)

func (l *Sphere) Sample(transformation *math.ComposedTransformation, p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample {

	l.prop.TransformationAt(time, transformation)

	sample := sampler.GenerateSample2D(0, subsample)
	ls := math.SampleHemisphereUniform(sample.X, sample.Y)

	z := p.Sub(transformation.Position).Normalized()
	cs := math.MakeCoordinateSystemMatrix3x3(z)

	n := cs.TransformVector3(ls)

	ws := transformation.Position.Add(n.Scale(transformation.Scale.X))

	v := ws.Sub(p)

	d := v.SquaredLength()
	t := math32.Sqrt(d)
	w := v.Div(t)

	nDotW := n.Dot(w.Scale(-1.0))

	if nDotW < 0.0 {
		// In this case no light will reach p, so we could make an early out.
		// I think it also means the sample we picked was bad, 
		// so this could probably be optimized away with the cone thingy described in pbrt.
		d = 0.0
	}

	radiusSquare := transformation.Scale.X * transformation.Scale.X

	result := Sample{Energy: l.color.Scale(l.lumen), L: w, T: t, Pdf: d / (math32.Abs(nDotW) * (radiusSquare * hemisphereArea))}

	return result

/*
	l.prop.TransformationAt(time, transformation)

	z := p.Sub(transformation.Position).Normalized()
	x, y := math.CoordinateSystem(z)

	radiusSquare := transformation.Scale.X * transformation.Scale.X

	sinThetaMax2 := radiusSquare / p.Sub(transformation.Position).SquaredLength()
	cosThetaMax := math32.Sqrt(math32.Max(0.0, 1.0 - sinThetaMax2))

	sample := sampler.GenerateSample2D(0, subsample)
	n := math.SampleOrientedConeUniform(sample.X, sample.Y, cosThetaMax, x, y, z)

	ws := transformation.Position.Add(n.Scale(transformation.Scale.X))

	v := ws.Sub(p)

	d := v.SquaredLength()
	t := math32.Sqrt(d)
	w := v.Div(t)

	nDotW := n.Dot(w.Scale(-1.0))

	if nDotW < 0.0 {
		// In this case no light will reach p, so we could make an early out.
		// I think it also means the sample we picked was bad, 
		// so this could probably be optimized away with the cone thingy described in pbrt.
		result := Sample{Energy: l.color.Scale(l.lumen), L: w, T: t, Pdf: 0.0}

		return result
	}

	result := Sample{Energy: l.color.Scale(l.lumen), L: w, T: t, Pdf: math.ConePdfUniform(cosThetaMax)}

	return result
		*/
}