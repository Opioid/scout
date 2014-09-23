package kd

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/base/math"
	_ "math"
	_ "fmt"
)

type Tree struct {
	root buildNode
}

type Intersection struct {
	T, U, V float32
	Index uint32
}

func (t *Tree) Intersect(ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, indices []uint32, vertices []geometry.Vertex, intersection *Intersection) bool {
	return t.root.intersect(ray, boundingMinT, boundingMaxT, indices, vertices, intersection)
}

func (t *Tree) IntersectP(ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, indices []uint32, vertices []geometry.Vertex) bool {
	return t.root.intersectP(ray, boundingMinT, boundingMaxT, indices, vertices)
}

func intersectTriangle(v0, v1, v2 math.Vector3, ray *math.OptimizedRay, thit, u, v *float32) bool {
	e1 := v1.Sub(v0)
	e2 := v2.Sub(v0)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1.0 / det

	tvec := ray.Origin.Sub(v0)
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

func intersectTriangleP(v0, v1, v2 math.Vector3, ray *math.OptimizedRay) bool {
	e1 := v1.Sub(v0)
	e2 := v2.Sub(v0)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1.0 / det

	tvec := ray.Origin.Sub(v0)
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