package bounding

import (
	"github.com/Opioid/scout/base/math"
    gomath "math"
	_ "fmt"
)

type AABB struct {
	Bounds [2]math.Vector3
}

func MakeEmptyAABB() AABB {
	b := AABB{}
    // min
	b.Bounds[0] = math.MakeVector3( gomath.MaxFloat32,  gomath.MaxFloat32,  gomath.MaxFloat32)
    
    // max
	b.Bounds[1] = math.MakeVector3(-gomath.MaxFloat32, -gomath.MaxFloat32, -gomath.MaxFloat32)

	return b
}

func MakeAABB(min, max math.Vector3) AABB {
	b := AABB{}
	b.Bounds[0] = min
	b.Bounds[1] = max
	return b
}

func (b *AABB) Position() math.Vector3 {
	return b.Bounds[0].Add(b.Bounds[1]).Scale(0.5)
}

func (b *AABB) Halfsize() math.Vector3 {
	return b.Bounds[1].Sub(b.Bounds[0]).Scale(0.5)
}

func (b *AABB) Transform(m *math.Matrix4x4, other *AABB) {
/*	*other = *b

	var xa = m.Right * boundingBox.Min.X;
    var xb = m.Right * boundingBox.Max.X;
 
    var ya = m.Up * boundingBox.Min.Y;
    var yb = m.Up * boundingBox.Max.Y;
 
    var za = m.Backward * boundingBox.Min.Z;
    var zb = m.Backward * boundingBox.Max.Z;
 
    return new BoundingBox(
        Vector3.Min(xa, xb) + Vector3.Min(ya, yb) + Vector3.Min(za, zb) + m.Translation,
        Vector3.Max(xa, xb) + Vector3.Max(ya, yb) + Vector3.Max(za, zb) + m.Translation
    );
*/
	right := m.Right()
	xa := right.Scale(b.Bounds[0].X)
	xb := right.Scale(b.Bounds[1].X)

	up := m.Up()
	ya := up.Scale(b.Bounds[0].Y)
	yb := up.Scale(b.Bounds[1].Y)

	dir := m.Direction()
	za := dir.Scale(b.Bounds[0].Z)
	zb := dir.Scale(b.Bounds[1].Z)

	translation := m.Translation()
	other.Bounds[0] = xa.Min(xb).Add(ya.Min(yb)).Add(za.Min(zb)).Add(translation)
	other.Bounds[1] = xa.Max(xb).Add(ya.Max(yb)).Add(za.Max(zb)).Add(translation)
}

func (b *AABB) IntersectP(ray *math.OptimizedRay) bool {
	tmin := (b.Bounds[    ray.Sign[0]].X - ray.Origin.X) * ray.ReciprocalDirection.X
	tmax := (b.Bounds[1 - ray.Sign[0]].X - ray.Origin.X) * ray.ReciprocalDirection.X

	tymin := (b.Bounds[    ray.Sign[1]].Y - ray.Origin.Y) * ray.ReciprocalDirection.Y
	tymax := (b.Bounds[1 - ray.Sign[1]].Y - ray.Origin.Y) * ray.ReciprocalDirection.Y

	if tmin > tymax || tymin > tmax {
		return false
	}

	if tymin > tmin {
		tmin = tymin
	}

	if tymax < tmax {
		tmax = tymax
	}

	tzmin := (b.Bounds[    ray.Sign[2]].Z - ray.Origin.Z) * ray.ReciprocalDirection.Z
	tzmax := (b.Bounds[1 - ray.Sign[2]].Z - ray.Origin.Z) * ray.ReciprocalDirection.Z

	if tmin > tzmax || tzmin > tmax {
		return false
	}

	if tzmin > tmin {
		tmin = tzmin
	}

	if tzmax < tmax {
		tmax = tzmax
	}

	return tmin < ray.MaxT && tmax > ray.MinT 
}

func (b *AABB) Intersect(ray *math.OptimizedRay) (bool, float32, float32) {
	tmin := (b.Bounds[    ray.Sign[0]].X - ray.Origin.X) * ray.ReciprocalDirection.X
	tmax := (b.Bounds[1 - ray.Sign[0]].X - ray.Origin.X) * ray.ReciprocalDirection.X

	tymin := (b.Bounds[    ray.Sign[1]].Y - ray.Origin.Y) * ray.ReciprocalDirection.Y
	tymax := (b.Bounds[1 - ray.Sign[1]].Y - ray.Origin.Y) * ray.ReciprocalDirection.Y

	if tmin > tymax || tymin > tmax {
		return false, 0.0, 0.0
	}

	if tymin > tmin {
		tmin = tymin
	}

	if tymax < tmax {
		tmax = tymax
	}

	tzmin := (b.Bounds[    ray.Sign[2]].Z - ray.Origin.Z) * ray.ReciprocalDirection.Z
	tzmax := (b.Bounds[1 - ray.Sign[2]].Z - ray.Origin.Z) * ray.ReciprocalDirection.Z

	if tmin > tzmax || tzmin > tmax {
		return false, 0.0, 0.0
	}

	if tzmin > tmin {
		tmin = tzmin
    }

	if tzmax < tmax {
		tmax = tzmax
	}

	return tmin < ray.MaxT && tmax > ray.MinT, tmin, tmax
}

func (b *AABB) Merge(other *AABB) AABB {
	return MakeAABB(b.Bounds[0].Min(other.Bounds[0]), b.Bounds[1].Max(other.Bounds[1]))
}