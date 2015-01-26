package light

import (
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/scene/shape/triangle"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/math32"
)

type Mesh struct {
	light

	mesh *triangle.Mesh
}

func NewMesh(shape shape.Shape) *Mesh {
	l := Mesh{}
	l.prop.Shape = shape
	return &l
}

func (l *Mesh) Samples(transformation *math.ComposedTransformation, p math.Vector3, time float32, subsample, maxSamples uint32, sampler sampler.Sampler, samples []Sample) []Sample {
	samples = samples[:0]
/*
	l.prop.TransformationAt(time, transformation)

	result := Sample{}
*/
	return samples	
}

func (l *Mesh) Sample(transformation *math.ComposedTransformation, p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample {
	l.prop.TransformationAt(time, transformation)

	sample1d := sampler.GenerateSample1D(0, subsample)
	sample2d := sampler.GenerateSample(0, subsample)

/*
	ls := math.HemisphereSample_uniform(sample.X, sample.Y)
	ws := transformation.Rotation.TransformVector3(ls).Scale(transformation.Scale.X)

	v := transformation.Position.Add(ws).Sub(p)

	d := v.SquaredLength()
	i := 1.0 / d
*/
	result := Sample{Energy: math.MakeVector3(sample2d.X, sample2d.Y, 0.0), L: math.MakeVector3(sample1d, 0.0, 0.0)}

	return result
}