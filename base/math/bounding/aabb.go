package bounding

import (
	"github.com/Opioid/scout/base/math"
    gomath "math"
	_ "fmt"
)

type AABB struct {
	Min, Max math.Vector3
}

func MakeAABB() AABB {
    min := math.Vector3{ gomath.MaxFloat32,  gomath.MaxFloat32,  gomath.MaxFloat32}
    max := math.Vector3{-gomath.MaxFloat32, -gomath.MaxFloat32, -gomath.MaxFloat32}

    return AABB{min, max}
}

func (b *AABB) Position() math.Vector3 {
    return b.Min.Add(b.Max).Scale(0.5)
}

func (b *AABB) Halfsize() math.Vector3 {
    return b.Max.Sub(b.Min).Scale(0.5)
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
    right := m.Row(0).Vector3()
    xa := right.Scale(b.Min.X)
    xb := right.Scale(b.Max.X)

    up := m.Row(1).Vector3()
    ya := up.Scale(b.Min.Y)
    yb := up.Scale(b.Max.Y)

    dir := m.Row(2).Vector3()
    za := dir.Scale(b.Min.Z)
    zb := dir.Scale(b.Max.Z)

    translation := m.Row(3).Vector3()
    other.Min = xa.Min(xb).Add(ya.Min(yb)).Add(za.Min(zb)).Add(translation)
    other.Max = xa.Max(xb).Add(ya.Max(yb)).Add(za.Max(zb)).Add(translation)
}

func (b *AABB) Intersect(ray *math.OptimizedRay) bool {
    tx1 := (b.Min.X - ray.Origin.X) * ray.ReciprocalDirection.X
    tx2 := (b.Max.X - ray.Origin.X) * ray.ReciprocalDirection.X

    tmin := math.Max(ray.MinT, math.Min(tx1, tx2))
    tmax := math.Min(ray.MaxT, math.Max(tx1, tx2))

    ty1 := (b.Min.Y - ray.Origin.Y) * ray.ReciprocalDirection.Y
    ty2 := (b.Max.Y - ray.Origin.Y) * ray.ReciprocalDirection.Y

    tmin = math.Max(tmin, math.Min(ty1, ty2))
    tmax = math.Min(tmax, math.Max(ty1, ty2))

    tz1 := (b.Min.Z - ray.Origin.Z) * ray.ReciprocalDirection.Z
    tz2 := (b.Max.Z - ray.Origin.Z) * ray.ReciprocalDirection.Z

    tmin = math.Max(tmin, math.Min(tz1, tz2))
    tmax = math.Min(tmax, math.Max(tz1, tz2))

    return tmax >= math.Max(0.0, tmin) 
}

func (b *AABB) Merge(other *AABB) AABB {
    return AABB{b.Min.Min(other.Min), b.Max.Max(other.Max)}
}