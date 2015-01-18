package triangle

import (
	"github.com/Opioid/scout/core/scene/shape/triangle/bvh"
	_ "github.com/Opioid/scout/core/scene/shape/triangle/kd"
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/shape/triangle/primitive"
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	_ "fmt"
)

type Mesh struct {
	aabb bounding.AABB

	tree bvh.Tree1
}

func NewMesh(aabb bounding.AABB, tree bvh.Tree1) *Mesh {
	m := Mesh{aabb: aabb, tree: tree}

	return &m
}

func (m *Mesh) Intersect(transformation *entity.ComposedTransformation, ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, 
						 dg *geometry.Differential) (bool, float32, float32) {
	oray := *ray
	oray.Origin = transformation.WorldToObject.TransformPoint(ray.Origin)
	oray.SetDirection(transformation.WorldToObject.TransformVector3(ray.Direction))

	intersection := primitive.Intersection{}
	intersection.T = ray.MaxT

	hit := m.tree.Intersect(&oray, boundingMinT, boundingMaxT, &intersection)

	if hit {
		thit := intersection.T
		epsilon := 5e-3 * thit

		dg.P = ray.Point(thit)

		dg.MaterialId = m.tree.Triangles[intersection.Index].MaterialId

		dg.N, dg.T, dg.UV = m.tree.Triangles[intersection.Index].Interpolate(intersection.U, intersection.V)

		dg.N = transformation.WorldToObject.TransposedTransformVector3(dg.N)
		dg.T = transformation.WorldToObject.TransposedTransformVector3(dg.T)

		dg.B = dg.N.Cross(dg.T)

		return hit, thit, epsilon
	}

	return false, 0.0, 0.0
}

func (m *Mesh) IntersectP(transformation *entity.ComposedTransformation, ray *math.OptimizedRay, boundingMinT, boundingMaxT float32) bool {
	oray := *ray
	oray.Origin = transformation.WorldToObject.TransformPoint(ray.Origin)
	oray.SetDirection(transformation.WorldToObject.TransformVector3(ray.Direction))

	return m.tree.IntersectP(&oray, boundingMinT, boundingMaxT)
}

func (m *Mesh) AABB() *bounding.AABB {
	return &m.aabb
}

func (m *Mesh) IsComplex() bool {
	return true
}

func (m *Mesh) IsFinite() bool {
	return true
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

func interpolateVertices(tri primitive.IndexTriangle, vertices []geometry.Vertex, u, v float32) (math.Vector3, math.Vector3, math.Vector2) {
	w := 1.0 - u - v
	
	return vertices[tri.A].N. Scale(w).Add(vertices[tri.B].N. Scale(u)).Add(vertices[tri.C].N. Scale(v)).Normalized(),
		   vertices[tri.A].T. Scale(w).Add(vertices[tri.B].T. Scale(u)).Add(vertices[tri.C].T. Scale(v)).Normalized(),
	       vertices[tri.A].UV.Scale(w).Add(vertices[tri.B].UV.Scale(u)).Add(vertices[tri.C].UV.Scale(v))
}