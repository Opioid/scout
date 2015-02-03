package primitive

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
)

type IndexTriangle struct {
	A, B, C uint32
	MaterialIndex uint32
}

func MakeIndexTriangle(a, b, c, materialIndex uint32) IndexTriangle {
	return IndexTriangle{a, b, c, materialIndex}
}

type Triangle struct {
	A, B, C geometry.Vertex
	MaterialIndex uint32
}

func MakeTriangle(a, b, c *geometry.Vertex, materialIndex uint32) Triangle {
	return Triangle{*a, *b, *c, materialIndex}
}

func (t *Triangle) Intersect(ray *math.OptimizedRay) (bool, Coordinates) {
	e1 := t.B.P.Sub(t.A.P)
	e2 := t.C.P.Sub(t.A.P)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1.0 / det

	tvec := ray.Origin.Sub(t.A.P)

	c := Coordinates{}
	c.U = tvec.Dot(pvec) * invDet

	if c.U < 0.0 || c.U > 1.0 {
		return false, c
	}

	qvec := tvec.Cross(e1)
	c.V = ray.Direction.Dot(qvec) * invDet

	if c.V < 0.0 || c.U + c.V > 1.0 {
		return false, c
	}

	c.T = e2.Dot(qvec) * invDet

	if c.T > ray.MinT && c.T < ray.MaxT {
		return true, c
	} 

	return false, c
}

func (t *Triangle) IntersectP(ray *math.OptimizedRay) bool {
	e1 := t.B.P.Sub(t.A.P)
	e2 := t.C.P.Sub(t.A.P)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1.0 / det

	tvec := ray.Origin.Sub(t.A.P)
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

func (t *Triangle) Interpolate(u, v float32) (math.Vector3, math.Vector3, math.Vector2) {
	w := 1.0 - u - v

	return t.A.N.Scale(w). Add(t.B.N.Scale(u)). Add(t.C.N.Scale(v)).Normalized(),
		   t.A.T.Scale(w). Add(t.B.T.Scale(u)). Add(t.C.T.Scale(v)).Normalized(),
		   t.A.UV.Scale(w).Add(t.B.UV.Scale(u)).Add(t.C.UV.Scale(v))
}

func (t *Triangle) InterpolatePosition(u, v float32) math.Vector3 {
	w := 1.0 - u - v

	return t.A.P.Scale(w).Add(t.B.P.Scale(u)).Add(t.C.P.Scale(v))
}