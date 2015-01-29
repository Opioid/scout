package light

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/scene/shape/triangle"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	_ "fmt" 
)

type Mesh struct {
	light

	mesh *triangle.Mesh

	numTriangles float32
}

func NewMesh(shape shape.Shape) *Mesh {
	l := Mesh{}
	l.prop.SetVisible(prop.IsLight, true)
	l.prop.Shape = shape

	l.mesh = shape.(*triangle.Mesh)

	l.numTriangles = float32(l.mesh.NumTriangles())

	return &l
}

func (l *Mesh) Sample(transformation *math.ComposedTransformation, p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample {
	l.prop.TransformationAt(time, transformation)

	sample1d := sampler.GenerateSample1D(0, subsample)
	sample2d := sampler.GenerateSample2D(0, subsample)

	index := uint32(l.numTriangles * sample1d - 0.00001)

	ls := l.mesh.InterpolatedPosition(index, sample2d.X, sample2d.Y)
	ws := transformation.ObjectToWorld.TransformPoint(ls)

	v := ws.Sub(p)

	d := v.SquaredLength()
	i := 1.0 / d
	t := math32.Sqrt(d)

	result := Sample{Energy: l.color.Scale(i * l.lumen), L: v.Div(t), T: t}

	return result
}