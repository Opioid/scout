package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
)

type vertex struct {
	p, n math.Vector3
	uv math.Vector2
}

type triangleMesh struct {
	indices []uint32

	vertices []vertex
}

func NewTriangleMesh(numIndices, numVertices uint32) *triangleMesh {
	m := new(triangleMesh)
	m.indices = make([]uint32, numIndices)
	m.vertices = make([]vertex, numVertices)
	return m
}

func (m *triangleMesh) Intersect(transformation *entity.ComposedTransformation, ray *math.Ray, thit *float32, epsilon *float32, dg *DifferentialGeometry) bool {
	closestHit := ray.MaxT
//	hitIndex := uint32(0)
	hit := false
	for i, len := uint32(0), uint32(len(m.indices)); i < len; i += 3 {
		var t float32
		if intersectTriangle(&m.vertices[m.indices[i + 0]], &m.vertices[m.indices[i + 1]], &m.vertices[m.indices[i + 2]], ray, &t) {
			if t <= closestHit {
				closestHit = t
		//		hitIndex = i
				hit = true
			}
		}
	}

	if hit {
		*thit = closestHit
		*epsilon = 5e-4 * *thit

		dg.P = ray.Point(*thit)
		dg.Nn = math.Vector3{0.0, 0.0, -1.0}

		return true
	}

	return false
}

func (m *triangleMesh) IntersectP(transformation *entity.ComposedTransformation, ray *math.Ray) bool {
	for i, len := uint32(0), uint32(len(m.indices)); i < len; i += 3 {
		if intersectTriangleP(&m.vertices[m.indices[i + 0]], &m.vertices[m.indices[i + 1]], &m.vertices[m.indices[i + 2]], ray) {
			return true
		}
	}

	return false
}

func (m *triangleMesh) setIndex(index, value uint32) {
	m.indices[index] = value
}

func (m *triangleMesh) setPosition(index uint32, p math.Vector3) {
	m.vertices[index].p = p
}

func intersectTriangle(v0, v1, v2 *vertex, ray *math.Ray, thit *float32) bool {
	e1 := v1.p.Sub(v0.p)
	e2 := v2.p.Sub(v0.p)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1.0 / det

	tvec := ray.Origin.Sub(v0.p)
	u := tvec.Dot(pvec) * invDet

	if u < 0.0 || u > 1.0 {
		return false
	}

	qvec := tvec.Cross(e1)
	v := ray.Direction.Dot(qvec) * invDet

	if v < 0.0 || u + v > 1.0 {
		return false
	}

	*thit = e2.Dot(qvec) * invDet

	if *thit >= ray.MinT && *thit < ray.MaxT {
		return true
	} 

	return false
}

func intersectTriangleP(v0, v1, v2 *vertex, ray *math.Ray) bool {
	e1 := v1.p.Sub(v0.p)
	e2 := v2.p.Sub(v0.p)

	pvec := ray.Direction.Cross(e2)

	det := e1.Dot(pvec)
	invDet := 1.0 / det

	tvec := ray.Origin.Sub(v0.p)
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

	if thit >= ray.MinT && thit < ray.MaxT {
		return true
	} 

	return false
}