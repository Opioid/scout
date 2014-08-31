package shape

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	gomath "math"
	_ "fmt"
)

type vertex struct {
	p, n math.Vector3
	uv math.Vector2
}

type triangleMesh struct {
	indices []uint32

	vertices []vertex

	aabb bounding.AABB
}

func NewTriangleMesh(numIndices, numVertices uint32) *triangleMesh {
	m := new(triangleMesh)
	m.indices = make([]uint32, numIndices)
	m.vertices = make([]vertex, numVertices)
	return m
}

func (m *triangleMesh) Intersect(transformation *entity.ComposedTransformation, ray *math.OptimizedRay, thit *float32, epsilon *float32, dg *DifferentialGeometry) bool {
	oray := *ray
	oray.Origin = transformation.WorldToObject.TransformPoint(ray.Origin)
	oray.Direction = transformation.WorldToObject.TransformVector(ray.Direction)

//	fmt.Println(transformation.WorldToObject)

	type intersectionResult struct {
		t, u, v float32
		index uint32
	}

	closestHit := intersectionResult{t: ray.MaxT}
	hasHit := false
	for i, len := uint32(0), uint32(len(m.indices)); i < len; i += 3 {
		hit := intersectionResult{index: i}
		if intersectTriangle(m.vertices[m.indices[i + 0]].p, m.vertices[m.indices[i + 1]].p, m.vertices[m.indices[i + 2]].p, &oray, &hit.t, &hit.u, &hit.v) {
			if hit.t <= closestHit.t {
				closestHit = hit
				hasHit = true
			}
		}
	}

	if hasHit {
		*thit = closestHit.t
		*epsilon = 5e-4 * *thit

		dg.P = ray.Point(*thit)

		interpolateVertices(&m.vertices[m.indices[closestHit.index + 0]],
			          	 	&m.vertices[m.indices[closestHit.index + 1]],
			           		&m.vertices[m.indices[closestHit.index + 2]],
			           		closestHit.u, closestHit.v,
			           		&dg.Nn, &dg.UV)

		dg.Nn = transformation.WorldToObject.TransposedTransformVector(dg.Nn)

		return true
	}

	return false
}

func (m *triangleMesh) IntersectP(transformation *entity.ComposedTransformation, ray *math.OptimizedRay) bool {
	oray := *ray
	oray.Origin = transformation.WorldToObject.TransformPoint(ray.Origin)
	oray.Direction = transformation.WorldToObject.TransformVector(ray.Direction)

	for i, len := uint32(0), uint32(len(m.indices)); i < len; i += 3 {
		if intersectTriangleP(m.vertices[m.indices[i + 0]].p, m.vertices[m.indices[i + 1]].p, m.vertices[m.indices[i + 2]].p, &oray) {
			return true
		}
	}

	return false
}

func (m *triangleMesh) AABB() *bounding.AABB {
	return &m.aabb
}

func (m *triangleMesh) IsComplex() bool {
	return true
}

func (m *triangleMesh) IsFinite() bool {
	return true
}

func (m *triangleMesh) setIndex(index, value uint32) {
	m.indices[index] = value
}

func (m *triangleMesh) setPosition(index uint32, p math.Vector3) {
	m.vertices[index].p = p
}

func (m *triangleMesh) setNormal(index uint32, n math.Vector3) {
	m.vertices[index].n = n
}

func (m *triangleMesh) setUV(index uint32, uv math.Vector2) {
	m.vertices[index].uv = uv
}

func (m *triangleMesh) compile() {
	min := math.Vector3{ gomath.MaxFloat32,  gomath.MaxFloat32,  gomath.MaxFloat32}
	max := math.Vector3{-gomath.MaxFloat32, -gomath.MaxFloat32, -gomath.MaxFloat32}
	
	for _, v := range m.vertices {
		min = v.p.Min(min)
		max = v.p.Max(max)
	}

	m.aabb = bounding.MakeAABB(min, max)
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

func interpolateVertices(a, b, c *vertex, u, v float32, n *math.Vector3, uv *math.Vector2) {
	w := 1.0 - u - v
	
	*n  = a.n.Scale(w).Add(b.n.Scale(u)).Add(c.n.Scale(v))
	*uv = a.uv.Scale(w).Add(b.uv.Scale(u)).Add(c.uv.Scale(v))
}