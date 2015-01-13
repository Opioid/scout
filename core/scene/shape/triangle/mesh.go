package triangle

import (
	"github.com/Opioid/scout/core/scene/shape/triangle/bvh"
	_ "github.com/Opioid/scout/core/scene/shape/triangle/kd"
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/shape/triangle/primitive"
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	gomath "math"
	_ "fmt"
)

type Mesh struct {
	triangles []primitive.IndexTriangle
	vertices []geometry.Vertex

	aabb bounding.AABB

	tree bvh.Tree
}

func NewMesh(numTriangles, numVertices uint32) *Mesh {
	m := new(Mesh)
	m.triangles = make([]primitive.IndexTriangle, numTriangles)
	m.vertices  = make([]geometry.Vertex, numVertices)
	return m
}

func (m *Mesh) Intersect(transformation *entity.ComposedTransformation, ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, 
						 dg *geometry.Differential) (bool, float32, float32) {
	oray := *ray
	oray.Origin = transformation.WorldToObject.TransformPoint(ray.Origin)
	oray.SetDirection(transformation.WorldToObject.TransformVector3(ray.Direction))

	intersection := primitive.Intersection{ T: ray.MaxT }

	hit := m.tree.Intersect(&oray, boundingMinT, boundingMaxT, m.triangles, m.vertices, &intersection)

	if hit {
		thit := intersection.T
		epsilon := 5e-4 * thit

		dg.P = ray.Point(thit)

		intersection.Triangle.Interpolate(intersection.U, intersection.V, &dg.N, &dg.T, &dg.UV)

	//	interpolateVertices(intersection.IndexTriangle, m.vertices, intersection.U, intersection.V, &dg.N, &dg.T, &dg.UV)

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

	return m.tree.IntersectP(&oray, boundingMinT, boundingMaxT, m.triangles, m.vertices)
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

func (m *Mesh) SetTriangle(index uint32, tri primitive.IndexTriangle) {
	m.triangles[index] = tri
}

func (m *Mesh) SetPosition(index uint32, p math.Vector3) {
	m.vertices[index].P = p
}

func (m *Mesh) SetNormal(index uint32, n math.Vector3) {
	m.vertices[index].N = n
}

func (m *Mesh) SetTangentAndSign(index uint32, t math.Vector3, s float32) {
	m.vertices[index].T = t
	m.vertices[index].BitangentSign = s
}

func (m *Mesh) SetUV(index uint32, uv math.Vector2) {
	m.vertices[index].UV = uv
}

func (m *Mesh) Compile() {
	min := math.MakeVector3( gomath.MaxFloat32,  gomath.MaxFloat32,  gomath.MaxFloat32)
	max := math.MakeVector3(-gomath.MaxFloat32, -gomath.MaxFloat32, -gomath.MaxFloat32)
	
	for _, v := range m.vertices {
		min = v.P.Min(min)
		max = v.P.Max(max)
	}

	m.aabb = bounding.MakeAABB(min, max)

	builder := bvh.Builder{}
	builder.Build(m.triangles, m.vertices, 8, &m.tree)

	m.triangles = nil
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

func interpolateVertices(tri *primitive.IndexTriangle, vertices []geometry.Vertex, u, v float32, n, t *math.Vector3, uv *math.Vector2) {
	w := 1.0 - u - v
	
	*n  = vertices[tri.A].N. Scale(w).Add(vertices[tri.B].N. Scale(u)).Add(vertices[tri.C].N. Scale(v)).Normalized()
	*t  = vertices[tri.A].T. Scale(w).Add(vertices[tri.B].T. Scale(u)).Add(vertices[tri.C].T. Scale(v)).Normalized()
	*uv = vertices[tri.A].UV.Scale(w).Add(vertices[tri.B].UV.Scale(u)).Add(vertices[tri.C].UV.Scale(v))
}