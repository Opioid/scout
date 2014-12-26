package primitive

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
)

type Triangle struct {
	A, B, C geometry.Vertex
}

func MakeTriangle(a, b, c *geometry.Vertex) Triangle {
	return Triangle{*a, *b, *c}

}

func (t *Triangle) Intersect(ray *math.OptimizedRay, thit, u, v *float32) bool {
	e1 := t.B.P.Sub(t.A.P)
	e2 := t.C.P.Sub(t.A.P)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1 / det

	tvec := ray.Origin.Sub(t.A.P)
	*u = tvec.Dot(pvec) * invDet

	if *u < 0 || *u > 1 {
		return false
	}

	qvec := tvec.Cross(e1)
	*v = ray.Direction.Dot(qvec) * invDet

	if *v < 0 || *u + *v > 1 {
		return false
	}

	*thit = e2.Dot(qvec) * invDet

	if *thit > ray.MinT && *thit < ray.MaxT {
		return true
	} 

	return false
}

func (t *Triangle) IntersectP(ray *math.OptimizedRay) bool {
	e1 := t.B.P.Sub(t.A.P)
	e2 := t.C.P.Sub(t.A.P)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1 / det

	tvec := ray.Origin.Sub(t.A.P)
	u := tvec.Dot(pvec) * invDet

	if u < 0 || u > 1 {
		return false
	}

	qvec := tvec.Cross(e1)
	v := ray.Direction.Dot(qvec) * invDet

	if v < 0 || u + v > 1 {
		return false
	}

	thit := e2.Dot(qvec) * invDet

	if thit > ray.MinT && thit < ray.MaxT {
		return true
	} 

	return false
}

func (tri *Triangle) Interpolate(u, v float32, n, t *math.Vector3, uv *math.Vector2) {
	w := 1 - u - v
	
	*n  = tri.A.N.Scale(w).Add(tri.B.N.Scale(u)).Add(tri.C.N.Scale(v)).Normalized()
	*t  = tri.A.T.Scale(w).Add(tri.B.T.Scale(u)).Add(tri.C.T.Scale(v)).Normalized()
	*uv = tri.A.UV.Scale(w).Add(tri.B.UV.Scale(u)).Add(tri.C.UV.Scale(v))
}