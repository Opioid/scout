package camera

import (
	"github.com/Opioid/scout/core/rendering/film"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	gomath "math"
	_ "fmt"
)

type Perspective struct {
	projectiveCamera
	lensRadius float32
	focalDistance float32	
	fov float32
	leftTop math.Vector3
	dx, dy math.Vector3
}

func NewPerspective(lensRadius, focalDistance, fov float32, dimensions math.Vector2, film film.Film) *Perspective {
	p := new(Perspective)
	p.lensRadius = lensRadius
	p.focalDistance = focalDistance
	p.fov = fov
	p.dimensions = calculateDimensions(dimensions, film)
	p.film = film
//	p.Entity.Transformation.ObjectToWorld.SetIdentity()
	return p
}

func (p *Perspective) UpdateView() {
//	p.Entity.Transformation.Update()

	ratio := p.dimensions.X / p.dimensions.Y

	z := ratio * gomath.Pi / p.fov * 0.5

	p.leftTop   = math.MakeVector3(-ratio,  1.0, z)
	rightTop   := math.MakeVector3( ratio,  1.0, z)
	leftBottom := math.MakeVector3(-ratio, -1.0, z)

	p.dx = rightTop.Sub(p.leftTop).Div(float32(p.film.Dimensions().X))
	p.dy = leftBottom.Sub(p.leftTop).Div(float32(p.film.Dimensions().Y))
}

func (p *Perspective) GenerateRay(sample *sampler.CameraSample, ray *math.OptimizedRay) {
	direction := p.leftTop.Add(p.dx.Scale(sample.Coordinates.X)).Add(p.dy.Scale(sample.Coordinates.Y))

	r := math.Ray{math.MakeVector3(0.0, 0.0, 0.0), direction, 0.0, 1000.0}

	if p.lensRadius > 0.0 {
		lensUv := math.DiskSample_uniform(sample.LensUv.X, sample.LensUv.Y).Scale(p.lensRadius)

		ft := p.focalDistance / r.Direction.Z
		focus := r.Point(ft)

		r.Origin = lensUv
		r.Direction = focus.Sub(r.Origin)
	}

	transformation := p.TransformationAt(sample.Time)

	ray.Origin = transformation.ObjectToWorld.TransformPoint(r.Origin)
	ray.SetDirection(transformation.ObjectToWorld.TransformVector3(r.Direction.Normalized()))

	ray.MaxT  = 1000.0
	ray.Depth = 0
}