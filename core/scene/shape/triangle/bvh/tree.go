package bvh

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/shape/triangle/primitive"
	"github.com/Opioid/scout/base/math"
	_ "math"
	_ "fmt"
)

type Tree struct {
	root buildNode
}


func (t *Tree) Intersect(ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, indices []uint32, vertices []geometry.Vertex, intersection *primitive.Intersection) bool {
	return t.root.intersect(ray, vertices, intersection)
}

func (t *Tree) IntersectP(ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, indices []uint32, vertices []geometry.Vertex) bool {
	return t.root.intersectP(ray, vertices)
}

func intersectTriangle(v0, v1, v2 math.Vector3, ray *math.OptimizedRay) (bool, primitive.Coordinates) {
	e1 := v1.Sub(v0)
	e2 := v2.Sub(v0)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1.0 / det

	tvec := ray.Origin.Sub(v0)

	c := primitive.Coordinates{}
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