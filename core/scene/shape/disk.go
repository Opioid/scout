package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
)

type Disk struct {
	aabb bounding.AABB
}

func NewDisk() *Disk {
	d := new(Disk)
	d.aabb = bounding.MakeAABB(math.MakeVector3(-1.0, -1.0, 0.0), math.MakeVector3(1.0, 1.0, 0.0))
	return d
}

func (disk *Disk) Intersect(transformation *entity.ComposedTransformation, ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, 
							intersection *geometry.Intersection) (bool, float32) {
	normal := transformation.Rotation.Row(2)

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	thit := -(numer / denom)
	
	if thit > ray.MinT && thit < ray.MaxT {
		p := ray.Point(thit)
		k := p.Sub(transformation.Position)
		l := k.Dot(k)

		radius := transformation.Scale.X

		if l <= radius * radius {
			intersection.Epsilon = 5e-4 * thit

			intersection.P = p
			intersection.T = transformation.Rotation.Row(0)
			intersection.B = transformation.Rotation.Row(1)	
			intersection.N = normal

			sk := k.Div(radius)

			u := transformation.Rotation.Row(0).Dot(sk)
			intersection.UV.X = (u + 1.0) * 0.5

			v := transformation.Rotation.Row(1).Dot(sk)
			intersection.UV.Y = (v + 1.0) * 0.5

			intersection.MaterialId = 0

			return true, thit
		}
	} 

	return false, 0.0
}

func (disk *Disk) IntersectP(transformation *entity.ComposedTransformation, ray *math.OptimizedRay, boundingMinT, boundingMaxT float32) bool {
	normal := transformation.Rotation.Row(2)

	d := -normal.Dot(transformation.Position)

	denom := normal.Dot(ray.Direction)

	numer := normal.Dot(ray.Origin) + d

	thit := -(numer / denom)
	
	if thit >= ray.MinT && thit < ray.MaxT {
		p := ray.Point(thit)
		k := p.Sub(transformation.Position)
		l := k.Dot(k)

		radius := transformation.Scale.X

		if l <= radius * radius {
			return true
		}
	} 

	return false	
}

func (d *Disk) AABB() *bounding.AABB {
	return &d.aabb
}

func (d *Disk) IsComplex() bool {
	return false
}

func (d *Disk) IsFinite() bool {
	return true
}