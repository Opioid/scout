package triangle

import (
	"github.com/Opioid/scout/core/scene/shape/triangle/bvh"
	_ "github.com/Opioid/scout/core/scene/shape/triangle/kd"
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/shape/triangle/primitive"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	_ "fmt"
)

type Mesh struct {
	aabb bounding.AABB

	tree bvh.Tree
}

func NewMesh(aabb bounding.AABB, tree bvh.Tree) *Mesh {
	m := Mesh{aabb: aabb, tree: tree}

	return &m
}

func (m *Mesh) Intersect(transformation *math.ComposedTransformation, ray, tray *math.OptimizedRay,
						 boundingMinT, boundingMaxT float32, intersection *geometry.Intersection) (bool, float32) {
	tray.Origin = transformation.WorldToObject.TransformPoint(ray.Origin)
	tray.SetDirection(transformation.WorldToObject.TransformVector3(ray.Direction))
	tray.MinT = ray.MinT
	tray.MaxT = ray.MaxT
	tray.Time = ray.Time
	tray.Depth = ray.Depth

	pi := primitive.Intersection{}
//	pi.T = ray.MaxT

	hit := m.tree.Intersect(tray, boundingMinT, boundingMaxT, &pi)

	if hit {
		thit := pi.T
		intersection.Epsilon = 3e-3 * thit

		intersection.P = ray.Point(thit)

		intersection.MaterialIndex = m.tree.Triangles[pi.Index].MaterialIndex
		n, t, uv := m.tree.Triangles[pi.Index].Interpolate(pi.U, pi.V)

	//	intersection.N = transformation.WorldToObject.TransposedTransformVector3(n).Normalized()
	//	intersection.T = transformation.WorldToObject.TransposedTransformVector3(t).Normalized()

		intersection.N = transformation.Rotation.TransformVector3(n)
		intersection.T = transformation.Rotation.TransformVector3(t)

		intersection.B = intersection.N.Cross(intersection.T)
		intersection.UV = uv

		return hit, thit
	}

	return false, 0.0
}

func (m *Mesh) IntersectP(transformation *math.ComposedTransformation, ray, tray *math.OptimizedRay,
					      boundingMinT, boundingMaxT float32) bool {
	tray.Origin = transformation.WorldToObject.TransformPoint(ray.Origin)
	tray.SetDirection(transformation.WorldToObject.TransformVector3(ray.Direction))
	tray.MinT = ray.MinT
	tray.MaxT = ray.MaxT

	return m.tree.IntersectP(tray, boundingMinT, boundingMaxT)
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

func (m *Mesh) NumTriangles() uint32 {
	return m.tree.NumTriangles()
}

func (m *Mesh) InterpolatedPosition(triangleIndex uint32, u, v float32) math.Vector3 {
	return m.tree.Triangles[triangleIndex].InterpolatePosition(u, v)
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