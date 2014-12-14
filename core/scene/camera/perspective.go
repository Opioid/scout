package camera

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/rendering/film"
	"github.com/Opioid/scout/base/math"
	gomath "math"
	_ "fmt"
)

type Perspective struct {
	genericCamera
	fov float32
	leftTop math.Vector3
	dx, dy math.Vector3
}

func NewPerspective(fov float32, dimensions math.Vector2, film film.Film) *Perspective {
	p := new(Perspective)
	p.film = film
	p.dimensions = calculateDimensions(dimensions, film)
	p.fov = fov
	p.Entity.Transformation.ObjectToWorld.SetIdentity()
	return p
}

func (p *Perspective) UpdateView() {
//	p.Entity.Transformation.Update()

	ratio := p.dimensions.X / p.dimensions.Y

	z := ratio * gomath.Pi / p.fov * 0.5

	corners := []math.Vector3{
		math.MakeVector3(-ratio,  1.0, z),
		math.MakeVector3( ratio,  1.0, z),
		math.MakeVector3(-ratio, -1.0, z),
	}

	p.leftTop   = p.Entity.Transformation.ObjectToWorld.TransformPoint(corners[0])
	rightTop   := p.Entity.Transformation.ObjectToWorld.TransformPoint(corners[1])
	leftBottom := p.Entity.Transformation.ObjectToWorld.TransformPoint(corners[2])

	p.dx = rightTop.Sub(p.leftTop)
	p.dy = leftBottom.Sub(p.leftTop)
}

func (p *Perspective) Transformation() *entity.ComposedTransformation {
	return &p.Entity.Transformation
}

func (p *Perspective) Position() math.Vector3 {
	return p.Entity.Transformation.Position
}

func (p *Perspective) Film() film.Film {
	return p.film
}

func (p *Perspective) GenerateRay(coordinates math.Vector2, ray *math.OptimizedRay) {
	x := coordinates.X / float32(p.film.Dimensions().X)
	y := coordinates.Y / float32(p.film.Dimensions().Y)

	ray.Origin = p.Position()

	direction := p.leftTop.Add(p.dx.Scale(x)).Add(p.dy.Scale(y)).Sub(ray.Origin)
	ray.SetDirection(direction.Normalized())

	ray.MaxT  = 1000.0
	ray.Depth = 0
}