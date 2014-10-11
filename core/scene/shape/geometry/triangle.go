package geometry

import (
	"github.com/Opioid/scout/base/math"
)

type Triangle struct {
	a, b, c Vertex
}

func MakeTriangle(a, b, c *Vertex) Triangle {
	return Triangle{*a, *b, *c}

}

func (t *Triangle) Intersect(ray *math.OptimizedRay, thit, u, v *float32) bool {
	e1 := t.b.P.Sub(t.a.P)
	e2 := t.c.P.Sub(t.a.P)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1.0 / det

	tvec := ray.Origin.Sub(t.a.P)
	*u = tvec.Dot(pvec) * invDet

	if *u < 0.0 || *u > 1.0 {
		return false
	}

	qvec := tvec.Cross(e1)
	*v = ray.Direction.Dot(qvec) * invDet

	if *v < 0.0 || *u + *v > 1.0 {
		return false
	}

	*thit = e2.Dot(qvec) * invDet

	if *thit > ray.MinT && *thit < ray.MaxT {
		return true
	} 

	return false
}

func (t *Triangle) IntersectP(ray *math.OptimizedRay) bool {
	e1 := t.b.P.Sub(t.a.P)
	e2 := t.c.P.Sub(t.a.P)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1.0 / det

	tvec := ray.Origin.Sub(t.a.P)
	u := tvec.Dot(pvec) * invDet

	if u < 0.0 || u > 1.0 {
		return false
	}

	qvec := tvec.Cross(e1)
	v := ray.Direction.Dot(qvec) * invDet

	if v < 0.0 || u + v > 1.0 {
		return false
	}

	thit := e2.Dot(qvec) * invDet

	if thit > ray.MinT && thit < ray.MaxT {
		return true
	} 

	return false
}

func (t *Triangle) Interpolate(u, v float32, n *math.Vector3, uv *math.Vector2) {
	w := 1.0 - u - v
	
	*n  = t.a.N.Scale(w).Add(t.b.N.Scale(u)).Add(t.c.N.Scale(v))
	*uv = t.a.UV.Scale(w).Add(t.b.UV.Scale(u)).Add(t.c.UV.Scale(v))
}